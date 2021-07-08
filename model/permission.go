package model

type Permission struct {
	ID         uint64 `json:"_id"`
	Key        string `json:"key"`
	Name       Object `json:"name"`
	Note       string `json:"note"`
	Status     bool   `json:"status"`
	CreateTime uint64 `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime uint64 `gorm:"autoUpdateTime" json:"update_time"`
}
