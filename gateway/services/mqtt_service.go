package services

import (
	"encoding/json"
	"gateway/globals"
	"gateway/models"
	"io"
	"log"
	"net/http"
)

func GetMQTTConfig() (models.MQTT, error) {

	var mqtt models.MQTT

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8585/api/v1/mqtt", nil)

	req.Header.Set("Api-key", globals.GenerateApiKey())

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("FAIL to read MQTT configuration from Backend: %s\n", err)
		return models.MQTT{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("FAIL to read MQTT configuration from Backend: %s\n", err)
		return models.MQTT{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FAIL to read response body.\n")
		return models.MQTT{}, err
	}

	err = json.Unmarshal(body, &mqtt)
	if err != nil {
		log.Printf("FAIL unmarshaling JSON for MQTT Configurations: %v\n", err)
		return models.MQTT{}, err
	}

	return mqtt, nil
}
