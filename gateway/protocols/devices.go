package protocols

import (
	"bytes"
	"gateway/globals"
	"gateway/messages"
	"gateway/readercomp"

	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func GetDevices() {

	var err error

	for _, protocol := range globals.Dictionary {
		globals.CurrentDevSettings, err = GetBackendDevSettings(protocol)

		if err != nil {
			log.Printf("Erro ao obter a informação de configurações dos dispositivos do endpoint: %s\n", err)
		}
	}
}

func GetBackendDevSettings(protocol string) ([]globals.Devices, error) {
	var devices []globals.Devices

	log.Printf("Buscando configurações de dispositivos do backend para o protocolo %s\n", protocol)

	// Get device configurations from the backend
	url := "http://localhost:8585/api/v1/configurationsbyprotocol/" + protocol

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Api-key", globals.GenerateApiKey())

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Erro ao obter a informação do endpoint: %s\n", err)
		return devices, err
	}
	defer resp.Body.Close()

	// Get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao obter a informação do body recebido do endpoint: %s\n", err)
		return devices, err
	}

	err = saveFileDevice(protocol, body)
	if err != nil {
		log.Printf("Erro ao salvar configurações de dispositivos do protocolo %s\n", err)
		return devices, err
	}

	return devices, nil
}

func saveFileDevice(protocol string, devices []byte) error {

	log.Printf("Salvando configurações de dispositivos do protocolo %s\n", protocol)

	// Official file is saved on the same directory that the protocol executable is
	filename := "../protocols/" + protocol + "/deviceConfig.json"

	// Temporary file, used just to compare with the previous official file
	tempfilename := "../protocols/" + protocol + "/tempDeviceConfig.json"

	// Open file for writing
	file, err := os.OpenFile(tempfilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create file: err[%v]\n", err)
		return err
	}
	defer file.Close()

	// Save JSON into the file
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, devices, "", "\t")
	_, err = file.Write(prettyJSON.Bytes())
	if err != nil {
		log.Printf("Failed to write encodedData to file: err[%v]\n", err)
		return err
	}

	file.Close()

	// Compare the new file with the previous one
	result, err := readercomp.FilesEqual(filename, tempfilename)
	if err != nil && result {
		log.Printf("Failed to compare files: err[%v]\n", err)
		return err
	}

	if !result {
		// Files are different. Remove old official file and rename the new file to the official one
		_, err := os.Stat(filename)
		if err == nil {
			err := os.Remove(filename)
			if err != nil {
				log.Printf("Failed to remove file: err[%v]\n", err)
				return err
			}
		}

		err = os.Rename(tempfilename, filename)
		if err != nil {
			log.Printf("Failed to rename file: err[%v]\n", err)

			return err
		}

		// Send MQTT message to the correspondent protocol to reload the configurations
		mqttData := messages.MQTTActionData{
			Action:   "update",
			Protocol: protocol,
		}

		msg := messages.MQTTActionMsg{
			MessageType: "action",
			Data:        mqttData,
		}

		// Notify the correspondent protocol that the configurations has changed via local MQTT
		m, _ := json.Marshal(msg)
		globals.LocalClient.Publish("general", 0, false, string(m))

	} else {
		// Files are equal, indicating that there are no changes in the configurations
		// Thus, just remove the temporary file
		os.Remove(tempfilename)
		return nil
	}

	return nil
}
