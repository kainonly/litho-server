package model

import (
	"time"
)

type Department struct {
	ID          uint64    `json:"id"`
	Parent      uint64    `gorm:"default:0;not null;comment:父节点" json:"parent"`
	Name        string    `gorm:"type:varchar;not null;comment:名称" json:"name"`
	Description string    `gorm:"type:varchar;comment:描述" json:"description"`
	Sort        int64     `gorm:"default:0;not null;comment:排序" json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
