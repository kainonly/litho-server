package model

type Policy struct {
	ID          uint64
	ResourceKey string   `gorm:"size:200;not null;uniqueIndex:key_unique"`
	AclKey      string   `gorm:"size:200;not null;uniqueIndex:key_unique"`
	Policy      uint8    `gorm:"type:tinyint(1);not null;default:0;comment:0:readonly,1:read & write"`
	Acl         Acl      `gorm:"foreignKey:AclKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource    Resource `gorm:"foreignKey:ResourceKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
