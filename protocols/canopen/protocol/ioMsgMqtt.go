package canopen

import (
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"os"
	"time"

	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/go-cmp/cmp" /// ADDED 14/10/2022 FROOM HOLDEN
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func MQTTSendMessageAutomatic(msg utilsPkg.NewDeviceExitPayloadJsonMsg) error {

	data := utilsPkg.MqttMsgStruct{
		Address: utilsPkg.Address{
			Address: msg.Address,
			Port:    msg.Port,
			Name:    msg.Name,
			Others:  msg.Others,
		},
		ReadTimeStamp: msg.ReadTimeStamp, //.Format(time.RFC3339),
		Protocol:      msg.Protocol,
		Id:            msg.Id,
		JsonDeviceData: utilsPkg.JsonDeviceData{
			SDO:  msg.SDO,
			PDO:  msg.PDO,
			NMT:  msg.NMT,
			SYNC: msg.SYNC,
			TIME: msg.TIME,
			EMCY: msg.EMCY,
		},
	}

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic:FAIL:err[%v]Fail to decode JSON data\ndata[%s]\n", err, data)
		return err
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: MQTT Client Disconencted from MQTT Brocker")
		return fmt.Errorf("MQTT Client Disconnneced from MQTT Brocker")
	} else {
		for _, topicMqtt := range msg.Topics {
			token := utilsPkg.MqttClient.Publish(topicMqtt, 0, false, jsonPayload)
			token.Wait()
			if token.Error() != nil {
				logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL send msg to MQTT Brocker topic[%s] err[%v]\n", topicMqtt, token.Error())
			} else {
				logPkg.CtsLog.Debug("ioMsgMqtt:MQTTSendMessageAutomatic: PASS sent msg to MQTT Brocker topic[%s]jsonPayload:\n%s\n\n", topicMqtt, jsonPayload)
			}
		}
	}
	return err
}

// msgType: MessageType = deviceStatus, status, error, warning, info, debug, etc...
func MQTTSendLogMessage(msg utilsPkg.MqttLogMsg, msgType string) bool {

	MqttLogMsg := utilsPkg.MqttLogMsg{
		Address:  msg.Address,
		Port:     msg.Port,
		Name:     msg.Name,
		Protocol: msg.Protocol,
		Id:       msg.Id,
		Time:     msg.Time,
		Message:  msg.Message,
		Code:     msg.Code,
		Source:   msg.Source,
	}

	msgLog := utilsPkg.MessageDeviceLog{
		MessageType: msgType, // deviceStatus, status, error, warning, info, debug, etc...
		Data:        MqttLogMsg,
	}

	jsonPayload, err := json.Marshal(msgLog)
	if err != nil {
		logPkg.CtsLog.Error("IoMsgMqtt:MQTTSendLogMessage:FAIL topic[log] msgType[%s] json.Marshal(msg[%v])\n", msgType, msg)
		return false
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("IoMsgMqtt:MQTTSendLogMessage:FAIL MQTT Publish topic[log]Disconnected!! msgType[%s]\njsonPayload[%s]\n", msgType, jsonPayload)
		return false
	} else {

		// Publishing on all Mqtt Channels set to CanOpen
		MqttToken = utilsPkg.MqttClient.Publish("log", 1, false, jsonPayload)
		if MqttToken.Error() != nil {
			logPkg.CtsLog.Error("IoMsgMqtt:MQTTSendLogMessage:FAIL to MQTT Publish topic[log] msgType[%s] err[%v] \n", msgType, MqttToken.Error())
			return false
		}
		logPkg.CtsLog.Debug("IoMsgMqtt:MQTTSendLogMessage:PASS to MQTT Publish topic[log] msgType[%s]\njsonPayload[%s] \n", msgType, jsonPayload)
	}
	return true
}

func MQTTSendProtocolStatus(status string) {
	var data utilsPkg.MessageStatusProtocol
	var topic string = "general"
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "canopen" ///TODO change to your protocol name
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:SendStatusProtocol: FAIL to connect with MQTT")
	} else {
		token := utilsPkg.MqttClient.Publish(topic, 0, false, jsonPayload)
		token.Wait()
		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:SendStatusProtocol: FAIL to send a MQTT Brocker topic[%s] err[%v]\n", topic, token.Error())
		}
	}
}

func MQTTSendGeneralDeviceStatus(dev utilsPkg.DevSettings, status string) {
	data := utilsPkg.MessageDeviceStatus{
		MessageType: "deviceStatus",
		Data: utilsPkg.MessageDataDeviceStatus{
			Protocol: "canopen",
			Id:       dev.Id,
			Device:   dev.Name,
			Status:   status,
		},
	}
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendGeneralDeviceStatus:FAIL to connect with MQTT")
	} else {
		token := utilsPkg.MqttClient.Publish("general", 0, false, jsonPayload)
		token.Wait()

		if token.Error() != nil {
			logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendGeneralDeviceStatus:FAIL to send a MQTT Brocker topic[general] err[%v]\n", token.Error())
		}
	}
}

// 14/10/2023 BEGIN FROM MODBUS HOLDEN - NOT USED
func Sender(data interface{}, topics []string) bool {

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		logPkg.CtsLog.Error("ioMsgMqtt:Sender: Unable to marshal json message. [%v]",
			err)
		return false
	}

	if !utilsPkg.MqttClient.IsConnected() {
		logPkg.CtsLog.Error("ioMsgMqtt:Sender: Unable to connect with MQTT")
		return false
	} else {
		for _, topic := range topics {
			token := utilsPkg.MqttClient.Publish(topic, 0, false, jsonPayload)
			token.Wait()

			if token.Error() != nil {
				logPkg.CtsLog.Error("ioMsgMqtt:Sender: Fail to send a MQTT message. [%v]",
					token.Error())
				return false
			}
		}

	}

	return true
}

