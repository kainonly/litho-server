package index

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	data := M{
		"hostname": os.Getenv("HOSTNAME"),
		"endpoint": "litho-server",
		"now":      time.Now(),
	}
	c.JSON(200, data)
}
