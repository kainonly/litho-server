package values

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	ValuesService *Service
}

func (x *Controller) Get(c *gin.Context) {}
