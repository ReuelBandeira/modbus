package mqtt

import (
	/// TODO: change canopen with your protocolname
	strPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"

	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MsgMQTTInput(msg mqtt.Message, client mqtt.Client) (strPkg.MessageInput, error) {

	var msgInf strPkg.MessageInput
	payload := msg.Payload()
	err := json.Unmarshal(payload, &msgInf)

	if err != nil {
		logPkg.CtsLog.Error("mqttInput:MsgMQTTInput:FAIL to decode JSON Input Msg MQTT - Data[%s] err[%v]\n",
			msg.Payload(), err)
		return msgInf, err
	}
	///logPkg.CtsLog.Debug("mqttInput:MsgMQTTInput: msgProtocol.Protocol    payload:\n%s\n\n", payload)
	return msgInf, err
}
