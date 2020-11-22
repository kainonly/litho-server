package role_policy

import (
	"gorm.io/gorm"
	"taste-api/config"
)

func Setup(db *gorm.DB, cfg *config.Config) error {
	opt := cfg.Database
	sql := "create view ? as "
	sql += "select rrr.role_key, p.acl_key, max(p.policy) as policy "
	sql += "from ? rrr "
	sql += "join ? p on rrrresource_key = p.resource_key "
	sql += "group by rrr.role_key, p.acl_key"
	return db.Exec(sql,
		opt.Table("role_policy"),
		opt.Table("role_resource_rel"),
		opt.Table("policy"),
	).Error
}
