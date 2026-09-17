package config

import (
	protocol "canopen/protocol"
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Start protocol
func StartProtocol() {
	logPkg.CtsLog.Warn("StartProtocol: Starting.....")
	utilsPkg.CANopenTcpServer_CreateChannel()
	devices, err := GetDeviceConfigInfo()
	if err != nil {
		if os.IsNotExist(err) {
			// deviceConfig.json File do not exist
			logPkg.CtsLog.Warn("err[%v]\n", err)
		} else if devices != nil {
			logPkg.CtsLog.Error("StartProtocol:FAIL:GetDeviceConfigInfo()\n err[%s]\n", err)
		}
	} else {
		logPkg.CtsLog.Debug("StartProtocol:PASS:GetDeviceConfigInfo() Calling StartRead... DevSettings[devices]:\n%v\n\n", devices)
	}
	StartRead(devices)
}

func StopDevicesProtocol() {
	logPkg.CtsLog.Warn("StopDevicesProtocol: Stopping.....")
	close(utilsPkg.StopChan)
	utilsPkg.WaitGroup.Wait()
	logPkg.CtsLog.Warn("StopDevicesProtocol:Stopped!!!")
}

func RestartDevicesProtocol() {
	logPkg.CtsLog.Info("RestartDevicesProtocol Updating.....")
}

func UpdateDevicesProtocol() {
	logPkg.CtsLog.Info("UpdateDevicesProtocol: Updating.....")
}

// ==============================================
// GetDeviceConfigInfo: >>> deviceConfig.json <<<
// ==============================================
// [
//
//		{
//		  "devices": [   // var devInf []utilsPkg.Devices
//		    {
//		      "address": "canopen", // var itemInf []utilsPkg.DevSettings
//		      "port": "0",
//		      "name": "Device Name",
//		      "protocol": "canopen",
//		      "readingTime": 5000,
//		      "topics": [
//		        "general",
//		        "log"
//		      ],
//		      "data": [ // []map[string]interface{}
//		        {
//	 =========================================
//	 TODO  according with protocol_config.json
//	 =========================================
//		        } // END: "data" {}
//		      ] // END: "data"
//		    } // END: "devices" {}
//		  ] // END: "devices" []
//		}
//
// ]
func GetDeviceConfigInfo() ([]utilsPkg.DevSettings, error) {
	var devInf []utilsPkg.Devices
	var itemInf []utilsPkg.DevSettings
	var devConfigFileName string = "deviceConfig.json"
	var byafileData []byte = nil
	_, err := os.Stat(devConfigFileName)
	if os.IsNotExist(err) {
		// File do not exist
		logPkg.CtsLog.Debug("JSON File[%s]Not Exist\nerr[%v]\n", devConfigFileName, err)
		return nil, err
	} else {
		// File Exists
		byafileData, err = os.ReadFile(devConfigFileName)
		if err != nil {
			logPkg.CtsLog.Error("File[%s]Fail to read\nerr[%v]\n", devConfigFileName, err)
			return nil, err
		} else if len(byafileData) == 0 {
			logPkg.CtsLog.Debug("File[%s]Is empty\n", devConfigFileName)
			return nil, fmt.Errorf("file[%s] is empty", devConfigFileName)
		}
	}

	err = json.Unmarshal(byafileData, &devInf)
	if err != nil {
		logPkg.CtsLog.Error("File[%s]JSON Format Invalid\nerr[%v]\n devInf[%v]\n", devConfigFileName, err, devInf)
		return itemInf, err
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			itemInf = append(itemInf,
				utilsPkg.DevSettings{
					Address:     item.Address,
					Port:        item.Port,
					Name:        item.Name,
					Protocol:    item.Protocol,
					Id:          item.Id,
					ReadingTime: item.ReadingTime,
					Topics:      item.Topics,
					Data:        item.Data, // []map[string]interface{}
				})
		}
	}
	return itemInf, nil
}

