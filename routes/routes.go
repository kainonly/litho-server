package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
)

func Initialize(
	r *gin.Engine,
	index *controller.Index,
) {
	r.GET("/", bit.Bind(index.Index))

	sys := r.Group("/sys")
	{
		sys.GET("/verify", bit.Bind(index.Verify))
	}
}
