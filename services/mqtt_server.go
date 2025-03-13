package services

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/util"
)

// 添加配置结构
type MQTTConfig struct {
	Server     ServerConfig     `toml:"server"`
	Connection ConnectionConfig `toml:"connection"`
	TCP        TCPConfig        `toml:"tcp"`
	Auth       AuthConfig       `toml:"auth"`
	MQTT       MQTTConfigDetail `toml:"mqtt"`
}

type TCPConfig struct {
	ID             string `toml:"id"`
	MaxConnections int    `toml:"max_connections"`
	BufferSize     int    `toml:"buffer_size"`
}

type AuthConfig struct {
	Type  string   `toml:"type"`
	Basic BaseAuth `toml:"basic"`
}

type BaseAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type MQTTConfigDetail struct {
	MaxMessageSize int  `toml:"max_message_size"`
	KeepAlive      int  `toml:"keep_alive"`
	SessionExpiry  int  `toml:"session_expiry"`
	AllowAnonymous bool `toml:"allow_anonymous"`
}

// 修改 mqttServer 结构体
type mqttServer struct {
	Server *server.Server
	Logger *logger.AppLogger
	config *MQTTConfig
}

func GetMqttServer(ctx context.Context, logger *logger.AppLogger, configPath string) *mqttServer {
	if defaultMqttServer != nil {
		return defaultMqttServer
	}

	mqttOnce.Do(func() {

		// 加载配置
		var config MQTTConfig
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			logger.LogFatal(ctx, "Failed to load MQTT config", "error", err)
			return
		}

		// 创建服务器实例
		s := server.New(&server.Options{
			Capabilities: &server.Capabilities{},
		})

		// 配置TCP监听器
		tcp := listeners.NewTCP(listeners.Config{
			ID:      config.TCP.ID,
			Address: config.Server.Address,
		})

		if err := s.AddListener(tcp); err != nil {
			logger.LogFatal(ctx, "Failed to add TCP listener", "error", err)
		}

		// 配置认证
		var hook server.Hook
		switch config.Auth.Type {
		case "basic":

		default:
			hook = &auth.AllowHook{}
		}

		if err := s.AddHook(hook, nil); err != nil {
			logger.LogFatal(ctx, "Failed to add authentication hook", "error", err)
		}

		mqttServer := &mqttServer{
			Server: s,
			Logger: logger,
			config: &config,
		}

		// 设置优雅退出
		go mqttServer.gracefulShutdown(ctx)

		// 启动MQTT服务器
		logger.LogInfo(ctx, "Starting MQTT server",
			"name", config.Server.Name,
			"address", config.Server.Address,
		)

		util.SafeGo(ctx, func() {
			if err := s.Serve(); err!= nil {
				logger.LogFatal(ctx, "Failed to start MQTT server", "error", err)
			}
		})
		

		// 关闭服务端时需要做的一些清理工作
		defaultMqttServer = mqttServer
	})
	return defaultMqttServer
}

var (
	defaultMqttServer *mqttServer
	mqttOnce          = sync.Once{}
)


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
