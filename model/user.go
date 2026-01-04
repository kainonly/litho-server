package model

import "time"

// User 用户
type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	Active    bool      `gorm:"column:active;not null;default:true"`         // 状态
	Email     string    `gorm:"column:email;type:text;not null;uniqueIndex"` // 电子邮件
	Phone     string    `gorm:"column:phone;type:text;not null;index"`       // 手机号
	Name      string    `gorm:"column:name;type:text;not null"`              // 姓名
	Password  string    `gorm:"column:password;type:text;not null"`          // 密码
	Avatar    string    `gorm:"column:avatar;type:text;not null"`            // 头像
}

func (User) TableName() string {
	return "user"
}
