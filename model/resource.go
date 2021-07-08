package model

type Resource struct {
	ID         uint64 `json:"_id"`
	Key        string `json:"key"`
	Parent     string `json:"parent"`
	Name       Object `json:"name"`
	Nav        bool   `json:"nav"`
	Router     bool   `json:"router"`
	Policy     bool   `json:"policy"`
	Icon       string `json:"icon"`
	Sort       uint   `json:"sort"`
	Status     bool   `json:"status"`
	CreateTime uint64 `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime uint64 `gorm:"autoUpdateTime" json:"update_time"`
}
