package model

type RoleResource struct {
	ID          uint64
	RoleKey     string    `gorm:"size:200;not null;uniqueIndex:ram"`
	ResourceKey string    `gorm:"size:200;not null;uniqueIndex:ram"`
	RoleBasic   RoleBasic `gorm:"foreignKey:RoleKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resource    Resource  `gorm:"foreignKey:ResourceKey;references:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
