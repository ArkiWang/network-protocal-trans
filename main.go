package main

import (
	"context"
	"path/filepath"

	applog "github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/services"
)

func main() {
    // 创建上下文
    ctx := context.Background()
    // 初始化日志
    configPath := filepath.Join("conf", "log.toml")

    applog.InitLogger(ctx, configPath)
    
    // 初始化服务
    services.InitServices(ctx)
}