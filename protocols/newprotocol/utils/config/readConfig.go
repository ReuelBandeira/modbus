package config

import (
	"encoding/json"
	"fmt"
	protocol "newprotocol/protocol"
	utilsPkg "newprotocol/utils"
	logPkg "newprotocol/utils/gologtofile"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Start protocol
func StartProtocol() {
	logPkg.CtsLog.Debug("StartProtocol: Starting.....")
	utilsPkg.CreateChannel()
	devices, err := GetDeviceConfigInfo()
	if err != nil {
		if os.IsNotExist(err) {
			// deviceConfig.json File do not exist
			logPkg.CtsLog.Warn("StartProtocol:WR:deviceConfig.json File do not exist\n err[%v]\n", err)
		} else if devices != nil {
			// deviceConfig.json File eist but no devices was found
			logPkg.CtsLog.Warn("StartProtocol:WR:GetDeviceConfigInfo()\n err[%s]\n", err)
		}
	} else {
		logPkg.CtsLog.Debug("StartProtocol:OK:GetDeviceConfigInfo() Calling StartRead... DevSettings[devices]:\n%v\n\n", devices)
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
//		      "address": "newprotocol", // var itemInf []utilsPkg.DevSettings
//		      "port": "0",
//		      "name": "Device Name",
//		      "protocol": "newprotocol",
//		      "readingTime": 5000,
//		      "topics": [
//		        "general",
//		        "log"
//		      ],
//		      "data": [ // []map[string]interface{}
//		        {
//	 =========================================
//	 TODO: according with protocol_config.json
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
		logPkg.CtsLog.Error("File[%s]JSON Format invalid\nerr[%v]\n devInf[%v]\n", devConfigFileName, err, devInf)
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

	for _, device := range devices {
		// Scanning Devices
		readHour := time.Now()
		readTime := readHour.Format("2006-01-02 15:04:05")
		statusDevice, conn := protocol.ProtocolConnect(device) //TODO:For CANopen replace device with CanPort
		if statusDevice {
			conn.Close()
			logPkg.CtsLog.Debug("StartRead:PASS Name[%s] Protocol[%s] Ip:Port[%s:%s] ONLINE\n",
				device.Name, device.Protocol, device.Address, device.Port)
			// Sending Info Log Device is Connected
			protocol.MQTTSendGeneralDeviceStatus(device, "connected")
			msgLog := utilsPkg.MqttLogMsg{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Id:       device.Id,
				Time:     readTime,
				Message:  "Device Connected & Running",
				Code:     "OK",
				Source:   "readConfig:StartRead",
			}
			_ = protocol.MQTTSendLogMessage(msgLog, "info")

			// Starting GoRoutine now
			utilsPkg.WaitGroup.Add(1)
			go func(device utilsPkg.DevSettings) {
				protocol.ContinuousRead(device, true) // Starting GoRoutine
			}(device)

		} else {
			//Fail to Connect with Device
			logPkg.CtsLog.Error("StartRead:FAIL to Connect\n\t DeviceName[%s]\n\t   Protocol[%s]\n\t    Ip:Port[%s:%s]\n",
				device.Name, device.Protocol, device.Address, device.Port)
			if utilsPkg.AnyDeviceIsConnected {

				utilsPkg.AnyDeviceIsConnected = false
				// Sending Status that Device is connected
				protocol.MQTTSendGeneralDeviceStatus(device, "diconnected")

				// Sending Error Log Device Disconnected
				msgLog := utilsPkg.MqttLogMsg{
					Address:  device.Address,
					Port:     device.Port,
					Name:     device.Name,
					Protocol: device.Protocol,
					Id:       device.Id,
					Time:     readTime,
					Message:  "disconnected",
					Code:     "200",
					Source:   "readConfig:StartRead",
				}
				_ = protocol.MQTTSendLogMessage(msgLog, "error")
			}

			// also starts GoRoutine with communication error for recovery
			utilsPkg.WaitGroup.Add(1)
			go func(device utilsPkg.DevSettings) {
				protocol.ContinuousRead(device, true) // Starting GoRoutine
			}(device)

		}
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

func CreateGoRotineChannel() {
	DoneChan = make(chan bool)
}

func StartWrite(devices []utilsPkg.DevSettings, client mqtt.Client, loop bool) {
	for _, device := range devices {
		// Scanning Devices
		readHour := time.Now()
		readTime := readHour.Format("2006-01-02 15:04:05")
		if device.Protocol == "newprotocol" {
			logPkg.CtsLog.Info("StartWrite: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
				device.Name, device.Protocol, device.Address, device.Port)
			statusDevice, conn := protocol.ProtocolConnect(device)
			if statusDevice {
				conn.Close()
				protocol.MQTTSendGeneralDeviceStatus(device, "connected")
				msgLog := utilsPkg.MqttLogMsg{
					Address:  device.Address,
					Port:     device.Port,
					Name:     device.Name,
					Protocol: device.Protocol,
					Id:       device.Id,
					Time:     readTime,
					Message:  "connected",
					Code:     "OK",
					Source:   "readConfig:StartWrite",
				}
				_ = protocol.MQTTSendLogMessage(msgLog, "info")
				protocol.ContinuousRead(device, false) // Running Single GoRoutine
			} else {
				protocol.MQTTSendGeneralDeviceStatus(device, "disconnected")
				msgLog := utilsPkg.MqttLogMsg{
					Address:  device.Address,
					Port:     device.Port,
					Name:     device.Name,
					Protocol: device.Protocol,
					Id:       device.Id,
					Time:     readTime,
					Message:  "disconnected",
					Code:     "201",
					Source:   "readConfig:StartWrite",
				}
				_ = protocol.MQTTSendLogMessage(msgLog, "error")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
