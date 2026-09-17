package globals

import (
	"encoding/base64"
	"gateway/models"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Devices struct {
	Devices []DevSettings `json:"devices"`
}

type DevSettings struct {
	Address     string                   `json:"address"`
	Port        string                   `json:"port"`
	Name        string                   `json:"name"`
	Protocol    string                   `json:"protocol"`
	ReadingTime int                      `json:"readingTime"`
	Topics      []string                 `json:"topics"`
	Data        []map[string]interface{} `json:"data"`
}

var (
	// Will contain all configurations from config file
	Conf *models.Configurations

	// Logger
	Log *log.Logger

	// api key needed to access backend
	APIK string

	// Protocols served by the gateway
	GatewayProtocols []models.GatewayProtocol

	// Local and Global MQTT clients
	LocalClient  mqtt.Client
	GlobalClient mqtt.Client

	// Processes running on the gateway
	Processes []models.Process

	// Devices configurations
	DevicesConf []DevSettings

	// Previous devices configurations
	PreviousDevSettings []Devices

	// Most recent devices configurations
	CurrentDevSettings []Devices

	// Dictionary of protocols' status
	Dictionary map[string]string

	// Path separator
	PathSeparator string
)

func GenerateApiKey() string {

	apik := "gateway-" + Conf.Backend.Host
	return base64.StdEncoding.EncodeToString([]byte(apik))

}
