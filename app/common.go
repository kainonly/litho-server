package app

import (
	"api/app/index"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	index.Provides,
	New,
)

func New(
	values *common.Values,
	index *index.Controller,
) *gin.Engine {
	r := middleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))
	return r
}
