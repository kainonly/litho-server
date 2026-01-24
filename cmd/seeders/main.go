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

type orgSeed struct {
	model.Org
	Users []model.User `json:"users"`
}

type menuSeed struct {
	model.Menu
	Children []routeSeed `json:"children"`
}

type routeSeed struct {
	model.Route
}

type resourceSeed struct {
	model.Resource
	Actions []resourceActionSeed `json:"actions"`
}

type resourceActionSeed struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sort int16  `json:"sort,omitempty"`
}

var seedModels = map[string]modelFactory{
	"menu":                  func() any { return model.Menu{} },
	"org":                   func() any { return model.Org{} },
	"permission":            func() any { return model.Permission{} },
	"resource_action":       func() any { return model.ResourceAction{} },
	"resource":              func() any { return model.Resource{} },
	"role":                  func() any { return model.Role{} },
	"role_menu":             func() any { return model.RoleMenu{} },
	"role_permission":       func() any { return model.RolePermission{} },
	"role_route":            func() any { return model.RoleRoute{} },
	"route":                 func() any { return model.Route{} },
	"route_resource_action": func() any { return model.RouteResourceAction{} },
	"user":                  func() any { return model.User{} },
	"user_org_role":         func() any { return model.UserOrgRole{} },
}

var modelAliases = map[string]string{
	"menus":                  "menu",
	"orgs":                   "org",
	"permissions":            "permission",
	"resources":              "resource",
	"roles":                  "role",
	"routes":                 "route",
	"users":                  "user",
	"resource_actions":       "resource_action",
	"role_menus":             "role_menu",
	"role_permissions":       "role_permission",
	"role_routes":            "role_route",
	"route_resource_actions": "route_resource_action",
	"user_org_roles":         "user_org_role",
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

	var jsonFiles []string
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(dir, entry.Name()))
		}
	}
	sort.Strings(jsonFiles)

	if len(jsonFiles) == 0 {
		fmt.Fprintln(os.Stdout, "未找到任何 .json 文件")
		return
	}

	skipFiles, err := computeSkipFiles(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "检测跳过文件失败: %v\n", err)
		os.Exit(1)
	}

	for _, filePath := range jsonFiles {
		if skipFiles[filePath] {
			continue
		}
		if err := seedFile(db, filePath); err != nil {
			fmt.Fprintf(os.Stderr, "导入失败 %s: %v\n", filePath, err)
			os.Exit(1)
		}
	}
}

func truncateSeedTables(db *gorm.DB) error {
	tables := []string{
		"menu",
		"org",
		"permission",
		"resource",
		"resource_action",
		"role",
		"route",
		"user",
		"user_org_role",
		"role_menu",
		"role_permission",
		"role_route",
		"route_resource_action",
	}

	quoted := make([]string, 0, len(tables))
	for _, table := range tables {
		quoted = append(quoted, `"`+table+`"`)
	}

	return db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, strings.Join(quoted, ", "))).Error
}

func computeSkipFiles(dataDir string) (map[string]bool, error) {
	skip := make(map[string]bool)

	orgsPath := filepath.Join(dataDir, "orgs.json")
	if hasEmbeddedUsers(orgsPath) {
		usersPath := filepath.Join(dataDir, "users.json")
		skip[usersPath] = true
	}

	menusPath := filepath.Join(dataDir, "menus.json")
	if hasEmbeddedRoutes(menusPath) {
		routesPath := filepath.Join(dataDir, "routes.json")
		skip[routesPath] = true
	}

	return skip, nil
}

func hasEmbeddedUsers(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	seeds, err := decodeOrgSeeds(data)
	if err != nil {
		return false
	}
	for _, seed := range seeds {
		if len(seed.Users) > 0 {
			return true
		}
	}
	return false
}

