package values

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

// Get 获取动态配置
func (x *Controller) Get(c *gin.Context) interface{} {
	var query struct {
		Keys []string `form:"keys"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Service.Get(ctx, query.Keys)
	if err != nil {
		return err
	}
	return data
}

// Set 设置动态配置
func (x *Controller) Set(c *gin.Context) interface{} {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.Set(ctx, body); err != nil {
		return err
	}
	return nil
}

// Delete 删除动态配置
func (x *Controller) Delete(c *gin.Context) interface{} {
	var uri struct {
		Key string `uri:"key"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.Delete(ctx, uri.Key); err != nil {
		return err
	}
	return nil
}
