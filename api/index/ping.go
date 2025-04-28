package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"os"
	"time"
)

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	data := M{
		"hostname": os.Getenv("HOSTNAME"),
		"endpoint": "flower-admin-api",
		"ip":       string(c.GetHeader(x.V.Ip)),
		"now":      time.Now(),
	}
	c.JSON(200, data)
}
