package values

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/server/utils/helper"
	"net/http"
)

type Controller struct {
	ValuesService *Service
}

func (x *Controller) In(r *gin.RouterGroup) {
	r.GET("", x.Get)
	r.PATCH("", x.Set)
	r.DELETE(":key", x.Remove)
}

// Get 获取动态配置
func (x *Controller) Get(c *gin.Context) {
	var query struct {
		// 动态配置键
		Keys string `form:"keys"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	keys := helper.ParseArray(query.Keys)
	data := x.ValuesService.Get(keys...)

	c.JSON(http.StatusOK, data)
}

// Set 设置动态配置
func (x *Controller) Set(c *gin.Context) {
	var body map[string]interface{}
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	if err := x.ValuesService.Set(ctx, body); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Remove 移除动态配置
func (x *Controller) Remove(c *gin.Context) {
	var params struct {
		Key string `uri:"key" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	if err := x.ValuesService.Remove(ctx, params.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
