package role_policy

import (
	"gorm.io/gorm"
	"taste-api/config"
)

func Setup(db *gorm.DB, cfg *config.Config) error {
	opt := cfg.Database
	sql := "create view ? as "
	sql += "select r.role_key, p.acl_key, max(p.policy) as policy "
	sql += "from ? r "
	sql += "join ? p on r.resource_key = p.resource_key "
	sql += "group by r.role_key, p.acl_key"
	return db.Exec(sql,
		opt.Table("role_policy"),
		opt.Table("role_resource_rel"),
		opt.Table("policy"),
	).Error
}
