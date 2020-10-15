package main

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"van-api/app/types"
)

var db *gorm.DB

func main() {
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		log.Fatalln("the configuration file does not exist")
	}
	buf, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("failed to read service configuration file", err)
	}
	config := types.Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatalln("service configuration file parsing failed", err)
	}
	db, err = gorm.Open(mysql.Open(config.Mysql.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Mysql.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Api Access Control Table'",
	//).AutoMigrate(&model.Acl{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//aclData := []model.Acl{
	//	{Key: "main", Name: types.JSON{"zh_cn": "公共模块", "en_us": "Common Module"}, Write: "uploads", Read: ""},
	//	{Key: "resource", Name: types.JSON{"zh_cn": "资源控制模块", "en_us": "Resource Module"}, Write: "add,edit,delete,sort", Read: "originLists,lists,get"},
	//	{Key: "acl", Name: types.JSON{"zh_cn": "访问控制模块", "en_us": "Acl Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
	//	{Key: "policy", Name: types.JSON{"zh_cn": "策略模块", "en_us": "Policy Module"}, Write: "add,delete", Read: "originLists"},
	//	{Key: "admin", Name: types.JSON{"zh_cn": "管理员模块", "en_us": "Admin Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
	//	{Key: "role", Name: types.JSON{"zh_cn": "权限组模块", "en_us": "Role Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
	//}
	//db.Create(&aclData)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Resource Access Control Table'",
	//).AutoMigrate(&model.Resource{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//resourceData := []model.Resource{
	//	{Key: "center", Parent: "origin", Name: types.JSON{"zh_cn": "个人中心", "en_us": "Center"}},
	//	{Key: "profile", Parent: "center", Name: types.JSON{"zh_cn": "信息修改", "en_us": "Profile"}, Router: true},
	//	{Key: "system", Parent: "origin", Name: types.JSON{"zh_cn": "系统设置", "en_us": "System"}, Nav: true, Icon: "setting"},
	//	{Key: "resource-index", Parent: "system", Name: types.JSON{"zh_cn": "资源控制管理", "en_us": "Resource"}, Nav: true, Router: true, Policy: true},
	//	{Key: "resource-add", Parent: "resource-index", Name: types.JSON{"zh_cn": "资源控制新增", "en_us": "Resource Add"}, Router: true},
	//	{Key: "resource-edit", Parent: "resource-index", Name: types.JSON{"zh_cn": "资源控制修改", "en_us": "Resource Edit"}, Router: true},
	//	{Key: "acl-index", Parent: "system", Name: types.JSON{"zh_cn": "访问控制管理", "en_us": "Acl"}, Nav: true, Router: true, Policy: true},
	//	{Key: "acl-add", Parent: "acl-index", Name: types.JSON{"zh_cn": "访问控制新增", "en_us": "Acl Add"}, Router: true},
	//	{Key: "acl-edit", Parent: "acl-index", Name: types.JSON{"zh_cn": "访问控制修改", "en_us": "Acl Edit"}, Router: true},
	//	{Key: "role-index", Parent: "system", Name: types.JSON{"zh_cn": "权限组", "en_us": "Role"}, Nav: true, Router: true, Policy: true},
	//	{Key: "role-add", Parent: "role-index", Name: types.JSON{"zh_cn": "权限组新增", "en_us": "Role Add"}, Router: true},
	//	{Key: "role-edit", Parent: "role-index", Name: types.JSON{"zh_cn": "权限组修改", "en_us": "Role Edit"}, Router: true},
	//	{Key: "admin-index", Parent: "system", Name: types.JSON{"zh_cn": "管理员", "en_us": "Admin"}, Nav: true, Router: true, Policy: true},
	//	{Key: "admin-add", Parent: "admin-index", Name: types.JSON{"zh_cn": "管理员新增", "en_us": "Admin Add"}, Router: true},
	//	{Key: "admin-edit", Parent: "admin-index", Name: types.JSON{"zh_cn": "管理员修改", "en_us": "Admin Edit"}, Router: true},
	//}
	//db.Create(&resourceData)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Policy Table'",
	//).AutoMigrate(&model.Policy{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//policyData := []model.Policy{
	//	{ResourceKey: "acl-index", AclKey: "acl", Policy: 1},
	//	{ResourceKey: "resource-index", AclKey: "resource", Policy: 1},
	//	{ResourceKey: "resource-index", AclKey: "policy", Policy: 1},
	//	{ResourceKey: "resource-index", AclKey: "acl", Policy: 0},
	//	{ResourceKey: "role-index", AclKey: "role", Policy: 1},
	//	{ResourceKey: "role-index", AclKey: "resource", Policy: 0},
	//	{ResourceKey: "admin-index", AclKey: "admin", Policy: 1},
	//	{ResourceKey: "admin-index", AclKey: "role", Policy: 0},
	//}
	//db.Create(&policyData)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Role Table'",
	//).AutoMigrate(&model.RoleBasic{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//roleData := model.RoleBasic{
	//	Key:  "*",
	//	Name: types.JSON{"zh_cn": "超级管理员", "en_us": "super"},
	//}
	//db.Create(&roleData)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Role Associated Resource Table'",
	//).AutoMigrate(&model.RoleResourceAssoc{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//roleResourceAssoc := []model.RoleResourceAssoc{
	//	{RoleKey: "*", ResourceKey: "system"},
	//	{RoleKey: "*", ResourceKey: "center"},
	//	{RoleKey: "*", ResourceKey: "profile"},
	//	{RoleKey: "*", ResourceKey: "acl-index"},
	//	{RoleKey: "*", ResourceKey: "acl-add"},
	//	{RoleKey: "*", ResourceKey: "acl-edit"},
	//	{RoleKey: "*", ResourceKey: "admin-index"},
	//	{RoleKey: "*", ResourceKey: "admin-add"},
	//	{RoleKey: "*", ResourceKey: "admin-edit"},
	//	{RoleKey: "*", ResourceKey: "resource-index"},
	//	{RoleKey: "*", ResourceKey: "resource-add"},
	//	{RoleKey: "*", ResourceKey: "resource-edit"},
	//	{RoleKey: "*", ResourceKey: "role-index"},
	//	{RoleKey: "*", ResourceKey: "role-add"},
	//	{RoleKey: "*", ResourceKey: "role-edit"},
	//}
	//db.Create(&roleResourceAssoc)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Admin Table'",
	//).AutoMigrate(&model.AdminBasic{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//pass, err := hash.Make("pass@VAN1234", hash.Option{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//adminData := model.AdminBasic{
	//	Username: "kain",
	//	Password: pass,
	//}
	//db.Create(&adminData)
	//err = db.Set(
	//	"gorm:table_options",
	//	"comment='Admin Associated Role Table'",
	//).AutoMigrate(&model.AdminRoleAssoc{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//adminRoleAssoc := model.AdminRoleAssoc{
	//	Username: "kain",
	//	RoleKey:  "*",
	//}
	//db.Create(&adminRoleAssoc)
	//var data []model.RolePolicy
	//db.Table("role_resource_assoc").
	//	Select("role_resource_assoc.role_key,policy.acl_key,max(policy.policy) as policy").
	//	Joins("join policy on policy.resource_key = role_resource_assoc.resource_key").
	//	Group("role_resource_assoc.role_key,policy.acl_key").
	//	Scan(&data)
	//log.Println(data)
}
