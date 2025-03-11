package services

import (
	"log"
	"sync"

	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type mqttServer struct {
	Server *server.Server
    Logger *log.Logger
}

var (
	defaultMqttServer *mqttServer
	once             = sync.Once{}
)

func GetMqttServer(logger *log.Logger) *mqttServer {
	if defaultMqttServer != nil {
		return defaultMqttServer
	}

	once.Do(func() {
		defaultMqttServer = initMqttServer(logger)
	})
	return defaultMqttServer
}

func initMqttServer(logger *log.Logger) *mqttServer {
	// 创建一个新的MQTT服务器实例
	s := server.New(nil)

	// 添加一个TCP监听器，监听本地的1883端口
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: ":1883",
	})
	if err := s.AddListener(tcp); err != nil {
		logger.Fatalf("Failed to add TCP listener: %v", err)
	}

	// 添加一个简单的认证钩子，这里使用匿名认证
	hook := &auth.AllowHook{}
	if err := s.AddHook(hook, nil); err != nil {
		logger.Fatalf("Failed to add authentication hook: %v", err)
	}

	// 启动MQTT服务器
	log.Println("Starting MQTT server on port 1883...")
	if err := s.Serve(); err != nil {
		logger.Fatalf("Failed to start MQTT server: %v", err)
	}
    return &mqttServer{
        Server: s,
        Logger: logger,
    }
}

func (m *mqttServer) Stop() {
    m.Server.Close()
}

func (m *mqttServer)Start() {
    
}