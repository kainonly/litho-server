package model

type Permission struct {
	ID   uint64 `json:"id"`
	Code string `gorm:"type:varchar(50);not null;unique;comment:授权代码"`
	Name string `gorm:"type:varchar(20);not null;comment:授权名称"`
}
