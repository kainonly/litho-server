package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"lab-api/app/api"
	"lab-api/app/system"
	"lab-api/bootstrap"
)

var Provides = wire.NewSet(
	bootstrap.HttpServer,
	bootstrap.InitializeDatabase,
	bootstrap.InitializeRedis,
	api.Provides,
	system.Provides,
	NewApp,
)

type App struct {
	*gin.Engine
}

func NewApp(
	r *gin.Engine,
	_ *api.Routes,
	_ *system.Routes,
) *App {
	return &App{Engine: r}
}
