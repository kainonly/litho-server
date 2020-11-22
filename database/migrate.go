package main

import (
	"go.uber.org/fx"
	"taste-api/bootstrap"
	"taste-api/database/role_policy"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
		),
		fx.Invoke(
			//acl.Setup,
			//resource.Setup,
			//policy.Setup,
			//role_basic.Setup,
			//role_resource_rel.Setup,
			//admin_basic.Setup,
			//admin_role_rel.Setup,
			role_policy.Setup,
		),
	).Done()
	//db.Exec(fmt.Sprint(
	//	"CREATE VIEW IF NOT EXISTS `", prefix, "role` AS ",
	//	"SELECT ",
	//	"`", prefix, "role_basic`.`id`,",
	//	"`", prefix, "role_basic`.`key`,",
	//	"`", prefix, "role_basic`.`name`,",
	//	"group_concat(distinct `", prefix, "role_resource_assoc`.`resource_key` separator ',') AS `resource`,",
	//	"group_concat(distinct concat(`", prefix, "role_policy`.`acl_key`, ':', `", prefix, "role_policy`.`policy`) separator ',') AS `acl`,",
	//	"`", prefix, "role_basic`.`note`,",
	//	"`", prefix, "role_basic`.`status`,",
	//	"`", prefix, "role_basic`.`create_time`,",
	//	"`", prefix, "role_basic`.`update_time` ",
	//	"FROM `", prefix, "role_basic` ",
	//	"LEFT JOIN `", prefix, "role_resource_assoc` ON `", prefix, "role_resource_assoc`.`role_key` = `", prefix, "role_basic`.`key` ",
	//	"LEFT JOIN `", prefix, "role_policy` ON `", prefix, "role_policy`.`role_key` = `", prefix, "role_basic`.`key` ",
	//	"GROUP BY ",
	//	"`", prefix, "role_basic`.`id`,",
	//	"`", prefix, "role_basic`.`key`,",
	//	"`", prefix, "role_basic`.`name`,",
	//	"`", prefix, "role_basic`.`note`,",
	//	"`", prefix, "role_basic`.`status`,",
	//	"`", prefix, "role_basic`.`create_time`,",
	//	"`", prefix, "role_basic`.`update_time`;",
	//))
	//db.Exec(fmt.Sprint(
	//	"CREATE VIEW IF NOT EXISTS `", prefix, "admin` AS ",
	//	"SELECT ",
	//	"`", prefix, "admin_basic`.`id`,",
	//	"`", prefix, "admin_basic`.`username`,",
	//	"`", prefix, "admin_basic`.`password`,",
	//	"group_concat(distinct `", prefix, "admin_role_assoc`.`role_key` separator ',') AS `role`,",
	//	"`", prefix, "admin_basic`.`call`,",
	//	"`", prefix, "admin_basic`.`email`,",
	//	"`", prefix, "admin_basic`.`phone`,",
	//	"`", prefix, "admin_basic`.`avatar`,",
	//	"`", prefix, "admin_basic`.`status`,",
	//	"`", prefix, "admin_basic`.`create_time`,",
	//	"`", prefix, "admin_basic`.`update_time` ",
	//	"FROM `", prefix, "admin_basic` ",
	//	"JOIN `", prefix, "admin_role_assoc` ON `", prefix, "admin_role_assoc`.`username` = `", prefix, "admin_basic`.`username` ",
	//	"GROUP BY ",
	//	"`", prefix, "admin_basic`.`id`,",
	//	"`", prefix, "admin_basic`.`username`,",
	//	"`", prefix, "admin_basic`.`password`,",
	//	"`", prefix, "admin_basic`.`call`,",
	//	"`", prefix, "admin_basic`.`email`,",
	//	"`", prefix, "admin_basic`.`phone`,",
	//	"`", prefix, "admin_basic`.`avatar`,",
	//	"`", prefix, "admin_basic`.`status`,",
	//	"`", prefix, "admin_basic`.`create_time`,",
	//	"`", prefix, "admin_basic`.`update_time`;",
	//))
}
