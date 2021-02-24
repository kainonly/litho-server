package common

import (
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/typ"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/application/redis"
	"lab-api/config"
	"net/http"
	"strings"
)

type Dependency struct {
	fx.In

	Config *config.Config
	Db     *gorm.DB
	Redis  *redis.Model
	Curd   *curd.Curd
}

func (c *Dependency) Inject(dependency interface{}) {
	dep := dependency.(Dependency)

	c.Config = dep.Config
	c.Db = dep.Db
	c.Redis = dep.Redis
	c.Curd = dep.Curd
}

var (
	SystemCookie = typ.Cookie{
		Name:     "system",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
)

func StringToSlice(value string, sep string) []string {
	if value == "" {
		return []string{}
	}
	return strings.Split(value, sep)
}
