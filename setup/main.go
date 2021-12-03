package main

import (
	"api/bootstrap"
	"api/model"
	"context"
	"github.com/weplanx/go/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"log"
	"os"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.SetValues,
			bootstrap.UseDatabase,
		),
		fx.Invoke(func(db *mongo.Database) {
			var err error
			ctx := context.Background()
			collection := db.Collection("schema")
			if _, err = collection.InsertMany(ctx, []interface{}{
				model.Schema{
					Key:         "page",
					Label:       "动态页面",
					Kind:        "manual",
					Description: "",
					System:      true,
					Fields:      nil,
				},
				model.Schema{
					Key:    "role",
					Label:  "权限组",
					Kind:   "collection",
					System: true,
					Fields: []model.Field{
						{
							Key:      "key",
							Label:    "权限代码",
							Type:     "text",
							Required: true,
							Unique:   true,
							System:   true,
						},
						{
							Key:      "name",
							Label:    "权限名称",
							Type:     "text",
							Required: true,
							System:   true,
						},
						{
							Key:      "status",
							Label:    "状态",
							Type:     "bool",
							Required: true,
							System:   true,
						},
						{
							Key:    "description",
							Label:  "描述",
							Type:   "text",
							System: true,
						},
						{
							Key:     "pages",
							Label:   "页面",
							Type:    "reference",
							Default: "'[]'",
							System:  true,
							Option: model.FieldOption{
								Mode:   "manual",
								Target: "page",
							},
						},
					},
				},
				model.Schema{
					Label: "成员",
					Key:   "admin",
					Kind:  "collection",
					Fields: []model.Field{
						{
							Key:      "username",
							Label:    "用户名",
							Type:     "text",
							Required: true,
							Unique:   true,
							System:   true,
						},
						{
							Key:      "password",
							Label:    "密码",
							Type:     "password",
							Required: true,
							Private:  true,
							System:   true,
						},
						{
							Key:      "status",
							Label:    "状态",
							Type:     "bool",
							Required: true,
							System:   true,
						},
						{
							Key:      "roles",
							Label:    "权限",
							Type:     "reference",
							Required: true,
							Default:  "'[]'",
							System:   true,
							Option: model.FieldOption{
								Mode:   "many",
								Target: "role",
								To:     "key",
							},
						},
						{
							Key:    "name",
							Label:  "姓名",
							Type:   "text",
							System: true,
						},
						{
							Key:    "email",
							Label:  "邮件",
							Type:   "email",
							System: true,
						},
						{
							Key:    "phone",
							Label:  "联系方式",
							Type:   "text",
							System: true,
						},
						{
							Key:     "avatar",
							Label:   "头像",
							Type:    "media",
							Default: "'[]'",
							System:  true,
						},
					},
					System: true,
				},
			}); err != nil {
				log.Fatalln(err)
			}
			if _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.M{
					"key": 1,
				},
				Options: options.
					Index().
					SetUnique(true).
					SetName("key_idx"),
			}); err != nil {
				log.Fatalln(err)
			}
			// >>>
			collection = db.Collection("page")
			if _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{"parent", 1},
					{"fragment", 1},
				},
				Options: options.Index().SetName("parent_fragment_idx").SetUnique(true),
			}); err != nil {
				log.Fatalln(err)
			}
			if _, err = collection.InsertOne(ctx, model.Page{
				Parent:   "root",
				Fragment: "dashboard",
				Name:     "仪表盘",
				Nav:      true,
				Icon:     "dashboard",
				Sort:     0,
				Router:   "manual",
			}); err != nil {
				log.Fatalln(err)
			}
			center, err := collection.InsertOne(ctx, model.Page{
				Parent:   "root",
				Fragment: "center",
				Name:     "个人中心",
				Nav:      false,
				Icon:     "",
				Sort:     0,
				Router:   "",
			})
			if err != nil {
				log.Fatalln(err)
			}
			if _, err = collection.InsertMany(ctx, []interface{}{
				model.Page{
					Parent:   center.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "profile",
					Name:     "我的信息",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "manual",
				},
				model.Page{
					Parent:   center.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "notification",
					Name:     "消息通知",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "manual",
				},
			}); err != nil {
				log.Fatalln(err)
			}
			settings, err := collection.InsertOne(ctx, model.Page{
				Parent:   "root",
				Fragment: "settings",
				Name:     "设置",
				Nav:      true,
				Icon:     "setting",
				Sort:     0,
				Router:   "",
			})
			if err != nil {
				log.Fatalln(err)
			}
			roleViewFields := []model.ViewFields{
				{
					Key:     "name",
					Label:   "权限名称",
					Display: true,
				},
				{
					Key:     "key",
					Label:   "权限代码",
					Display: true,
				},
				{
					Key:     "description",
					Label:   "描述",
					Display: true,
				},
				{
					Key:     "pages",
					Label:   "页面",
					Display: true,
				},
			}
			adminViewFields := []model.ViewFields{
				{
					Key:     "username",
					Label:   "用户名",
					Display: true,
				},
				{
					Key:     "password",
					Label:   "密码",
					Display: true,
				},
				{
					Key:     "status",
					Label:   "状态",
					Display: true,
				},
				{
					Key:     "roles",
					Label:   "权限",
					Display: true,
				},
				{
					Key:     "name",
					Label:   "姓名",
					Display: true,
				},
				{
					Key:     "email",
					Label:   "邮件",
					Display: true,
				},
				{
					Key:     "phone",
					Label:   "联系方式",
					Display: true,
				},
				{
					Key:     "avatar",
					Label:   "头像",
					Display: true,
				},
			}
			if _, err = collection.InsertMany(ctx, []interface{}{
				model.Page{
					Parent:   settings.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "schema",
					Name:     "模式管理",
					Nav:      true,
					Icon:     "",
					Sort:     0,
					Router:   "manual",
				},
				model.Page{
					Parent:   settings.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "page",
					Name:     "页面管理",
					Nav:      true,
					Icon:     "",
					Sort:     0,
					Router:   "manual",
				},
				model.Page{
					Parent:   settings.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "role",
					Name:     "权限管理",
					Nav:      true,
					Icon:     "",
					Sort:     0,
					Router:   "table",
					Option: model.PageOption{
						Schema: "role",
						Fields: roleViewFields,
					},
				},
				model.Page{
					Parent:   settings.InsertedID.(primitive.ObjectID).Hex(),
					Fragment: "admin",
					Name:     "成员管理",
					Nav:      true,
					Icon:     "",
					Sort:     0,
					Router:   "table",
					Option: model.PageOption{
						Schema: "admin",
						Fields: adminViewFields,
					},
				},
			}); err != nil {
				log.Fatalln(err)
			}
			var role map[string]interface{}
			if err = collection.FindOne(ctx, bson.M{
				"parent":   settings.InsertedID.(primitive.ObjectID).Hex(),
				"fragment": "role",
			}).Decode(&role); err != nil {
				log.Fatalln(err)
			}
			if _, err = collection.InsertMany(ctx, []interface{}{
				model.Page{
					Parent:   role["_id"].(primitive.ObjectID).Hex(),
					Fragment: "create",
					Name:     "创建资源",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "form",
					Option: model.PageOption{
						Schema: "role",
						Fetch:  false,
						Fields: roleViewFields,
					},
				},
				model.Page{
					Parent:   role["_id"].(primitive.ObjectID).Hex(),
					Fragment: "update",
					Name:     "更新资源",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "form",
					Option: model.PageOption{
						Schema: "role",
						Fetch:  true,
						Fields: roleViewFields,
					},
				},
			}); err != nil {
				log.Fatalln(err)
			}
			var admin map[string]interface{}
			if err = collection.FindOne(ctx, bson.M{
				"parent":   settings.InsertedID.(primitive.ObjectID).Hex(),
				"fragment": "admin",
			}).Decode(&admin); err != nil {
				log.Fatalln(err)
			}
			if _, err = collection.InsertMany(ctx, []interface{}{
				model.Page{
					Parent:   admin["_id"].(primitive.ObjectID).Hex(),
					Fragment: "create",
					Name:     "创建资源",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "form",
					Option: model.PageOption{
						Schema: "admin",
						Fields: adminViewFields,
					},
				},
				model.Page{
					Parent:   admin["_id"].(primitive.ObjectID).Hex(),
					Fragment: "update",
					Name:     "更新资源",
					Nav:      false,
					Icon:     "",
					Sort:     0,
					Router:   "form",
					Option: model.PageOption{
						Schema: "admin",
						Fetch:  true,
						Fields: adminViewFields,
					},
				},
			}); err != nil {
				log.Fatalln(err)
			}
			// >>>
			if _, err = db.Collection("role").InsertOne(ctx, bson.M{
				"key":         "*",
				"name":        "超级管理员",
				"status":      true,
				"description": "",
				"pages":       bson.A{},
			}); err != nil {
				log.Fatalln(err)
			}
			var hash string
			if hash, err = password.Create("pass@VAN1234"); err != nil {
				log.Fatalln(err)
			}
			if _, err = db.Collection("admin").InsertOne(ctx, bson.M{
				"username": "admin",
				"password": hash,
				"status":   true,
				"roles":    bson.A{"*"},
				"name":     "超级管理员",
				"email":    "",
				"phone":    "",
				"avatar":   "",
			}); err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}),
	)
}
