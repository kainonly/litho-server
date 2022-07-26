package locker

import (
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	SessionsService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}
