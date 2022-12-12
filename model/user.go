package model

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID         uint64        `json:"id"`
	Username   string        `gorm:"type:varchar;not null;uniqueIndex;comment:用户名" json:"username"`
	Password   string        `gorm:"type:varchar;not null;comment:密码" json:"-"`
	Email      string        `gorm:"type:varchar;not null;index;comment:电子邮件" json:"email"`
	Name       string        `gorm:"type:varchar;default:'';not null;comment:称呼" json:"name"`
	Avatar     string        `gorm:"type:varchar;default:'';not null;comment:头像" json:"avatar"`
	Department uint64        `gorm:"default:0;not null;comment:所属部门" json:"department,omitempty"`
	Roles      pq.Int64Array `gorm:"type:int8[];default:array[]::int8[];not null;comment:权限组" json:"roles,omitempty"`
	Status     bool          `gorm:"default:true;not null;comment:状态" json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
