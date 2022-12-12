package model

import (
	"github.com/lib/pq"
	"time"
)

type Project struct {
	ID         uint64         `json:"id"`
	Name       string         `gorm:"type:varchar;not null;comment:项目名称" json:"name"`
	Namespace  string         `gorm:"type:varchar;uniqueIndex;not null;comment:命名空间" json:"namespace"`
	Secret     string         `gorm:"type:varchar;not null;comment:密钥" json:"secret"`
	Entry      pq.StringArray `gorm:"type:varchar[];default:array[]::varchar[];comment:后端入口" json:"entry"`
	ExpireTime int64          `gorm:"default:0;comment:有效期" json:"expire_time"`
	Status     bool           `gorm:"default:true;not null;comment:状态" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
