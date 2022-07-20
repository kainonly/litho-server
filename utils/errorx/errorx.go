package errorx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		switch any := err.Err.(type) {
		case validator.ValidationErrors:
			var message []string
			for _, field := range any {
				message = append(message, field.Error())
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": message,
			})
			break
		case Public:
			c.JSON(http.StatusBadRequest, any)
		default:
			c.Status(http.StatusInternalServerError)
		}
	}
}

type Public struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

func (x Public) Error() string {
	return x.Message
}

func NewPublic(code int, message string) Public {
	return Public{Code: code, Message: message}
}
