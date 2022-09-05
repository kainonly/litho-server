package departments

import (
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	DepartmentsService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}
