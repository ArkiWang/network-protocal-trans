package services

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/networkProtocalTrans/logger"
)

type mqttServer struct {
	Server *server.Server
	Logger *logger.AppLogger
}

var (
	defaultMqttServer *mqttServer
	once              = sync.Once{}
)

func GetMqttServer(ctx context.Context, logger *logger.AppLogger) *mqttServer {
	if defaultMqttServer != nil {
		return defaultMqttServer
	}

	once.Do(func() {
		defaultMqttServer = initMqttServer(ctx, logger)
	})
	return defaultMqttServer
}


func initMqttServer(ctx context.Context, logger *logger.AppLogger) *mqttServer {
    // 创建一个新的MQTT服务器实例
    s := server.New(nil)

    // 添加一个TCP监听器，监听本地的1883端口
    tcp := listeners.NewTCP(listeners.Config{
        ID:      "t1",
        Address: ":1883",
    })
    if err := s.AddListener(tcp); err != nil {
        logger.LogFatal(ctx, "Failed to add TCP listener", "error", err)
    }

    // 添加一个简单的认证钩子，这里使用匿名认证
    hook := &auth.AllowHook{}
    if err := s.AddHook(hook, nil); err != nil {
        logger.LogFatal(ctx, "Failed to add authentication hook", "error", err)
    }

    mqttServer := &mqttServer{
        Server: s,
        Logger: logger,
    }

    // 设置优雅退出
    go mqttServer.gracefulShutdown(ctx)

    // 启动MQTT服务器
    logger.LogInfof(ctx, "Starting MQTT server on port 1883...")
    if err := s.Serve(); err != nil {
        logger.LogFatal(ctx, "Failed to start MQTT server", "error", err)
    }

    return mqttServer
}

func (m *mqttServer) gracefulShutdown(ctx context.Context) {
    // 监听系统信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // 等待信号
    sig := <-sigChan
    m.Logger.LogInfo(ctx, "Received shutdown signal", "signal", sig.String())

    // 关闭 MQTT 服务器
    m.Logger.LogInfo(ctx, "Shutting down MQTT server...")
    m.Stop()
    m.Logger.LogInfo(ctx, "MQTT server stopped")
}

func (m *mqttServer) Stop() {
    if err := m.Server.Close(); err != nil {
        m.Logger.LogError(context.Background(), "Error closing MQTT server", "error", err)
    }
}