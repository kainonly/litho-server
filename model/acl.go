package model

type Acl struct {
	ID         uint64 `json:"_id"`
	Key        string `gorm:"type:varchar(200);unique" json:"key"`
	Name       Object `gorm:"type:json" json:"name"`
	Read       string `gorm:"type:longtext" json:"read"`
	Write      string `gorm:"type:longtext" json:"write"`
	Status     bool   `gorm:"default:1" json:"status"`
	CreateTime uint64 `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime uint64 `gorm:"autoUpdateTime" json:"update_time"`
}
