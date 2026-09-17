package pkg

import (
	"encoding/json"
	"errors"
	"fmt"

	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"
	"os"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Slice to store detected changes
var changes []utilsPkg.DeviceChange

func StartReadFile() {
	for {
		originalData, err := readJSON(utilsPkg.FilePath)
		if err != nil {
			logger.Log.Errorf("checkUpDevice:StartCheckDevice: Falha ao ler o arquivo JSON: [%v]", err)
			dataLog := utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Falha ao ler o arquivo JSON. [%v]", err),
				Code:     "",
				Source:   "checkUpdDevice",
			}
			msgLog := utilsPkg.MessageLog{
				MessageType: "error",
				Data:        dataLog,
			}
			Sender(msgLog, utilsPkg.Log)
			time.Sleep(time.Second * 5) // Aguarde um pouco antes de tentar novamente.
		} else {
			logger.Log.Debug("checkUpDevice:StartCheckDevice: Leitura do arquivo JSON concluída.")
			utilsPkg.OrigiNalDataID = make(map[rune]utilsPkg.Device)
			for _, device := range originalData.Device {
				utilsPkg.OrigiNalDataID[device.ID] = device
			}
			break // Saia do loop se a leitura for bem-sucedida.
		}
	}
}


func readJSON(filePath string) (utilsPkg.DeviceStruct, error) {
	// Read the JSON file
	var devices []utilsPkg.DeviceStruct
	jsonData, err := os.ReadFile(utilsPkg.FilePath)
	if err != nil {
		logger.Log.Errorf("checkUpDevice:readJSON: Unable to read a new JSON device file. Error [%v]",
			err)
		return utilsPkg.DeviceStruct{}, err
	}

	// err = json.Unmarshal(jsonData, &devices)
	// if err != nil {
	// 	logger.Log.Errorf("checkUpDevice:readJSON: Unable to Unmarshal a new Json device file. Error [%v]",
	// 		err)
	// 	return utilsPkg.DeviceStruct{}, err
	// }

	// Unmarshal the JSON data into the slice of DeviceStruct
	err = json.Unmarshal(jsonData, &devices)
	if err != nil {
		// Provide a more specific error message if the JSON is not properly formatted
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			logger.Log.Errorf("readJSON: JSON file is not properly formatted at byte %v. Error: %v", jsonErr.Offset, err)
		} else {
			logger.Log.Errorf("readJSON: Unable to unmarshal the JSON device file. Error: %v", err)
		}
		return utilsPkg.DeviceStruct{}, err
	}

	// Check if the devices slice is empty
	if len(devices) == 0 {
		err = errors.New("readJSON: No device data found in the JSON file")
		logger.Log.Error(err)
		return utilsPkg.DeviceStruct{}, err
	}

	return devices[0], nil
}

func DiffDevices() []utilsPkg.DeviceChange {
	//Read New File Version
	updatedData, err := readJSON(utilsPkg.FilePath)
	if err != nil {
		logger.Log.Errorf("CheckUpdDevice:DiffDevices: Error reading new JSON file. Error [%v]", err)
		return nil
	}

	// Reset changes for each iteration
	changes = nil

	//Check New and Updated devices
	for _, updatedDevice := range updatedData.Device {
		originalDevice, ok := utilsPkg.OrigiNalDataID[updatedDevice.ID]
		if !ok {
			logger.Log.Infof("CheckUpdDevice:DiffDevices: New device found. ID [%v]", updatedDevice.ID)

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
			logger.Log.Infof("CheckUpdDevice:DiffDevices: Update on device [%v] found.", updatedDevice.ID)
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
			logger.Log.Infof("CheckUpdDevice:DiffDevices: Device deleted [%v] from list.", originalDevice.ID)
		}
	}

	// Update originalDevicesByID for the next comparison
	utilsPkg.OrigiNalDataID = make(map[rune]utilsPkg.Device)
	for _, updatedDevice := range updatedData.Device {
		utilsPkg.OrigiNalDataID[updatedDevice.ID] = updatedDevice
	}

	return changes
}
