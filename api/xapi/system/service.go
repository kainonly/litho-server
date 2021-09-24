package system

import (
	"github.com/gin-gonic/gin"
	"lab-api/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func (x *Service) Version() interface{} {
	return gin.H{
		"version": "1.0",
	}
}
