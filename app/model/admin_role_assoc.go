package model

type AdminRoleAssoc struct {
	ID         uint64
	Username   string     `gorm:"size:30;not null;uniqueIndex:role_unique"`
	RoleKey    string     `gorm:"size:200;not null;uniqueIndex:role_unique"`
	AdminBasic AdminBasic `gorm:"foreignKey:Username;references:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoleBasic  RoleBasic  `gorm:"foreignKey:RoleKey;references:Keyid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
