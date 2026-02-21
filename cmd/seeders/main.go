package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"server/bootstrap"
	"server/common"
	"server/model"

	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type modelFactory func() any

type userSeed struct {
	model.User
	Org  string `json:"org"`
	Role string `json:"role"`
}

type routeSeed struct {
	model.Route
	Children []routeSeed `json:"children"`
}

var seedModels = map[string]modelFactory{
	"cap":   func() any { return model.Cap{} },
	"org":   func() any { return model.Org{} },
	"role":  func() any { return model.Role{} },
	"route": func() any { return model.Route{} },
	"user":  func() any { return model.User{} },
}

var modelAliases = map[string]string{
	"caps":   "cap",
	"orgs":   "org",
	"roles":  "role",
	"routes": "route",
	"users":  "user",
}

func main() {
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	dataDir := flag.String("data", "cmd/seeders/data", "种子数据目录")
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

	if err := truncateSeedTables(db); err != nil {
		fmt.Fprintf(os.Stderr, "清空种子表失败: %v\n", err)
		os.Exit(1)
	}

	dir, err := resolvePath(*dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "数据目录错误: %v\n", err)
		os.Exit(1)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取数据目录失败: %v\n", err)
		os.Exit(1)
	}

	// 按依赖顺序排列：routes 须先于 roles（roles 动态注入 route ID）
	seedOrder := []string{"caps", "orgs", "routes", "roles", "users"}
	orderIndex := func(name string) int {
		base := strings.TrimSuffix(name, filepath.Ext(name))
		for i, s := range seedOrder {
			if strings.EqualFold(base, s) {
				return i
			}
		}
		return len(seedOrder)
	}

	var jsonFiles []string
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(dir, entry.Name()))
		}
	}
	sort.Slice(jsonFiles, func(i, j int) bool {
		bi := filepath.Base(jsonFiles[i])
		bj := filepath.Base(jsonFiles[j])
		oi, oj := orderIndex(bi), orderIndex(bj)
		if oi != oj {
			return oi < oj
		}
		return jsonFiles[i] < jsonFiles[j]
	})

	if len(jsonFiles) == 0 {
		fmt.Fprintln(os.Stdout, "未找到任何 .json 文件")
		return
	}

	for _, filePath := range jsonFiles {
		if err := seedFile(db, filePath); err != nil {
			fmt.Fprintf(os.Stderr, "导入失败 %s: %v\n", filePath, err)
			os.Exit(1)
		}
	}
}

func truncateSeedTables(db *gorm.DB) error {
	tables := []string{
		"cap",
		"org",
		"role",
		"route",
		"user",
	}

	quoted := make([]string, 0, len(tables))
	for _, table := range tables {
		quoted = append(quoted, `"`+table+`"`)
	}

	return db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, strings.Join(quoted, ", "))).Error
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

func seedFile(db *gorm.DB, filePath string) error {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	modelKey := normalizeModelKey(base)
	factory, ok := seedModels[modelKey]
	if !ok {
		return seedCompositeFile(db, filePath)
	}

	if modelKey == "user" {
		return seedUsersWithLookup(db, filePath, base)
	}

	if modelKey == "route" {
		return seedRoutesWithTree(db, filePath, base)
	}

	if modelKey == "role" {
		return seedRolesWithRoutes(db, filePath, base)
	}

	return seedWithFactory(db, filePath, base, modelKey, factory)
}

func seedCompositeFile(db *gorm.DB, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("未注册的种子模型: %s", filepath.Base(filePath))
	}

	if len(payload) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空文件: %s\n", filePath)
		return nil
	}

	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		modelKey := normalizeModelKey(key)
		factory, ok := seedModels[modelKey]
		if !ok {
			return fmt.Errorf("未注册的种子模型: %s", key)
		}

		if err := seedWithFactory(db, filePath, key, modelKey, factory, payload[key]); err != nil {
			return err
		}
	}

	return nil
}

func seedWithFactory(db *gorm.DB, filePath, label, modelKey string, factory modelFactory, raw ...json.RawMessage) error {
	var data []byte
	var err error
	if len(raw) > 0 {
		data = raw[0]
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return err
		}
	}

	records, err := decodeJSON(data, factory())
	if err != nil {
		return err
	}

	if err := applySeedTransforms(modelKey, records); err != nil {
		return err
	}

	count := countRecords(records)
	if count == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(records).Error; err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "导入成功 %s (%s): %d\n", filePath, label, count)
	return nil
}

func normalizeModelKey(key string) string {
	if v, ok := modelAliases[key]; ok {
		return v
	}
	return key
}

