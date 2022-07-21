package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/weplanx/server/utils/errs"
	"io"
	"strings"
)

func BindAndValidate(r io.Reader, obj interface{}, validate string) (err error) {
	if err = json.NewDecoder(r).Decode(obj); err != nil {
		return errs.NewPublic(0, "data type must be JSON")
	}
	if err = validator.New().Var(obj, validate); err != nil {
		return
	}
	return
}

func ParseArray(v string) (data []string) {
	if len(v) == 0 {
		return
	}
	return strings.Split(v, ",")
}
