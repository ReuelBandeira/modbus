package ethip

import (
	"encoding/json"

	structPkg "ethernetip/utils"

	logPkg "ethernetip/utils/gologtofile"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func MQTTSendMessageAutomatic(msg structPkg.ExitPayloadMsg) {

	data := structPkg.MqttMsgStruct{
		Address: structPkg.Address{
			Address: msg.Address,
			Port:    msg.Port,
			Name:    msg.Name,
			Others:  msg.Others,
		},
		ReadTimeStamp: msg.ReadTimeStamp, //.Format(time.RFC3339),
		Protocol:      msg.Protocol,
		Data: structPkg.Data{
			BitMemories:  msg.BitMemories,
			WordMemories: msg.WordMemories,
		},
	}
	channelMqtt := "ethernetip"
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !structPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		token := structPkg.MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
		token.Wait()

		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
			logPkg.CtsLog.Error("%v", token.Error())
		}
	}
}

func SendStatusProtocol(status string) {
	var data structPkg.MessageStatusProtocol
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "ethernetip"
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !structPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		//for _, channelMqtt := range topics {
		//token := structPkg.MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
		token := structPkg.MqttClient.Publish("general", 0, false, jsonPayload)
		token.Wait()

		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
			logPkg.CtsLog.Error("%v", token.Error())
		}

		//}
	}
}

func SendStatusDevice(device structPkg.DevSettings, status string) {
	data := structPkg.MessageDeviceStatus{
		MessageType: "deviceStatus",
		Data: structPkg.MessageDataDeviceStatus{
			Protocol: "ethernetip",
			Device:   device.Name,
			Id:       device.Id,
			Status:   status,
		},
	}
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !structPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		//for _, channelMqtt := range topics {
		//token := structPkg.MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
		token := structPkg.MqttClient.Publish("general", 0, false, jsonPayload)
		token.Wait()

		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
			logPkg.CtsLog.Error("%v", token.Error())
		}

		//}
	}

}

func SendMsgLog(data interface{}, topic string) bool {

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		logPkg.CtsLog.Error("ioMsgMqtt:Sender: Unable to marshal json message. [%v]",
			err)
		return false
	}

	if !structPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:Sender: Unable to connect with MQTT")
		return false
	} else {

		token := structPkg.MqttClient.Publish(topic, 0, false, jsonPayload)
		token.Wait()

		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:Sender: Fail to send a MQTT message. [%v]",
				token.Error())
			return false
		}
	}

	return true
}
