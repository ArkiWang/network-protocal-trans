package main

import (
	"context"
	"path/filepath"

	applog "github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/router"
	"github.com/networkProtocalTrans/services"
)

func main() {
	// 创建上下文
	ctx := context.Background()
	// 初始化日志
	configPath := filepath.Join("conf", "log.toml")

	log := applog.InitLogger(ctx, configPath)
	log.LogInfo(ctx, "init logger successfully")
	// log.LogError(ctx, "init logger failed")

	// 初始化服务
	services.InitServices(ctx)
	// 初始化路由
	r := router.InitRouter(ctx)
	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.LogFatal(ctx, "Server startup failed", "error", err)
		return
	}
	log.LogInfo(ctx, "Server started successfully")
}