func SendStatusProtocol(status string) {
	logPkg.CtsLog.Debug("sendMsg:SendStatusProtocol: Send status Protocol to MQTT: [%s]", status)
	var data utilsPkg.MessageStatusProtocol
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "modbus"
	Sender(data, utilsPkg.General)
}

func SendStatusDevice(device utilsPkg.Device, status string) {
	logPkg.CtsLog.Error("sendMsg:SendStatusDevice: Send status Device [%v] to MQTT: [%s]", device.Address, status)

	data := utilsPkg.MessageDeviceStatus{
		MessageType: "deviceStatus",
		Data: utilsPkg.MessageDataDeviceStatus{
			Protocol: device.Protocol,
			Device:   device.Name,
			Status:   status,
		},
	}

	Sender(data, utilsPkg.General)
}

// Slice to store detected changes
var changes []utilsPkg.DeviceChange

func StartReadFile() {
	//utilsPkg.CANopenTcpServer_CreateChannel()
	//Initialize the originalData with the first read of JSON data
	originalData, err := readJSON(utilsPkg.FilePath)
	if err != nil {
		logPkg.CtsLog.Error("checkUpDevice:StartCheckDevice: Fail to read JSON File: [%v]", err)

		dataLog := utilsPkg.DataLog{
			Address:  "",
			Port:     "",
			Name:     "",
			Protocol: "",
			Time:     time.Now().Format("2006-01-02 15:04:05"),
			Message:  fmt.Sprintf("Fail to read JSON File. [%v]", err),
			Code:     "",
			Source:   "checkUpdDevice",
		}

		msgLog := utilsPkg.MessageLog{
			MessageType: "error",
			Data:        dataLog,
		}

		Sender(msgLog, utilsPkg.Log)
		return
	}
	logPkg.CtsLog.Debug("checkUpDevice:StartCheckDevice: Json file reading done.")

	// CANopenTcpServer_Create a map to keep track of the originalData devices by ID
	utilsPkg.OrigiNalDataID = make(map[rune]utilsPkg.Device)
	for _, device := range originalData.Device {
		utilsPkg.OrigiNalDataID[device.ID] = device
	}
}

func readJSON(filePath string) (utilsPkg.DeviceStruct, error) {
	// Read the JSON file
	var devices []utilsPkg.DeviceStruct
	jsonData, err := os.ReadFile(utilsPkg.FilePath)
	if err != nil {
		logPkg.CtsLog.Error("checkUpDevice:readJSON: Unable to read a new JSON device file. Error [%v]",
			err)
		return devices[0], err
	}

	err = json.Unmarshal(jsonData, &devices)
	if err != nil {
		logPkg.CtsLog.Error("checkUpDevice:readJSON: Unable to Unmarshal a new Json device file. Error [%v]",
			err)
		return utilsPkg.DeviceStruct{}, err
	}

	return devices[0], nil
}

func DiffDevices() []utilsPkg.DeviceChange {
	//Read New File Version
	updatedData, err := readJSON(utilsPkg.FilePath)
	if err != nil {
		logPkg.CtsLog.Error("CheckUpdDevice:DiffDevices: Error reading new JSON file. Error [%v]", err)
		return nil
	}

	// Reset changes for each iteration
	changes = nil

	//Check New and Updated devices
	for _, updatedDevice := range updatedData.Device {
		originalDevice, ok := utilsPkg.OrigiNalDataID[updatedDevice.ID]
		if !ok {
			logPkg.CtsLog.Error("CheckUpdDevice:DiffDevices: New device found. ID [%v] originalDevice[%v]", updatedDevice.ID, originalDevice)

			dataLog := utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("New device found. ID [%v]", updatedDevice.ID),
				Code:     "",
				Source:   "checkUpdDevice",
			}

			msgLog := utilsPkg.MessageLog{
				MessageType: "error",
				Data:        dataLog,
			}
			Sender(msgLog, utilsPkg.Log)

			changes = append(changes, utilsPkg.DeviceChange{Device: updatedDevice, Type: utilsPkg.New})
			continue
		}

		if !cmp.Equal(originalDevice, updatedDevice) {
			changes = append(changes, utilsPkg.DeviceChange{Device: updatedDevice, Type: utilsPkg.Updated})
			logPkg.CtsLog.Error("CheckUpdDevice:DiffDevices: Update on device [%v] found.", updatedDevice.ID)
		}
	}

	// Find deleted devices
	for id, originalDevice := range utilsPkg.OrigiNalDataID {
		found := false
		for _, updatedDevice := range updatedData.Device {
			if id == updatedDevice.ID {
				found = true
				break
			}
		}
		if !found {
			changes = append(changes, utilsPkg.DeviceChange{Device: originalDevice, Type: utilsPkg.Deleted})
			logPkg.CtsLog.Error("CheckUpdDevice:DiffDevices: Device deleted [%v] from list.", originalDevice.ID)
		}
	}

	// Update originalDevicesByID for the next comparison
	utilsPkg.OrigiNalDataID = make(map[rune]utilsPkg.Device)
	for _, updatedDevice := range updatedData.Device {
		utilsPkg.OrigiNalDataID[updatedDevice.ID] = updatedDevice
	}

	return changes
}

// 14/10/2023   END FROM MODBUS HOLDEN - NOTUSED
