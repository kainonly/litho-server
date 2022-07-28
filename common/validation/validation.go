package validation

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Extend() {
	binding.MustRegValidateFunc("mongoId", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		if !primitive.IsValidObjectID(args[0].(string)) {
			return primitive.ErrInvalidHex
		}
		return nil
	})
}
