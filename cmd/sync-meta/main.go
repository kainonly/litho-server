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
	"unicode"

	"server/bootstrap"
	"server/common"

	"gorm.io/gorm"
)

type actionDef struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type resourceDef struct {
	Key     string
	Label   string
	Actions []actionDef
}

type capDef struct {
	Key         string
	Description string
}

func main() {
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	apiDir := flag.String("api-dir", "api", "API 模块目录")
	commonFile := flag.String("common", "common/common.go", "能力标识定义文件路径")
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

	apiDirPath, err := resolvePath(*apiDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API 目录路径错误: %v\n", err)
		os.Exit(1)
	}

	capsFilePath, err := resolvePath(*commonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "能力标识文件路径错误: %v\n", err)
		os.Exit(1)
	}

	// 扫描各模块 common.go（Resource/Label）及操作文件（I* 常量）生成资源定义
	resources, err := parseModules(apiDirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析模块失败: %v\n", err)
		os.Exit(1)
	}

	if err := syncResources(db, resources); err != nil {
		fmt.Fprintf(os.Stderr, "同步资源失败: %v\n", err)
		os.Exit(1)
	}

	// 从 common/caps.go 中读取 C* 常量生成能力标识定义
	caps, err := extractCapConstants(capsFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析能力标识失败: %v\n", err)
		os.Exit(1)
	}

	if err := syncCaps(db, caps); err != nil {
		fmt.Fprintf(os.Stderr, "同步能力标识失败: %v\n", err)
		os.Exit(1)
	}
}

// parseModules 扫描 apiDir 下各子模块，通过 common.go 的 Key/Label 常量
// 和各操作文件中的 I* 常量，构建资源定义列表。
func parseModules(apiDir string) ([]resourceDef, error) {
	entries, err := os.ReadDir(apiDir)
	if err != nil {
		return nil, err
	}

	var defs []resourceDef
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		moduleDir := filepath.Join(apiDir, entry.Name())
		commonFile := filepath.Join(moduleDir, "common.go")
		if _, err := os.Stat(commonFile); err != nil {
			continue
		}

		key, label, err := parseCommonConstants(commonFile)
		if err != nil || key == "" || label == "" {
			continue
		}

		actions, err := parseModuleActions(moduleDir)
		if err != nil {
			return nil, fmt.Errorf("模块 %s: %w", entry.Name(), err)
		}

		defs = append(defs, resourceDef{
			Key:     key,
			Label:   label,
			Actions: actions,
		})
	}
	return defs, nil
}

// parseCommonConstants 从 common.go 中提取 Key 和 Label 常量值。
func parseCommonConstants(filePath string) (key, label string, err error) {
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
			case "Key":
				key = val
			case "Label":
				label = val
			}
		}
		return true
	})
	return
}

// extractCapConstants 从 caps.go 中提取所有顶层字符串常量，常量名作为 key，常量值作为 description。
func extractCapConstants(filePath string) ([]capDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var result []capDef
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
			desc := strings.Trim(lit.Value, `"`)
			result = append(result, capDef{Key: name.Name, Description: desc})
		}
		return true
	})
	return result, nil
}

// parseModuleActions 扫描模块目录下所有 .go 文件，提取 I* 常量，
// 将常量名（如 ICreate）转为 action value（如 create），常量值作为 label。
func parseModuleActions(moduleDir string) ([]actionDef, error) {
	goFiles, err := filepath.Glob(filepath.Join(moduleDir, "*.go"))
	if err != nil {
		return nil, err
	}

	var actions []actionDef
	for _, file := range goFiles {
		if strings.HasSuffix(file, "common.go") {
			continue
		}
		found, err := extractIConstants(file)
		if err != nil {
			return nil, err
		}
		actions = append(actions, found...)
	}
	return actions, nil
}

// extractIConstants 从单个 Go 文件中提取以 I 开头的字符串常量，
// 将常量名转为 snake_case action value，常量字符串值作为 label。
func extractIConstants(filePath string) ([]actionDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var result []actionDef
	ast.Inspect(f, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range spec.Names {
			if !strings.HasPrefix(name.Name, "I") {
				continue
			}
			if i >= len(spec.Values) {
				continue
			}
			lit, ok := spec.Values[i].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			label := strings.Trim(lit.Value, `"`)
			value := pascalToSnake(name.Name[1:]) // 去掉前缀 I，转 snake_case
			result = append(result, actionDef{Label: label, Value: value})
		}
		return true
	})
	return result, nil
}

// pascalToSnake 将 PascalCase 字符串转换为 snake_case，例如 SetRoles -> set_roles。
func pascalToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func syncResources(db *gorm.DB, defs []resourceDef) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, def := range defs {
			actions := def.Actions
			if actions == nil {
				actions = []actionDef{}
			}
			actionsJSON, err := json.Marshal(actions)
			if err != nil {
				return err
			}

			var count int64
			if err := tx.Raw(`SELECT COUNT(*) FROM "resource" WHERE id = ?`, def.Key).Scan(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				if err := tx.Exec(
					`UPDATE "resource" SET label = ?, actions = ?::jsonb WHERE id = ?`,
					def.Label, string(actionsJSON), def.Key,
				).Error; err != nil {
					return err
				}
				fmt.Printf("更新资源 %s (%s)\n", def.Key, def.Label)
			} else {
				if err := tx.Exec(
					`INSERT INTO "resource" (id, label, actions) VALUES (?, ?, ?::jsonb)`,
					def.Key, def.Label, string(actionsJSON),
				).Error; err != nil {
					return err
				}
				fmt.Printf("新增资源 %s (%s)\n", def.Key, def.Label)
			}
		}
		return nil
	})
}

func syncCaps(db *gorm.DB, defs []capDef) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, def := range defs {
			var count int64
			if err := tx.Raw(`SELECT COUNT(*) FROM "cap" WHERE id = ?`, def.Key).Scan(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				if err := tx.Exec(
					`UPDATE "cap" SET description = ? WHERE id = ?`,
					def.Description, def.Key,
				).Error; err != nil {
					return err
				}
				fmt.Printf("更新能力标识 %s\n", def.Key)
			} else {
				if err := tx.Exec(
					`INSERT INTO "cap" (id, description) VALUES (?, ?)`,
					def.Key, def.Description,
				).Error; err != nil {
					return err
				}
				fmt.Printf("新增能力标识 %s\n", def.Key)
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
