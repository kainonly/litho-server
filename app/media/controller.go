package media

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

func (x *Controller) FindLabels(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	values, err := x.Service.FindLabels(ctx)
	if err != nil {
		return err
	}
	return values
}
