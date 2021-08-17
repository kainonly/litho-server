package model

import "time"

type Resource struct {
	ID         uint64     `json:"id"`
	Status     *bool      `gorm:"default:true" json:"status"`
	CreateTime time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	Parent     uint64     `gorm:"default:0;comment:父节点" json:"parent"`
	Fragment   string     `gorm:"type:varchar(20);not null;comment:URL片段" json:"fragment"`
	Name       string     `gorm:"type:varchar(20);not null;comment:资源名称" json:"name"`
	Nav        *bool      `gorm:"default:false;comment:是否在导航中显示" json:"nav"`
	Router     *bool      `gorm:"default:false;comment:是否为路由页面" json:"router"`
	Icon       string     `gorm:"type:varchar(200);comment:导航字体图标" json:"icon"`
	Sort       uint8      `gorm:"default:0;comment:排序" json:"sort"`
	Acl        JSONObject `gorm:"type:jsonb;default:'{}';comment:访问控制" json:"-"`
}
