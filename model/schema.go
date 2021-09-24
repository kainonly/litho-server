package model

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type Schema struct {
	ID      int64   `json:"id"`
	Key     string  `gorm:"type:varchar;not null;unique" json:"key"`
	Kind    string  `gorm:"type:varchar;not null" json:"kind"`
	Columns Columns `gorm:"type:jsonb;default:'{}'" json:"columns"`
	System  *bool   `gorm:"default:false" json:"system"`
}

type Columns map[string]Column

func (x *Columns) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Columns) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}

type Column struct {
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Default  string   `json:"default,omitempty"`
	Unique   bool     `json:"unique,omitempty"`
	Require  bool     `json:"require,omitempty"`
	Relation Relation `json:"relation,omitempty"`
	Private  bool     `json:"private,omitempty"`
	System   bool     `json:"system,omitempty"`
}

type Relation struct {
	Mode       string `json:"mode,omitempty"`
	Target     string `json:"target,omitempty"`
	References string `json:"references,omitempty"`
}

func GenerateSchema(tx *gorm.DB) (err error) {
	if tx.Migrator().HasTable(&Schema{}) {
		if err = tx.Migrator().DropTable(&Schema{}); err != nil {
			return
		}
	}
	if err = tx.AutoMigrate(&Schema{}); err != nil {
		return
	}
	tx.Exec("create index columns_gin on schema using gin(columns)")
	data := []Schema{
		{
			Key:  "resource",
			Kind: "manual",
		},
		{
			Key:  "role",
			Kind: "collection",
			Columns: map[string]Column{
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
					Relation: Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
				"permissions": {
					Label:   "策略",
					Type:    "rel",
					Default: "'[]'",
					Relation: Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
			},
			System: True(),
		},
		{
			Key:  "admin",
			Kind: "collection",
			Columns: map[string]Column{
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
					Relation: Relation{
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
					Relation: Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
				"permissions": {
					Label:   "策略",
					Type:    "rel",
					Default: "'[]'",
					Relation: Relation{
						Mode:   "manual",
						Target: "resource",
					},
					System: true,
				},
			},
			System: True(),
		},
	}
	return tx.Create(&data).Error
}
