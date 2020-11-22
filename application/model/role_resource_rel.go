package model

type RoleResourceRel struct {
	ID          uint64
	RoleKey     string    `gorm:"size:200;not null;uniqueIndex:role_key_resource_key_unique"`
	ResourceKey string    `gorm:"size:200;not null;uniqueIndex:role_key_resource_key_unique"`
	RoleBasic   RoleBasic `gorm:"foreignKey:RoleKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource    Resource  `gorm:"foreignKey:ResourceKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
