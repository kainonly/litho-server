package model

type Acl struct {
	ID         uint64     `json:"_id"`
	Key        string     `json:"key" binding:"required"`
	Name       JSONObject `json:"name" binding:"required"`
	Read       Array      `json:"read"`
	Write      Array      `json:"write"`
	Status     *bool      `json:"status"`
	CreateTime int64      `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime int64      `gorm:"autoUpdateTime" json:"update_time"`
}
