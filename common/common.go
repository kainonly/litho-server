package common

import (
	"github.com/kainonly/go/captcha"
	"github.com/kainonly/go/locker"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Inject struct {
	V       *Values
	Db      *gorm.DB
	RDb     *redis.Client
	Captcha *captcha.Captcha
	Locker  *locker.Locker
}

type HandleFunc func(do *gorm.DB) *gorm.DB

type IAMUser struct {
	ID     string `json:"id"`
	Status *bool  `json:"status"`
}
