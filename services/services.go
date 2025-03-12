package services

import (
	"context"

	"github.com/networkProtocalTrans/logger"
)
var (
	MqttServer *mqttServer
	appLogger = logger.DefaultLogger
)

func InitServices(ctx context.Context) {
	MqttServer = GetMqttServer(ctx, appLogger)
}
