package model

import (
	"github.com/google/uuid"
	"time"
)

type Admin struct {
	ID          uint64     `json:"id"`
	Status      *bool      `gorm:"default:true" json:"status"`
	CreateTime  time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	UUID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();comment:用户唯一标识"`
	Username    string     `gorm:"type:varchar(20);not null;unique;comment:用户名" json:"username"`
	Password    string     `gorm:"type:varchar(255);not null;comment:密码，建议采用argon2id" json:"-"`
	Name        string     `gorm:"type:varchar(20);comment:称呼" json:"name"`
	Email       string     `gorm:"type:varchar(255);comment:电子邮件" json:"email"`
	Phone       string     `gorm:"type:varchar(20);comment:联系电话" json:"phone"`
	Avatar      string     `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	Permissions Array      `gorm:"type:json;comment:特殊授权"`
	Roles       []Role     `gorm:"many2many:admin_role_rel;constraint:OnDelete:CASCADE"`
	Resources   []Resource `gorm:"many2many:admin_resource_rel;constraint:OnDelete:CASCADE"`
}
