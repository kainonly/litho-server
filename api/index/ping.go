package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"time"
)

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	x.Csrf.SetToken(c)
	r := M{
		"name": x.V.Hostname,
		"now":  time.Now(),
	}
	if !x.V.IsRelease() {
		r["values"] = x.V
	}
	c.JSON(200, r)
}
