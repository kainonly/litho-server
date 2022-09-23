package model

type User struct {
	Model
	Username string `gorm:"type:varchar;not null;uniqueIndex;comment:用户名" json:"username"`
	Password string `gorm:"type:varchar;not null;comment:密码" json:"password,omitempty"`
	Email    string `gorm:"type:varchar;not null;index;comment:电子邮件" json:"email"`
	Avatar   string `gorm:"type:varchar;not null;comment:头像" json:"avatar"`
	Status   *bool  `gorm:"default:true;comment:状态" json:"status"`
}
