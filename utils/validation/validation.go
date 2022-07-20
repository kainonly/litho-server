package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

func Extend() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("objectId", func(fl validator.FieldLevel) bool {
			return primitive.IsValidObjectID(fl.Field().String())
		})
		v.RegisterValidation("key", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[a-z_]+$`, fl.Field().String())
			return matched
		})
		v.RegisterValidation("sort", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[a-z_]+\.(1|-1)$`, fl.Field().String())
			return matched
		})
	}
}
