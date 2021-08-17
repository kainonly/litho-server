package model

import (
	"github.com/google/uuid"
	"time"
)

type Admin struct {
	ID         uint64     `json:"id"`
	Status     *bool      `gorm:"default:true" json:"status"`
	CreateTime time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	UID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username   string     `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password   string     `gorm:"type:varchar(255);not null" json:"-"`
	Super      *bool      `gorm:"default:false" json:"-"`
	Role       []string   `gorm:"type:varchar(20)[]" json:"-"`
	Name       string     `gorm:"type:varchar(20)" json:"name"`
	Email      string     `gorm:"type:varchar(255)" json:"email"`
	Phone      string     `gorm:"type:varchar(20)" json:"phone"`
	Avatar     string     `gorm:"type:varchar(255)" json:"avatar"`
	Roles      []Role     `gorm:"many2many:admin_role_rel;References:Key;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Resources  []Resource `gorm:"many2many:admin_resource_rel;constraint:OnDelete:CASCADE"`
}
