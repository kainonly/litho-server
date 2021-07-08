package model

type Acl struct {
	ID         uint64 `json:"_id"`
	Key        string `json:"key"`
	Name       Object `json:"name"`
	Read       string `json:"read"`
	Write      string `json:"write"`
	Status     bool   `json:"status"`
	CreateTime int64  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime int64  `gorm:"autoUpdateTime" json:"update_time"`
}
