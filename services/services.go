package services

import "hvuit.com/mqtt2ws/start"

var (
	MqttServer *mqttServer
)

func InitServices() {
	MqttServer = GetMqttServer(start.Logger)
}
