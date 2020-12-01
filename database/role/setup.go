package role

import (
	"gorm.io/gorm"
	"lab-api/config"
)

func Setup(db *gorm.DB, cfg *config.Config) error {
	opt := cfg.Database
	sql := "create view ? as "
	sql += "select rb.id,rb.key,rb.name,"
	sql += "json_agg(distinct rrr.resource_key) as resource,"
	sql += "json_agg(distinct concat(rp.acl_key, ':', rp.policy)) as acl,"
	sql += "rb.note,rb.status,rb.create_time,rb.update_time "
	sql += "from ? rb "
	sql += "left join ? rrr on rb.key = rrr.role_key "
	sql += "left join ? rp on rb.key = rp.role_key "
	sql += "group by rb.id, rb.key, rb.name, rb.note, rb.status, rb.create_time, rb.update_time"
	return db.Exec(sql,
		opt.Table("role"),
		opt.Table("role_basic"),
		opt.Table("role_resource_rel"),
		opt.Table("role_policy"),
	).Error
}
