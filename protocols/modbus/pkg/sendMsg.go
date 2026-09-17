package pkg

import (
	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"
)

func SendStatusProtocol(status string) {
	logger.Log.Debugf("sendMsg:SendStatusProtocol: Send status Protocol to MQTT: [%s]", status)
	var data utilsPkg.MessageStatusProtocol
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "modbus"
	Sender(data, utilsPkg.General)
}

func SendStatusDevice(device utilsPkg.Device, status string) {
	logger.Log.Debugf("sendMsg:SendStatusDevice: Send status Device [%v] to MQTT: [%s]", device.Address, status)

	data := utilsPkg.MessageDeviceStatus{
		MessageType: "deviceStatus",
		Data: utilsPkg.MessageDataDeviceStatus{
			Protocol: device.Protocol,
			Device:   device.Name,
			Status:   status,
			ID: device.ID,
		},
	}

	Sender(data, utilsPkg.General)
}
