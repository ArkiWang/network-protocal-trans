package services

import (
	"context"
	"flag"
	"net/http"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/networkProtocalTrans/logger"
	"github.com/networkProtocalTrans/util"
)
var addr = flag.String("addr", "localhost:8080", "http service address")

// 添加配置结构
type WSConfig struct {
	Server     ServerConfig     `toml:"server"`
	Websocket  WebsocketConfig  `toml:"websocket"`
	Connection ConnectionConfig `toml:"connection"`
	CORS       CORSConfig       `toml:"cors"`
}

type ServerConfig struct {
	Name    string `toml:"name"`
	Address string `toml:"address"`
	Debug   bool   `toml:"debug"`
}

type WebsocketConfig struct {
	ReadBufferSize    int   `toml:"read_buffer_size"`
	WriteBufferSize   int   `toml:"write_buffer_size"`
	MaxMessageSize    int64 `toml:"max_message_size"`
	EnableCompression bool  `toml:"enable_compression"`
}

type ConnectionConfig struct {
	MaxConnections   int `toml:"max_connections"`
	HeartbeatTimeout int `toml:"heartbeat_timeout"`
	WriteTimeout     int `toml:"write_timeout"`
	ReadTimeout      int `toml:"read_timeout"`
}

type CORSConfig struct {
	AllowedOrigins []string `toml:"allowed_origins"`
	AllowAll       bool     `toml:"allow_all"`
}

// 修改 websocketServer 结构体
type websocketServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	logger     *logger.AppLogger
	config     *WSConfig
	maxClients int
}

// 创建websocket服务器实例
func GetWebsocketServer(ctx context.Context, log *logger.AppLogger, configPath string) *websocketServer {
	if defaultWSserver != nil {
		return defaultWSserver
	}

	wsOnce.Do(func() {
		// 加载配置
		var config WSConfig
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			log.LogFatal(ctx, "Failed to load websocket config", "error", err)
			return
		}

		// 更新 upgrader 配置
		upgrader = websocket.Upgrader{
			ReadBufferSize:    config.Websocket.ReadBufferSize,
			WriteBufferSize:   config.Websocket.WriteBufferSize,
			EnableCompression: config.Websocket.EnableCompression,
			CheckOrigin: func(r *http.Request) bool {
				if config.CORS.AllowAll {
					return true
				}
				origin := r.Header.Get("Origin")
				for _, allowed := range config.CORS.AllowedOrigins {
					if allowed == "*" || allowed == origin {
						return true
					}
				}
				return false
			},
		}

		defaultWSserver = &websocketServer{
			clients:    make(map[*websocket.Conn]bool),
			broadcast:  make(chan []byte),
			logger:     log,
			config:     &config,
			maxClients: config.Connection.MaxConnections,
		}
	})
	return defaultWSserver
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	defaultWSserver *websocketServer
	wsOnce          = sync.Once{}
)



// HandleConnections 处理websocket连接
// TODO add panic handler
func (s *websocketServer) HandleConnections(c *gin.Context) {
	ctx := c.Request.Context()
    w,r := c.Writer, c.Request
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        s.logger.LogErrorf(ctx,"websocketServer Upgrade failed: %v", err)
        return
    }
    defer ws.Close()

    s.clients[ws] = true

    for {
        messageType, message, err := ws.ReadMessage()
        if err != nil {
            s.logger.LogErrorf(ctx,"websocketServer ReadMessage failed: %v", err)
            delete(s.clients, ws)
            break
        }

        // 广播消息给所有客户端
        for client := range s.clients {
            if err := client.WriteMessage(messageType, message); err != nil {
                s.logger.LogErrorf(ctx,"websocketServer WriteMessage failed: %v", err)
                client.Close()
                delete(s.clients, client)
            }
        }
    }
}

func (s *websocketServer)Test(c *gin.Context){
	util.HomeTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/ws/echo")
}

