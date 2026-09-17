package pkg

import (
	"encoding/json"
	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func Sender(data interface{}, topics []string) bool {

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("ioMsgMqtt:Sender: Unable to marshal json message. [%v]",
			err)
		return false
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logger.Log.Errorf("ioMsgMqtt:Sender: Unable to connect with MQTT")
		return false
	} else {
		for _, topic := range topics {
			token := utilsPkg.MqttClient.Publish(topic, 0, false, jsonPayload)
			token.Wait()

			if token.Error() != nil {
				logger.Log.Errorf("ioMsgMqtt:Sender: Fail to send a MQTT message. [%v]",
					token.Error())
				return false
			}
		}

	}

	return true
}
