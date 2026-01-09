package common

import (
	"github.com/kainonly/go/captcha"
	"github.com/kainonly/go/locker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Inject struct {
	fx.In

	V       *Values
	Db      *gorm.DB
	RDb     *redis.Client
	Captcha *captcha.Captcha
	Locker  *locker.Locker
}

type HandleFunc func(do *gorm.DB) *gorm.DB

type IAMUser struct {
	ID     string `json:"id"`
	OrgID  string `json:"org_id"`
	RoleID string `json:"role_id"`
	Active bool   `json:"active"`
	Ip     string `json:"-"`
}
