package index

import "github.com/gin-gonic/gin"

func (x *Service) Version() interface{} {
	return gin.H{
		"version": "1.0",
	}
}
