package configs

import (
	"encoding/json"
	//"fmt"
	"os"
	"sync"
	"time"

	"icts/modbus/configs/logger"
	packages "icts/modbus/pkg"
	utilsPkg "icts/modbus/utils"
)

func StartModbus() {
	logger.Log.Debug("modbusInit:StartModbus: Set channel to control go routines.")
	utilsPkg.CreateChannel()

	//Read a file from root
	_, err := GetDevConfig()
	if err != nil {
		logger.Log.Errorf("modbusInit:StartModbus: Unable to capture data from deviceConfig.json file. Error [%s]",
			err)
		return
	}
}

func GetDevConfig() ([]utilsPkg.Device, error) {

	var deviceStruct []utilsPkg.DeviceStruct
	file, err := os.ReadFile(utilsPkg.FilePath)
	if err != nil {
		logger.Log.Errorf("modbusInit:GetDevConfig: Fail to read JSON File: [%v]", err)
		return nil, err
	}
	logger.Log.Debug("modbusInit:GetDevConfig: Json file reading done.")

	err = json.Unmarshal(file, &deviceStruct)
	if err != nil {
		logger.Log.Errorf("modbusInit:GetDevConfig: Fail to decode JSON file: [%v]", err)
		return nil, err
	}
	logger.Log.Debug("modbusInit:GetDevConfig: Json file unmarshal done.")

	var filteredDevices []utilsPkg.Device
	logger.Log.Debug("modbusInit:GetDevConfig: Applying Modbus filter")
	for _, entry := range deviceStruct {
		for _, device := range entry.Device {
			if device.Protocol == "modbus" {
				filteredDevices = append(filteredDevices, device)
			}
		}
	}

	return filteredDevices, nil
}

func StartRead() {
	for _, device := range utilsPkg.OrigiNalDataID {
		logger.Log.Infof("modbusInit:StartRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]",
			device.Name, device.Protocol, device.Address, device.Port)

		wg := &sync.WaitGroup{} // Create a new WaitGroup for this routine
		utilsPkg.WaitGroups[device.ID] = wg
		wg.Add(1)
		utilsPkg.StopChannels[device.ID] = make(chan struct{})
		go func(device utilsPkg.Device) {
			//
			//Check PLC connection
			statusPlc := packages.ConnModbus(device)
			if statusPlc {
				packages.ReadInfoMdbs(device, device.ID)
			} else {
				logger.Log.Errorf("Device [%s] is Offline.", device.Name)
				utilsPkg.MsgLog = utilsPkg.MessageLog{
					MessageType: "error",
					Data: utilsPkg.DataLog{
						Address:  device.Address,
						Port:     device.Port,
						Name:     device.Name,
						Protocol: device.Protocol,
						Time:     time.Now().Format("2006-01-02 15:04:05"),
						Message:  "Device is Offline",
						Code:     "",
						Source:   "modbusInit:StartRead",
					},
				}
				packages.Sender(utilsPkg.MsgLog, utilsPkg.Log)
				defer wg.Done()
			}
		}(device)

		time.Sleep(100 * time.Millisecond)
	}
}

func StartIndividualDevice(device utilsPkg.Device) {
	logger.Log.Debugf("modbusInit:StartIndividualDevice: Name[%s] Protocol[%s] Ip:Port[%s:%s]",
		device.Name, device.Protocol, device.Address, device.Port)

	wg := &sync.WaitGroup{} 								// Create a new WaitGroup for this routine
	utilsPkg.WaitGroups[device.ID] = wg						// Add this WG to global variable
	wg.Add(1)												// Add a item
	utilsPkg.StopChannels[device.ID] = make(chan struct{})	// Create channel to control this routine	
	go func(device utilsPkg.Device) {
											//Start go routine		
		statusPlc := packages.ConnModbus(device)			//Check PLC connection
		if statusPlc {
			packages.ReadInfoMdbs(device, device.ID)
		} else {
			logger.Log.Errorf("Device [%s] is Off.", device.Name)
			utilsPkg.MsgLog = utilsPkg.MessageLog{
				MessageType: "error",
				Data: utilsPkg.DataLog{
					Address:  device.Address,
					Port:     device.Port,
					Name:     device.Name,
					Protocol: device.Protocol,
					Time:     time.Now().Format("2006-01-02 15:04:05"),
					Message:  "Device is Offline",
					Code:     "",
					Source:   "modbusInit:StartRead",
				},
			}
			packages.Sender(utilsPkg.MsgLog, utilsPkg.Log)
			defer wg.Done()	
		}
	}(device)
	time.Sleep(100 * time.Millisecond)
}
