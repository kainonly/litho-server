package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/res"
	"log"
	"taste-api/application/common"
)

type Controller struct {
}

type LoginBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
}

func (c *Controller) Login(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body LoginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	result, err := app.Cache.AdminGet(body.Username)
	if err != nil {
		return res.Error(err)
	}
	log.Println(result)
	return res.Ok()
}

func (c *Controller) Verify() interface{} {
	return res.Ok()
}
