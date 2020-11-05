package model

type RoleResourceAssoc struct {
	ID          uint64
	RoleKey     string    `gorm:"size:200;not null;uniqueIndex:key_unique"`
	ResourceKey string    `gorm:"size:200;not null;uniqueIndex:key_unique"`
	RoleBasic   RoleBasic `gorm:"foreignKey:RoleKey;references:Keyid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource    Resource  `gorm:"foreignKey:ResourceKey;references:Keyid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
