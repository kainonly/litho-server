package validate

import (
	"github.com/go-playground/validator/v10"
)

type Message map[string]map[string]string

func Make(value interface{}, message Message) (errs map[string]string) {
	validate := validator.New()
	err := validate.Struct(value)
	if err != nil {
		errs = make(map[string]string)
		for _, ve := range err.(validator.ValidationErrors) {
			if message[ve.Field()] != nil && message[ve.Field()][ve.Tag()] != "" {
				errs[ve.Field()] = message[ve.Field()][ve.Tag()]
			} else {
				errs[ve.Field()] = ve.Error()
			}
		}
	}
	return
}
