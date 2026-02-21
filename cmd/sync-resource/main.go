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

	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type resourceDef struct {
	Path    string
	Label   string
	Actions []common.ActionDef
}

func main() {
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	apiFile := flag.String("api", "api/api.go", "API 路由文件路径")
	apiDir := flag.String("api-dir", "api", "API 模块目录")
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

	apiDirPath, err := resolvePath(*apiDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API 目录路径错误: %v\n", err)
		os.Exit(1)
	}

	// 解析各模块 common.go 获取 Resource -> Label 映射
	labelMap, err := parseModuleLabels(apiDirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析模块标签失败: %v\n", err)
		os.Exit(1)
	}

	defs, err := parseRoutes(apiPath, labelMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析路由失败: %v\n", err)
		os.Exit(1)
	}

	if err := syncResources(db, defs); err != nil {
		fmt.Fprintf(os.Stderr, "同步资源失败: %v\n", err)
		os.Exit(1)
	}
}

// parseModuleLabels 扫描 apiDir 下各子目录的 common.go，
// 提取 Resource 和 Label 常量，返回 resource path -> label 映射。
func parseModuleLabels(apiDir string) (map[string]string, error) {
	result := make(map[string]string)

	entries, err := os.ReadDir(apiDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		commonFile := filepath.Join(apiDir, entry.Name(), "common.go")
		if _, err := os.Stat(commonFile); err != nil {
			continue
		}

		resource, label, err := parseCommonConstants(commonFile)
		if err != nil || resource == "" || label == "" {
			continue
		}
		result[resource] = label
	}

	return result, nil
}

// parseCommonConstants 从 common.go 中提取 Resource 和 Label 常量值。
func parseCommonConstants(filePath string) (resource, label string, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return
	}

	ast.Inspect(f, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range spec.Names {
			if i >= len(spec.Values) {
				continue
			}
			lit, ok := spec.Values[i].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			val := strings.Trim(lit.Value, `"`)
			switch name.Name {
			case "Resource":
				resource = val
			case "Label":
				label = val
			}
		}
		return true
	})
	return
}

// parseRoutes 解析 api.go 中所有 POST m.POST("/:resource/:action", ...) 调用，
// 按 path（/:resource）分组收集 actions，并附加模块 label。
func parseRoutes(apiFilePath string, labelMap map[string]string) ([]resourceDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, apiFilePath, nil, 0)
	if err != nil {
		return nil, err
	}

	// path -> []action value
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
		actions := make([]common.ActionDef, 0, len(index[path]))
		for _, v := range index[path] {
			if def, ok := common.ActionLabels[v]; ok {
				actions = append(actions, def)
			} else {
				// 未在 ActionLabels 注册的 action，使用 value 作为 label
				actions = append(actions, common.ActionDef{Label: v, Value: v})
			}
		}
		defs = append(defs, resourceDef{
			Path:    path,
			Label:   labelMap[path],
			Actions: actions,
		})
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

			name := def.Label
			if name == "" {
				name = strings.TrimPrefix(def.Path, "/")
			}

			var existing struct{ ID string }
			err = tx.Raw(`SELECT id FROM "resource" WHERE path = ? LIMIT 1`, def.Path).Scan(&existing).Error
			if err == nil && existing.ID == "" {
				err = gorm.ErrRecordNotFound
			}
			if err == nil {
				if err := tx.Exec(
					`UPDATE "resource" SET updated_at = ?, name = ?, actions = ?::jsonb WHERE id = ?`,
					now, name, string(actionsJSON), existing.ID,
				).Error; err != nil {
					return err
				}
				fmt.Printf("更新资源 %s (%s): %v\n", def.Path, name, def.Actions)
			} else if err == gorm.ErrRecordNotFound {
				active := true
				id := help.SID()
				if err := tx.Exec(
					`INSERT INTO "resource" (id, created_at, updated_at, active, name, path, actions) VALUES (?, ?, ?, ?, ?, ?, ?::jsonb)`,
					id, now, now, active, name, def.Path, string(actionsJSON),
				).Error; err != nil {
					return err
				}
				fmt.Printf("新增资源 %s (%s): %v\n", def.Path, name, def.Actions)
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