func StartRead(devices []utilsPkg.DevSettings) {
	var emptyFile bool = false
	canBusPortAddress := fmt.Sprintf("%d", protocol.CanBusPortAddress)
	readHour := time.Now()
	readTime := readHour.Format("2006-01-02 15:04:05")
	canBusConnected := protocol.CanBusConnect(protocol.CanBusPortName, protocol.CanBusPortAddress)
	if canBusConnected {
		// CAN Bus is OK
		protocol.CanBusIsConnected = true
		logPkg.CtsLog.Debug("StartRead:PASS CanBusPortName[%s%d]CAN Bus is OK\n", protocol.CanBusPortName, protocol.CanBusPortAddress)
		protocol.MQTTSendGeneralDeviceStatus(devices[0], "canbus:connected")
		msgLog := utilsPkg.MqttLogMsg{
			Address:  protocol.CanBusPortName,
			Port:     canBusPortAddress,
			Name:     "CanBus",
			Protocol: "canopen",
			Id:       0,
			Time:     readTime,
			Message:  "connected",
			Code:     "OK",
			Source:   "readConfig:StartRead",
		}
		_ = protocol.MQTTSendLogMessage(msgLog, "info")
	} else {
		protocol.CanBusIsConnected = false
		// Sending error Log Can Bus Disconnected
		logPkg.CtsLog.Error("StartRead:FAIL CanBusPortName[%s%d]CAN Bus error\n", protocol.CanBusPortName, protocol.CanBusPortAddress)
		protocol.MQTTSendGeneralDeviceStatus(devices[0], "canbus:disconnected")
		msgLog := utilsPkg.MqttLogMsg{
			Address:  protocol.CanBusPortName,
			Port:     canBusPortAddress,
			Name:     "CanBus",
			Protocol: "canopen",
			Id:       0,
			Time:     readTime,
			Message:  "disconnected",
			Code:     "CANOPEN_402",
			Source:   "readConfig:StartRead",
		}
		_ = protocol.MQTTSendLogMessage(msgLog, "error")
	}
	if devices == nil {
		for {
			// Waiting for an existing and valid deviceConfig.json file
			newDevices, err := GetDeviceConfigInfo()
			if err != nil {
				if os.IsNotExist(err) {
					// File do not exist
					logPkg.CtsLog.Debug("err[%v]\n", err)
				} else if newDevices != nil {
					logPkg.CtsLog.Error("StartProtocol:FAIL:GetDeviceConfigInfo()\n err[%s]\n", err)
					return
				} else {
					if !emptyFile {
						logPkg.CtsLog.Warn("err[%v]\n", err)
						emptyFile = true
					}
				}
			} else {
				// Found a valid protocolConfig.json
				devices = newDevices
				logPkg.CtsLog.Warn("validDeviceConfig.json found!!!\n")
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	// Create Go Routines for all existing devices
	for _, device := range devices {
		utilsPkg.WaitGroup.Add(1)
		go func(device utilsPkg.DevSettings) {
			protocol.ContinuousRead(device, true) // Starting GoRoutine
		}(device)
		time.Sleep(1 * time.Second)
	}
	logPkg.CtsLog.Debug("StartRead:Exiting devices[%v]\n", devices)
}

// GoRotine:Channels Control
var (
	DoneChan      chan bool
	SyncOnceChan  sync.Once
	SyncWaitGroup sync.WaitGroup
)

func CANopenTcpServer_CreateGoRotineChannel() {
	DoneChan = make(chan bool)
}

func StartWrite(netConn net.Conn, devices []utilsPkg.DevSettings, client mqtt.Client, loop bool) {
	for _, dev := range devices {
		readHour := time.Now()
		readTime := readHour.Format("2006-01-02 15:04:05")
		if dev.Protocol == "canopen" {
			logPkg.CtsLog.Info("StartWrite: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
				dev.Name, dev.Protocol, dev.Address, dev.Port)
			netConn, err := protocol.CANopenTcpClient_ConnectToServer(netConn, dev.Address, dev.Port)
			if err == nil {
				logPkg.CtsLog.Info("StartWrite: Name[%s] Protocol[%s] Ip:Port[%s:%s]Connected netConn[%v]\n",
					dev.Name, dev.Protocol, dev.Address, dev.Port, netConn)
				msgLog := utilsPkg.MqttLogMsg{
					Address:  dev.Address,
					Port:     dev.Port,
					Name:     dev.Name,
					Protocol: dev.Protocol,
					Id:       dev.Id,
					Time:     readTime,
					Message:  "connected",
					Code:     "OK",
					Source:   "readConfig:StartWrite",
				}
				_ = protocol.MQTTSendLogMessage(msgLog, "info")
				protocol.ContinuousRead(dev, false) // Running Single GoRoutine
			} else {
				logPkg.CtsLog.Error("StartWrite: Name[%s] Protocol[%s] Ip:Port[%s:%s]Not sConnected err[%v]\n",
					dev.Name, dev.Protocol, dev.Address, dev.Port, err)
				msgLog := utilsPkg.MqttLogMsg{
					Address:  dev.Address,
					Port:     dev.Port,
					Name:     dev.Name,
					Protocol: dev.Protocol,
					Id:       dev.Id,
					Time:     readTime,
					Message:  "disconnected",
					Code:     "CANOPEN_403",
					Source:   "readConfig:StartWrite",
				}
				_ = protocol.MQTTSendLogMessage(msgLog, "error")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
