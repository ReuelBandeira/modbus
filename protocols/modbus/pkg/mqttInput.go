package pkg

import (
	"encoding/json"
	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MsgMQTTInput(msg mqtt.Message, client mqtt.Client) (utilsPkg.MessageInput, error) {

	var msgInput utilsPkg.MessageInput
	err := json.Unmarshal(msg.Payload(), &msgInput)

	if err != nil {
		logger.Log.Errorf("mqttInput:MsgMQTTInput: Unable to decode JSON Input Msg MQTT - Data[%s] err[%v]",
			msg.Payload(), err)
		return msgInput, err
	}
	return msgInput, err
}
