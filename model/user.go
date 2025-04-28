package model

import (
	"server/common"
	"time"
)

type User struct {
	ID           string         `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	Status       *bool          `gorm:"column:status;type:boolean;not null;default:true" json:"status"`
	CreateTime   *time.Time     `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime   *time.Time     `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	DepartmentID string         `gorm:"column:department_id;type:bigint;not null;index" json:"department_id"`
	Pid          string         `gorm:"column:pid;type:bigint;not null;index" json:"pid"`
	RoleID       string         `gorm:"column:role_id;type:bigint;not null;index" json:"role_id"`
	Name         string         `gorm:"column:name;type:character varying;not null" json:"name"`
	Email        string         `gorm:"column:email;type:character varying;not null;uniqueIndex" json:"email"`
	Password     string         `gorm:"column:password;type:character varying;not null" json:"password"`
	Totp         string         `gorm:"column:totp;type:character varying;not null" json:"totp"`
	Avatar       string         `gorm:"column:avatar;type:character varying;not null" json:"avatar"`
	Sessions     int32          `gorm:"column:sessions;type:integer;not null" json:"sessions"`
	History      common.History `gorm:"column:history;type:jsonb;not null;default:'{}'" json:"history"`
}

func (User) TableName() string {
	return "user"
}
