package schedules

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

func (x *Controller) GetKeys(c *gin.Context) interface{} {
	keys, err := x.Service.GetKeys()
	if err != nil {
		return err
	}
	return keys
}
