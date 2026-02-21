package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/bootstrap"
	"server/common"
	"server/model"

	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type resourceDef struct {
	Path    string
	Actions []string
}

func main() {
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	apiFile := flag.String("api", "api/api.go", "API 路由文件路径")
	flag.Parse()

	values, err := loadValues(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	db, err := bootstrap.UseGorm(values)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	apiPath, err := resolvePath(*apiFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API 文件路径错误: %v\n", err)
		os.Exit(1)
	}

	defs, err := parseRoutes(apiPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析路由失败: %v\n", err)
		os.Exit(1)
	}

	if err := syncResources(db, defs); err != nil {
		fmt.Fprintf(os.Stderr, "同步资源失败: %v\n", err)
		os.Exit(1)
	}
}

// parseRoutes 解析 api.go 中所有 POST m.POST("/:resource/:action", ...) 调用，
// 按 path（/:resource）分组收集 actions。
func parseRoutes(apiFilePath string) ([]resourceDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, apiFilePath, nil, 0)
	if err != nil {
		return nil, err
	}

	// path -> []action
	index := make(map[string][]string)
	order := []string{}

	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if sel.Sel.Name != "POST" {
			return true
		}
		if len(call.Args) < 1 {
			return true
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		route := strings.Trim(lit.Value, `"`)
		// 匹配 /:resource/:action 模式（两段，以 / 开头，第二段不以 _ 开头且不含 :）
		parts := strings.Split(strings.TrimPrefix(route, "/"), "/")
		if len(parts) != 2 {
			return true
		}
		resource := "/" + parts[0]
		action := parts[1]
		if strings.HasPrefix(action, ":") || strings.HasPrefix(action, "_") {
			return true
		}

		if _, exists := index[resource]; !exists {
			order = append(order, resource)
		}
		index[resource] = append(index[resource], action)
		return true
	})

	defs := make([]resourceDef, 0, len(order))
	for _, path := range order {
		defs = append(defs, resourceDef{Path: path, Actions: index[path]})
	}
	return defs, nil
}

func syncResources(db *gorm.DB, defs []resourceDef) error {
	now := time.Now()

	return db.Transaction(func(tx *gorm.DB) error {
		for _, def := range defs {
			actionsJSON, err := json.Marshal(def.Actions)
			if err != nil {
				return err
			}

			var existing model.Resource
			err = tx.Where("path = ?", def.Path).Take(&existing).Error
			if err == nil {
				if err := tx.Exec(
					`UPDATE "resource" SET updated_at = ?, actions = ?::jsonb WHERE id = ?`,
					now, string(actionsJSON), existing.ID,
				).Error; err != nil {
					return err
				}
				fmt.Printf("更新资源 %s: %v\n", def.Path, def.Actions)
			} else if err == gorm.ErrRecordNotFound {
				active := true
				id := help.SID()
				if err := tx.Exec(
					`INSERT INTO "resource" (id, created_at, updated_at, active, name, path, actions) VALUES (?, ?, ?, ?, ?, ?, ?::jsonb)`,
					id, now, now, active, strings.TrimPrefix(def.Path, "/"), def.Path, string(actionsJSON),
				).Error; err != nil {
					return err
				}
				fmt.Printf("新增资源 %s: %v\n", def.Path, def.Actions)
			} else {
				return err
			}
		}
		return nil
	})
}

func loadValues(path string) (*common.Values, error) {
	absPath, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	return bootstrap.LoadStaticValues(absPath)
}

func resolvePath(path string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}

	candidates := []string{
		filepath.Join(filepath.Dir(execPath), "..", "..", path),
		filepath.Join(filepath.Dir(execPath), path),
		path,
	}

	for _, candidate := range candidates {
		absPath, _ := filepath.Abs(candidate)
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("路径不存在: %s", path)
}
