package devtools

import (
	"bytes"
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/thoas/go-funk"
	"lab-api/model"
	"strings"
	"text/template"
)

func (x *Service) MigrateSchema(ctx context.Context) (err error) {
	tx := x.Db.WithContext(ctx)
	if tx.Migrator().HasTable(&model.Schema{}) {
		if err = tx.Migrator().DropTable(&model.Schema{}); err != nil {
			return
		}
	}
	if err = tx.AutoMigrate(&model.Schema{}); err != nil {
		return
	}
	tx.Exec("create index columns_gin on schema using gin(columns)")
	data := []model.Schema{
		{
			Key:    "resource",
			Kind:   "manual",
			System: model.True(),
		},
		{
			Key:  "role",
			Kind: "collection",
			Columns: map[string]model.Column{
				"key": {
					Label:   "权限代码",
					Type:    "varchar",
					Require: true,
					Unique:  true,
					System:  true,
				},
				"name": {
					Label:   "权限名称",
					Type:    "varchar",
					Require: true,
					System:  true,
				},
				"description": {
					Label:  "描述",
					Type:   "text",
					System: true,
				},
				"routers": {
					Label:   "路由",
					Type:    "rel",
					Default: "'[]'",
					Relation: model.Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
				"permissions": {
					Label:   "策略",
					Type:    "rel",
					Default: "'[]'",
					Relation: model.Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
			},
			System: model.True(),
		},
		{
			Key:  "admin",
			Kind: "collection",
			Columns: map[string]model.Column{
				"uuid": {
					Label:   "唯一标识",
					Type:    "uuid",
					Default: "uuid_generate_v4()",
					Require: true,
					Unique:  true,
					Private: true,
					System:  true,
				},
				"username": {
					Label:   "用户名",
					Type:    "varchar",
					Require: true,
					Unique:  true,
					System:  true,
				},
				"password": {
					Label:   "密码",
					Type:    "varchar",
					Require: true,
					Private: true,
					System:  true,
				},
				"roles": {
					Label:   "权限",
					Type:    "rel",
					Require: true,
					Default: "'[]'",
					Relation: model.Relation{
						Mode:       "many",
						Target:     "role",
						References: "key",
					},
					System: true,
				},
				"name": {
					Label:  "姓名",
					Type:   "varchar",
					System: true,
				},
				"email": {
					Label:  "邮件",
					Type:   "varchar",
					System: true,
				},
				"phone": {
					Label:  "联系方式",
					Type:   "varchar",
					System: true,
				},
				"avatar": {
					Label:   "头像",
					Type:    "array",
					Default: "'[]'",
					System:  true,
				},
				"routers": {
					Label:   "路由",
					Type:    "rel",
					Default: "'[]'",
					Relation: model.Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
				"permissions": {
					Label:   "策略",
					Type:    "rel",
					Default: "'[]'",
					Relation: model.Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
			},
			System: model.True(),
		},
	}
	return tx.Create(&data).Error
}

func (x *Service) title(s string) string {
	return strings.Title(s)
}

func (x *Service) typ(val string) string {
	switch val {
	case "int":
		return "int32"
	case "int8":
		return "int64"
	case "decimal":
		return "float64"
	case "float8":
		return "float64"
	case "varchar":
		return "string"
	case "text":
		return "string"
	case "bool":
		return "*bool"
	case "timestamptz":
		return "time.Time"
	case "uuid":
		return "uuid.UUID"
	case "object":
		return "Object"
	case "array":
		return "Array"
	case "rel":
		return "Array"
	}
	return val
}

func (x *Service) columns(columns model.Columns) string {
	var b strings.Builder
	for k, v := range columns {
		b.WriteString(x.title(k))
		b.WriteString(" ")
		b.WriteString(x.typ(v.Type))
		b.WriteString(" `")
		b.WriteString(`gorm:"type:`)
		if funk.Contains([]string{"object", "array", "rel"}, v.Type) {
			b.WriteString("jsonb")
		} else {
			b.WriteString(v.Type)
		}
		if v.Require {
			b.WriteString(`;not null`)
		}
		if v.Unique {
			b.WriteString(`;unique`)
		}
		if v.Default != "" {
			b.WriteString(`;default:`)
			b.WriteString(v.Default)
		}
		b.WriteString(`" json:"`)
		if v.Private {
			b.WriteString(`-`)
		} else {
			b.WriteString(k)
		}
		b.WriteString(`"`)
		b.WriteString("`\n")
	}
	return b.String()
}

func (x *Service) CreateModels(ctx context.Context) (buf bytes.Buffer, err error) {
	tx := x.Db.WithContext(ctx)
	var schemas []model.Schema
	if err = tx.
		Where("kind <> ?", "manual").
		Find(&schemas).Error; err != nil {
		return
	}
	var tmpl *template.Template
	if tmpl, err = template.
		New("model").
		Funcs(template.FuncMap{
			"title":   x.title,
			"columns": x.columns,
		}).
		Parse(modelTpl); err != nil {
		return
	}
	if err = tmpl.Execute(&buf, schemas); err != nil {
		return
	}
	return
}

func (x *Service) MigrateResource(ctx context.Context) (err error) {
	tx := x.Db.WithContext(ctx)
	if tx.Migrator().HasTable(&model.Resource{}) {
		if err = tx.Migrator().DropTable(&model.Resource{}); err != nil {
			return
		}
	}
	if err = tx.AutoMigrate(&model.Resource{}); err != nil {
		return
	}
	tx.Exec("create index router_gin on resource using gin(router)")
	data := []model.Resource{
		{
			Name:   "仪表盘",
			Path:   "dashboard",
			Parent: "root",
			Router: model.Router{
				Template: "manual",
			},
			Nav:  model.True(),
			Icon: "dashboard",
		},
		{
			Name:   "个人中心",
			Path:   "center",
			Parent: "root",
		},
		{
			Name:   "我的信息",
			Path:   "center/profile",
			Parent: "center",
			Router: model.Router{
				Template: "manual",
			},
		},
		{
			Name:   "消息通知",
			Path:   "center/notification",
			Parent: "center",
			Router: model.Router{
				Template: "manual",
			},
		},
		{
			Name:   "设置",
			Path:   "setting",
			Parent: "root",
			Icon:   "setting",
		},
		{
			Name:   "资源管理",
			Path:   "setting/resource",
			Parent: "setting",
			Router: model.Router{
				Template: "manual",
			},
			Nav: model.True(),
		},
		{
			Name:   "权限管理",
			Path:   "setting/role",
			Parent: "setting",
			Router: model.Router{
				Template: "list",
				Schema:   "role",
			},
			Nav: model.True(),
		},
		{
			Name:   "创建资源",
			Path:   "setting/role/create",
			Parent: "setting/role",
			Router: model.Router{
				Template: "page",
				Schema:   "role",
			},
		},
		{
			Name:   "更新资源",
			Path:   "setting/role/update",
			Parent: "setting/role",
			Router: model.Router{
				Template: "page",
				Schema:   "role",
			},
		},
		{
			Name:   "成员管理",
			Path:   "setting/admin",
			Parent: "setting",
			Router: model.Router{
				Template: "list",
				Schema:   "admin",
			},
			Nav: model.True(),
		},
		{
			Name:   "创建资源",
			Path:   "setting/admin/create",
			Parent: "setting/admin",
			Router: model.Router{
				Template: "page",
				Schema:   "admin",
			},
		},
		{
			Name:   "更新资源",
			Path:   "setting/admin/update",
			Parent: "setting/admin",
			Router: model.Router{
				Template: "page",
				Schema:   "admin",
			},
		},
	}
	return tx.Create(&data).Error
}

func (x *Service) Seeder(ctx context.Context) (err error) {
	tx := x.Db.WithContext(ctx)
	var routers model.Array
	if err = tx.Model(&model.Resource{}).Pluck("id", &routers).Error; err != nil {
		return
	}
	roles := []map[string]interface{}{
		{
			"key":         "*",
			"name":        "超级管理员",
			"description": "超级管理员拥有完整权限不能编辑，若不使用可以禁用该权限",
			"routers":     model.Array{},
			"permissions": model.Array{},
		},
		{
			"key":         "admin",
			"name":        "管理员",
			"description": "分配管理用户",
			"routers":     routers,
			"permissions": model.Array{
				"resource:*",
				"role:*",
				"admin:*",
			},
		},
	}
	if err = tx.Table("role").Create(&roles).Error; err != nil {
		return
	}
	var password string
	if password, err = argon2id.CreateHash(
		"pass@VAN1234",
		argon2id.DefaultParams,
	); err != nil {
		return
	}
	admins := []map[string]interface{}{
		{
			"username": "admin",
			"password": password,
			"roles":    model.Array{"*"},
		},
	}
	if err = tx.Table("admin").Create(&admins).Error; err != nil {
		return
	}
	return
}