func seedRolesWithRoutes(db *gorm.DB, filePath, label string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var seeds []model.Role
	if err := json.Unmarshal(data, &seeds); err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	var routeIDs []string
	if err := db.Model(&model.Route{}).Pluck("id", &routeIDs).Error; err != nil {
		return fmt.Errorf("查询 route ID 失败: %w", err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := range seeds {
			if seeds[i].ID == "" {
				seeds[i].ID = help.SID()
			}
			if seeds[i].Strategy == nil {
				seeds[i].Strategy = &common.Object{}
			}
			(*seeds[i].Strategy)["routes"] = routeIDs
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seeds[i]).Error; err != nil {
				return err
			}
		}
		fmt.Fprintf(os.Stdout, "导入成功 %s (%s): %d\n", filePath, label, len(seeds))
		return nil
	})
}

func seedRoutesWithTree(db *gorm.DB, filePath, label string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var seeds []routeSeed
	if err := json.Unmarshal(data, &seeds); err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		count, err := insertRoutes(tx, seeds, "0", "")
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "导入成功 %s (%s): %d\n", filePath, label, count)
		return nil
	})
}

func insertRoutes(tx *gorm.DB, seeds []routeSeed, pid string, nav string) (int, error) {
	count := 0
	for i := range seeds {
		if seeds[i].ID == "" {
			seeds[i].ID = help.SID()
		}
		if seeds[i].Pid == "" {
			seeds[i].Pid = pid
		}
		if seeds[i].Nav == "" {
			seeds[i].Nav = nav
		}
		if seeds[i].Sort == 0 {
			seeds[i].Sort = int16(i + 1)
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seeds[i].Route).Error; err != nil {
			return count, err
		}
		count++

		if len(seeds[i].Children) > 0 {
			n, err := insertRoutes(tx, seeds[i].Children, seeds[i].ID, seeds[i].Nav)
			if err != nil {
				return count, err
			}
			count += n
		}
	}
	return count, nil
}

func seedUsersWithLookup(db *gorm.DB, filePath, label string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var seeds []userSeed
	if err := json.Unmarshal(data, &seeds); err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := range seeds {
			if seeds[i].ID == "" {
				seeds[i].ID = help.SID()
			}

			if seeds[i].Org != "" && seeds[i].OrgID == "" {
				var org model.Org
				if err := tx.Where("name = ?", seeds[i].Org).Take(&org).Error; err != nil {
					return fmt.Errorf("找不到组织 %q: %w", seeds[i].Org, err)
				}
				seeds[i].OrgID = org.ID
			}

			if seeds[i].Role != "" && seeds[i].RoleID == "" {
				var role model.Role
				if err := tx.Where("name = ?", seeds[i].Role).Take(&role).Error; err != nil {
					return fmt.Errorf("找不到角色 %q: %w", seeds[i].Role, err)
				}
				seeds[i].RoleID = role.ID
			}

			if seeds[i].Password != "" {
				hash, err := passlib.Hash(seeds[i].Password)
				if err != nil {
					return err
				}
				seeds[i].Password = hash
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seeds[i].User).Error; err != nil {
				return err
			}
		}

		fmt.Fprintf(os.Stdout, "导入成功 %s (%s): %d\n", filePath, label, len(seeds))
		return nil
	})
}

func applySeedTransforms(modelKey string, records any) error {
	return fillMissingIDs(records)
}

func fillMissingIDs(records any) error {
	value := reflect.ValueOf(records)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("记录类型不合法")
	}
	value = value.Elem()
	if value.Kind() != reflect.Slice {
		return fmt.Errorf("记录类型不合法")
	}

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			continue
		}
		idField := item.FieldByName("ID")
		if !idField.IsValid() || !idField.CanSet() || idField.Kind() != reflect.String {
			continue
		}
		if idField.String() == "" {
			idField.SetString(help.SID())
		}
	}

	return nil
}

func decodeJSON(data []byte, prototype any) (any, error) {
	elemType := reflect.TypeOf(prototype)
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	sliceType := reflect.SliceOf(elemType)
	slicePtr := reflect.New(sliceType)

	if err := json.Unmarshal(data, slicePtr.Interface()); err == nil {
		return slicePtr.Interface(), nil
	}

	itemPtr := reflect.New(elemType)
	if err := json.Unmarshal(data, itemPtr.Interface()); err != nil {
		return nil, err
	}

	singleSlice := reflect.MakeSlice(sliceType, 1, 1)
	singleSlice.Index(0).Set(itemPtr.Elem())
	singleSlicePtr := reflect.New(sliceType)
	singleSlicePtr.Elem().Set(singleSlice)
	return singleSlicePtr.Interface(), nil
}

func countRecords(records any) int {
	value := reflect.ValueOf(records)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.Slice {
		return value.Len()
	}
	return 0
}
