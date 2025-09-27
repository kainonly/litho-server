package bootstrap

import (
	"os"
	"regexp"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passport"
	"gopkg.in/yaml.v3"
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
	h.Use(help.ErrorHandler())

	return
}
