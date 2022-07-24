package pages

import (
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	PagesService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}
