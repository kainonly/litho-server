package roles

import (
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	RolesService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}
