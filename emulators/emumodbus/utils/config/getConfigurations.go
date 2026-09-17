package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	/// TODO: change newprotocol with your protocolname
	utilsPkg "emumodbus/utils"
	golog "emumodbus/utils/gologtofile"

	"os"
	"reflect"
	"time"
)

// default:Logs
var DefaultLogFilePath string = "."
var DefaultLogFileName string = "emu_modbus"
var DefaultLogLevel int = 4 // InfoLevel

var DeleteExistingLogFiles bool = true

// default:Backend
var DefaultBackEndListenAddr string = "127.0.0.1"
var DefaultBackEndListenPort string = "8585"

// default:configPath
var DefaultConfigPath string = "."
var DefaultConfigFileName string = "deviceConfig.json"

// default:Devices Config EndPoint and file
var DefaultBackEndDevConfEndPoint string = "api/v1/configurations"

func GetBackEndDevicesConfig() {
	var previousDevSettings []utilsPkg.Devices

	for {
		currentDevSettings, err := getDevSettings()

		if err != nil {
			fmt.Printf("Erro ao obter a informação do endpoint: %s\n", err)
			time.Sleep(5 * time.Second) // Aguarde um tempo antes de tentar novamente
			continue
		}

		isEqual := reflect.DeepEqual(currentDevSettings, previousDevSettings)

		if !isEqual {
			err := saveFileDevice(currentDevSettings)
			if err != nil {
				fmt.Printf("Erro ao salvar a informação em um arquivo: %s\n", err)
			} else {
				previousDevSettings = currentDevSettings
			}
		}
	}

}

func getDevSettings() ([]utilsPkg.Devices, error) {
	var devices []utilsPkg.Devices

	// get http request from Backend ex: url:http://127.0.0.1:8585/api/v1/configurations
	url := "http://" + DefaultBackEndListenAddr + ":" + DefaultBackEndListenPort + "/" + DefaultBackEndDevConfEndPoint
	resp, err := http.Get(url)
	if err != nil {
		golog.CtsLog.Error("GetBackEndDevicesConfig:FAIL to access backend server: err[%v]\n", err)
		return devices, err
	}
	defer resp.Body.Close()

	// read response body from http Response it MUST be a JSON formatted as devices
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		golog.CtsLog.Error("Failed to read response body: err[%v]\n", err)
		return devices, err
	}

	// Extract response body JSON format into devices[utilsPkg.Devices]
	err = json.Unmarshal(body, &devices)
	if err != nil {
		golog.CtsLog.Error("Failed to Unmarshal PLCINFO configuration: err[%v]\n", err)
		return devices, err
	}
	return devices, nil
}

func saveFileDevice(device []utilsPkg.Devices) error {
	// Open or Local MqttConfigFileName at ConfigPath
	file, err := os.Create(DefaultConfigPath + "/" + DefaultConfigFileName)
	if err != nil {
		golog.CtsLog.Error("Failed to create file:[%s/%s] err[%v]\n", DefaultConfigPath, DefaultConfigFileName, err)
		return err
	}
	defer file.Close()

	// Encode  devices as JSON
	encodedData, err := json.MarshalIndent(device, "", "    ")
	if err != nil {
		golog.CtsLog.Error("Failed to encode devices as JSON: err[%v]\n", err)
		return err
	}

	// Save  JSON encoded data into devices
	_, err = file.Write(encodedData)
	if err != nil {
		golog.CtsLog.Error("Failed to write encodedData to file: err[%v]\n", err)
		return err
	}
	return nil

}
