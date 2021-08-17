package model

import "time"

type Resource struct {
	ID         uint64    `json:"id"`
	Status     *bool     `gorm:"default:true" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
	Parent     uint64    `gorm:"default:0;comment:父节点" json:"parent"`
	Fragment   string    `gorm:"type:varchar(20);not null;comment:URL片段" json:"fragment"`
	Name       string    `gorm:"type:varchar(20);not null;comment:资源名称" json:"name"`
	Nav        *bool     `gorm:"default:false;comment:是否在导航显示" json:"nav"`
	Router     *bool     `gorm:"default:false;comment:是否为路由页面" json:"router"`
	Policy     *bool     `gorm:"default:false;comment:是否为策略节点" json:"policy"`
	Icon       string    `gorm:"type:varchar(200);comment:导航字体图标" json:"icon"`
	Sort       uint8     `gorm:"default:0;comment:排序" json:"sort"`
	Acls       []Acl     `gorm:"many2many:policy;References:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Policy struct {
	ResourceID uint64 `gorm:"primaryKey"`
	AclKey     string `gorm:"type:varchar(20);primaryKey"`
	Act        uint8  `gorm:"default:0"`
}
