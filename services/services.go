package services

import (
	"context"

	"github.com/networkProtocalTrans/logger"
)

var (
	MqttServer *mqttServer
	WsServer   *websocketServer
)

func InitServices(ctx context.Context) {
	MqttServer = GetMqttServer(ctx, logger.DefaultLogger,"./conf/servicer/mqtt-test.toml")
	WsServer = GetWebsocketServer(ctx, logger.DefaultLogger,"./conf/servicer/web-socket-test.toml")
}
