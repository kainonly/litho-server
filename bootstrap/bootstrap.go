package bootstrap

import (
	"context"
	"database/sql"
	"os"
	"regexp"
	"server/common"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
	"github.com/hertz-contrib/cors"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadStaticValues(path string) (v *common.Values, err error) {
	v = new(common.Values)
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &v); err != nil {
		return
	}
	return
}

func UseGorm(v *common.Values) (orm *gorm.DB, err error) {
	if orm, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  v.Database.Url,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		},
	); err != nil {
		return
	}
	var db *sql.DB
	if db, err = orm.DB(); err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return
}

func UseRedis(v *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(v.Database.Redis)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

func UsePassport(v *common.Values) *passport.Passport {
	return passport.New(
		passport.SetKey(v.Key),
		passport.SetIssuer(v.Domain),
	)
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
		csrf.SetDomain(v.Domain),
	)
}

func UseLocker(client *redis.Client) *locker.Locker {
	return locker.New(client)
}

func UseCaptcha(client *redis.Client) *captcha.Captcha {
	return captcha.New(client)
}

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	if v.Address == "" {
		return
	}
	vd := go_playground.NewValidator()
	vd.SetValidateTag("vd")
	vdx := vd.Engine().(*validator.Validate)
	vdx.RegisterValidation("snake", func(fl validator.FieldLevel) bool {
		matched, errX := regexp.MatchString("^[a-z_]+$", fl.Field().Interface().(string))
		if errX != nil {
			return false
		}
		return matched
	})
	vdx.RegisterValidation("sort", func(fl validator.FieldLevel) bool {
		matched, errX := regexp.MatchString("^[a-z_.]+:(-1|1)$", fl.Field().Interface().(string))
		if errX != nil {
			return false
		}
		return matched
	})

	opts := []config.Option{
		server.WithHostPorts(v.Address),
		server.WithCustomValidator(vd),
	}

	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		help.ErrorHandler(),
		cors.New(cors.Config{
			AllowOrigins: v.Cors.Origins,
			AllowMethods: []string{"GET", "POST"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type",
				"X-Page", "X-Pagesize", "X-Requested-With",
			},
			ExposeHeaders:    []string{"X-Total"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	return
}
