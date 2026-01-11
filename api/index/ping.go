package index

import (
    "context"
    "os"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
    data := utils.H{
        "hostname": os.Getenv("HOSTNAME"),
        "endpoint": "litho-server",
        "now":      time.Now(),
    }

    c.JSON(200, data)
}
