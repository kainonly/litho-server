package index

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	data := M{
		"hostname": os.Getenv("HOSTNAME"),
		"endpoint": "litho-api",
		"ip":       string(c.GetHeader(x.V.Ip)),
		"now":      time.Now(),
	}
	c.JSON(200, data)
}
