package model

type Resource struct {
	ID             uint64 `json:"id"`
	Status         *bool  `gorm:"default:true" json:"status"`
	Parent         uint64 `gorm:"index:idx_parent_fragment;default:0;comment:父节点" json:"parent"`
	Fragment       string `gorm:"type:varchar(20);index:idx_parent_fragment,unique;not null;comment:URL片段" json:"fragment"`
	Name           string `gorm:"type:varchar(20);not null;comment:资源名称" json:"name"`
	Nav            *bool  `gorm:"default:false;comment:是否为导航" json:"nav"`
	Router         *bool  `gorm:"default:false;comment:是否为路由页面" json:"router"`
	Strategy       *bool  `gorm:"default:false;comment:策略节点，可绑定多个访问控制" json:"strategy"`
	Icon           string `gorm:"type:varchar(200);comment:导航节点的字体图标" json:"icon"`
	Sort           uint8  `gorm:"default:0;comment:导航节点排序" json:"sort"`
	ResourceAclRel []Acl  `gorm:"many2many:resource_acl_rel;references:Path;joinReferences:Path;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ResourceAclRel struct {
	ResourceID uint64 `gorm:"primaryKey"`
	Path       string `gorm:"type:varchar(20);primaryKey"`
	Mode       uint8  `gorm:"default:0;comment:只读或全部"`
}
