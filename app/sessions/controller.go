package sessions

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

// Gets 获取会话
func (x *Controller) Gets(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	values, err := x.Service.Gets(ctx)
	if err != nil {
		return err
	}
	return values
}

// Delete 删除会话
func (x *Controller) Delete(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.Delete(ctx, uri.Id); err != nil {
		return err
	}
	return nil
}

// BulkDelete 删除所有会话
func (x *Controller) BulkDelete(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	if err := x.Service.BulkDelete(ctx); err != nil {
		return err
	}
	return nil
}
