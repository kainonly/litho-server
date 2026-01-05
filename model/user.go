package model

import "time"

// User 用户表
type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Active    bool      `gorm:"column:active;not null;default:true;comment:状态"`           // 状态
	Email     string    `gorm:"column:email;type:text;not null;uniqueIndex;comment:电子邮件"` // 电子邮件
	Phone     string    `gorm:"column:phone;type:text;not null;index;comment:手机号"`        // 手机号
	Name      string    `gorm:"column:name;type:text;not null;comment:姓名"`                // 姓名
	Password  string    `gorm:"column:password;type:text;not null;comment:密码"`            // 密码
	Avatar    string    `gorm:"column:avatar;type:text;not null;comment:头像"`              // 头像
}

func (User) TableName() string {
	return "user"
}
