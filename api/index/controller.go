package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
}

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}
