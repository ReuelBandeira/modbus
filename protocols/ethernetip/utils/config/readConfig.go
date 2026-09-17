package config

import (
	"encoding/json"
	"errors"
	ethernetipPkg "ethernetip/internal/ethip"
	utilsPkg "ethernetip/utils"
	logPkg "ethernetip/utils/gologtofile"
	"os"
	"time"
)

// Start protocol Ethernet/IP
func StartEthernetip() {
	var err error
	utilsPkg.CreateChannel()
	logPkg.CtsLog.Info("Starting Ethernet/IP...")
	utilsPkg.DevicesInUse, err = GetDevConfig()
	if err != nil {
		logPkg.CtsLog.Error("FAIL - Unable to capture data from deviceConfig.json file. Error [%s]", err)
		return
	}
	utilsPkg.NumberDevice = 0
	StartRead(utilsPkg.DevicesInUse)

}

// GetDevConfig will get all informations in file
func GetDevConfig() ([]utilsPkg.DevSettings, error) {
	var devInf []utilsPkg.Devices
	var itemInf []utilsPkg.DevSettings
	var devConfig string = "deviceConfig.json"

	for {
		file, err := os.ReadFile(devConfig)
		if err != nil {
			logPkg.CtsLog.Error("Falha ao ler o arquivo JSON: %v", err)
			time.Sleep(time.Second * 5) // Aguarda 5 segundos antes de tentar novamente
			continue
		}

		err = json.Unmarshal(file, &devInf)
		if err != nil {
			//logPkg.CtsLog.Error("Falha ao decodificar o arquivo JSON: %v", err)
			// Aguarda 5 segundos antes de tentar novamente

			var jsonErr *json.SyntaxError
			if errors.As(err, &jsonErr) {
				logPkg.CtsLog.Error("readJSON: JSON file is not properly formatted at byte %v. Error: %v", jsonErr.Offset, err)
			} else {
				logPkg.CtsLog.Error("readJSON: Unable to unmarshal the JSON device file. Error: %v", err)
			}
			time.Sleep(time.Second * 5)
			continue
		}
		// Check if the devices slice is empty
		if len(devInf) == 0 {
			err = errors.New("readJSON: No device data found in the JSON file")
			logPkg.CtsLog.Error("%v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		
		// Se o arquivo foi lido e deserializado com sucesso, saia do loop
		break
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			itemInf = append(itemInf, utilsPkg.DevSettings{
				Id:          item.Id,
				Address:     item.Address,
				Port:        item.Port,
				Name:        item.Name,
				Protocol:    item.Protocol,
				ReadingTime: item.ReadingTime,
				Topics:      item.Topics,
				Data:        item.Data,
			})
		}
	}
	return itemInf, nil
}

func StartRead(devices []utilsPkg.DevSettings) {
	for _, device := range devices {
		logPkg.CtsLog.Info("StartRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
			device.Name, device.Protocol, device.Address, device.Port)
		statusPlc := ethernetipPkg.ConnEthernetIp(device)

		if statusPlc {
			ethernetipPkg.SendStatusDevice(device, "connected")
			logPkg.CtsLog.Info("Ok - Dispositivo está online")
		} else {
			ethernetipPkg.SendStatusDevice(device, "not connected")
			logPkg.CtsLog.Info("Erro - Dispositivo não está online")
		}
		utilsPkg.Wg.Add(1)
		utilsPkg.NumberDevice++
		go func(device utilsPkg.DevSettings) {
			ethernetipPkg.ReadInfoEth(device)
		}(device)
		time.Sleep(1 * time.Second)
	}
	ethernetipPkg.SendStatusProtocol("running")
}

func CheckNewDevices(devices []utilsPkg.DevSettings) bool {
	findDevice := false
	devNew, err := GetDevConfig()
	if err != nil {
		logPkg.CtsLog.Error("FAIL - Unable to capture data from deviceConfig.json file. Error [%s]", err)
		return false
	}
	for _, device := range devNew {
		for _, devActual := range devices {
			if (devActual.Name == device.Name) && (devActual.Protocol == device.Protocol) {
				findDevice = true
			}
		}
		if !findDevice {
			return true
		}
		findDevice = false
	}
	return false
}

/*
func StartWrite(devices []utilsPkg.DevSettings, client mqtt.Client, loop bool) {
	for _, device := range devices {

		if device.Protocol == "ethernetip" {
			logPkg.CtsLog.Info("StartWrite: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
				device.Name, device.Protocol, device.Address, device.Port)
			statusPlc := ethernetipPkg.ConnEthernetIp(device, client)
			if statusPlc {
				ethernetipPkg.MQTTSendEthStatusDevice(device, SettingsMqtt, statusPlc)
				ethernetipPkg.ReadInfoEth(device, SettingsMqtt, client, loop, "write")
			} else {
				ethernetipPkg.MQTTSendEthStatusDevice(device, SettingsMqtt, statusPlc)
			}
		}
		time.Sleep(1 * time.Second)
	}

}
*/
