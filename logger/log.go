package logger

import "context"

func InitServices(ctx context.Context) {
	// 初始化服务
	InitLogger(ctx, "conf/log.toml")
}