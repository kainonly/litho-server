package model

type Admin struct {
	Common

	Username string `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Super    *bool  `gorm:"default:false" json:"-"`
	Name     string `gorm:"type:varchar(20)" json:"name"`
	Email    string `gorm:"type:varchar(255)" json:"email"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
}