func hasEmbeddedRoutes(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	seeds, err := decodeMenuSeeds(data)
	if err != nil {
		return false
	}
	for _, seed := range seeds {
		if len(seed.Children) > 0 {
			return true
		}
	}
	return false
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

	if modelKey == "org" {
		return seedOrgsWithUsers(db, filePath, base)
	}

	if modelKey == "menu" {
		return seedMenusWithRoutes(db, filePath, base)
	}

	if modelKey == "resource" {
		return seedResourcesWithActions(db, filePath, base)
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

func seedMenusWithRoutes(db *gorm.DB, filePath, label string, raw ...json.RawMessage) error {
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

	seeds, err := decodeMenuSeeds(data)
	if err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range seeds {
			if seed.ID == "" {
				seed.ID = help.SID()
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seed.Menu).Error; err != nil {
				return err
			}

			if len(seed.Children) == 0 {
				continue
			}

			for i := range seed.Children {
				if seed.Children[i].ID == "" {
					seed.Children[i].ID = help.SID()
				}
				if seed.Children[i].MenuID == "" {
					seed.Children[i].MenuID = seed.ID
				}
				if seed.Children[i].PID == "" {
					seed.Children[i].PID = "0"
				}
				if seed.Children[i].CreatedAt.IsZero() {
					seed.Children[i].CreatedAt = seed.CreatedAt
				}
				if seed.Children[i].UpdatedAt.IsZero() {
					seed.Children[i].UpdatedAt = seed.UpdatedAt
				}
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seed.Children).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func decodeMenuSeeds(data []byte) ([]menuSeed, error) {
	var seeds []menuSeed
	if err := json.Unmarshal(data, &seeds); err == nil {
		return seeds, nil
	}

	var single menuSeed
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, err
	}
	return []menuSeed{single}, nil
}

func seedResourcesWithActions(db *gorm.DB, filePath, label string, raw ...json.RawMessage) error {
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

	seeds, err := decodeResourceSeeds(data)
	if err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range seeds {
			resourceID, err := upsertResource(tx, &seed.Resource)
			if err != nil {
				return err
			}

			if len(seed.Actions) == 0 {
				continue
			}

			actions := make([]model.ResourceAction, 0, len(seed.Actions))
			for i, action := range seed.Actions {
				sortValue := action.Sort
				if sortValue == 0 {
					sortValue = int16(i + 1)
				}
				actions = append(actions, model.ResourceAction{
					ID:         help.SID(),
					ResourceID: resourceID,
					CreatedAt:  seed.CreatedAt,
					UpdatedAt:  seed.UpdatedAt,
					Sort:       sortValue,
					Active:     true,
					Name:       action.Name,
					Code:       action.Code,
				})
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&actions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func decodeResourceSeeds(data []byte) ([]resourceSeed, error) {
	var seeds []resourceSeed
	if err := json.Unmarshal(data, &seeds); err == nil {
		return seeds, nil
	}

	var single resourceSeed
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, err
	}
	return []resourceSeed{single}, nil
}

func upsertResource(tx *gorm.DB, resource *model.Resource) (string, error) {
	if resource.Code == "" {
		return "", fmt.Errorf("resource code 不能为空")
	}

	var existing model.Resource
	if err := tx.Where("code = ?", resource.Code).Take(&existing).Error; err == nil {
		return existing.ID, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if resource.ID == "" {
		resource.ID = help.SID()
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(resource).Error; err != nil {
		return "", err
	}

	return resource.ID, nil
}

func seedOrgsWithUsers(db *gorm.DB, filePath, label string, raw ...json.RawMessage) error {
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

	seeds, err := decodeOrgSeeds(data)
	if err != nil {
		return err
	}

	if len(seeds) == 0 {
		fmt.Fprintf(os.Stdout, "跳过空数据: %s (%s)\n", filePath, label)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range seeds {
			if seed.ID == "" {
				seed.ID = help.SID()
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seed.Org).Error; err != nil {
				return err
			}

			role := model.Role{
				ID:          help.SID(),
				OrgID:       seed.ID,
				Sort:        1,
				Active:      true,
				Name:        "超级管理",
				Description: "系统超级管理员",
			}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
				return err
			}

			if len(seed.Users) == 0 {
				continue
			}

			if err := applySeedTransforms("user", &seed.Users); err != nil {
				return err
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&seed.Users).Error; err != nil {
				return err
			}

			uors := make([]model.UserOrgRole, 0, len(seed.Users))
			for _, user := range seed.Users {
				uors = append(uors, model.UserOrgRole{
					ID:     help.SID(),
					UserID: user.ID,
					OrgID:  seed.ID,
					RoleID: role.ID,
				})
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&uors).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func decodeOrgSeeds(data []byte) ([]orgSeed, error) {
	var seeds []orgSeed
	if err := json.Unmarshal(data, &seeds); err == nil {
		return seeds, nil
	}

	var single orgSeed
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, err
	}
	return []orgSeed{single}, nil
}

func applySeedTransforms(modelKey string, records any) error {
	if err := fillMissingIDs(records); err != nil {
		return err
	}

	switch modelKey {
	case "user":
		users, ok := records.(*[]model.User)
		if !ok {
			return fmt.Errorf("用户数据类型不匹配")
		}
		for i := range *users {
			if (*users)[i].Password == "" {
				continue
			}
			hash, err := passlib.Hash((*users)[i].Password)
			if err != nil {
				return err
			}
			(*users)[i].Password = hash
		}
	}
	return nil
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
