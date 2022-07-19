package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	AppService *Service
}

func (x *Controller) In(r *gin.RouterGroup) {
	r.GET("/", x.Index)
	r.POST("/auth", x.AuthLogin)
}

func (x *Controller) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time": x.AppService.Index(),
		"ip":   c.GetHeader("X-Forwarded-For"),
	})
}

func (x *Controller) AuthLogin(c *gin.Context) {
	var body struct {
		Identity string `json:"identity" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	//if err := x.AppService.Test(); err != nil {
	//	c.Error(err)
	//	return
	//}

	//c.SetCookie("access_token", ts, 0, "", "", true, true)
	//c.SetSameSite(http.SameSiteStrictMode)

	c.Status(http.StatusNoContent)
}
