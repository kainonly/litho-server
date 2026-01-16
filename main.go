package main

import (
	"context"
	"os"
	"os/signal"
	"server/bootstrap"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	ctx := context.Background()

	// 加载配置
	values, err := bootstrap.LoadStaticValues("config/values.yml")
	if err != nil {
		panic(err)
	}

	// 使用 wire 生成 API
	api, err := bootstrap.NewAPI(values)
	if err != nil {
		panic(err)
	}

	// 初始化服务器
	var h *server.Hertz
	if h, err = api.Initialize(ctx); err != nil {
		panic(err)
	}

	// 启动服务器
	go h.Spin()

	// 处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err = h.Shutdown(ctx); err != nil {
		panic(err)
	}
}
