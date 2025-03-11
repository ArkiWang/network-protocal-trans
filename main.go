package main

import (
	"context"
	"log"
	"path/filepath"

	"hvuit.com/mqtt2ws/services"
	"hvuit.com/mqtt2ws/start"
)

func main() {
    // 创建上下文
    ctx := context.Background()
    // 初始化日志
    configPath := filepath.Join("conf", "log.toml")
    if err := start.InitLogger(ctx, configPath); err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
	

    // 使用日志
    start.Logger.Println("应用启动")
    if start.AccessLogger != nil {
        start.AccessLogger.Println("访问日志记录")
    }

    // 初始化服务
    services.InitServices()
    // 启动mqtt服务
    services.MqttServer.Start()
}