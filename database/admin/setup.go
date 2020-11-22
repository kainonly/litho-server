package admin

import (
	"gorm.io/gorm"
	"taste-api/config"
)

func Setup(db *gorm.DB, cfg *config.Config) error {
	opt := cfg.Database
	sql := "create view ? as "
	sql += "select ab.id, ab.username, ab.password,"
	sql += "array_agg(distinct arr.role_key) as role,"
	sql += "ab.call, ab.email, ab.phone, ab.avatar, ab.status, ab.create_time,ab.update_time "
	sql += "from ? ab "
	sql += "join ? arr on ab.username = arr.username "
	sql += "group by ab.id, ab.username, ab.password, ab.call, ab.email, ab.phone, ab.avatar, ab.status, ab.create_time,ab.update_time"
	return db.Exec(sql,
		opt.Table("admin"),
		opt.Table("admin_basic"),
		opt.Table("admin_role_rel"),
	).Error
}
