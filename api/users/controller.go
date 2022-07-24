package users

import (
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	UsersService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}
