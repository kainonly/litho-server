package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
)

func Initialize(r *gin.Engine, s *controller.Controllers) {
	index := s.Index
	r.GET("/", bit.Bind(index.Index))
}
