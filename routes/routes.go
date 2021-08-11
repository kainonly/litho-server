package routes

import (
	"github.com/gin-gonic/gin"
	"lab-api/controller"
)

func Initialize(r *gin.Engine, s *controller.Controllers) {
	index := s.Index
	r.GET("/", index.Index)
}
