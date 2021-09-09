package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"lab-api/app/api"
	"lab-api/app/system"
	"lab-api/common"
)

var Provides = wire.NewSet(
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

func Run(_ *common.Set) (*App, error) {
	wire.Build(
		common.HttpServer,
		common.InitializeDatabase,
		common.InitializeRedis,
		common.InitializeCookie,
		common.InitializeAuthx,
		common.InitializeCipher,
		wire.Struct(new(common.Dependency), "*"),
		Provides,
	)
	return &App{}, nil
}
