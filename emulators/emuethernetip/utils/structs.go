package utils

import (
	//"emuethernetip/protocols/ethernetip"
	//"emuethernetip/protocols/ethernetip/commonIndustrialProtocol"
	//_type "emuethernetip/protocols/ethernetip/type"
	//"net"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var DevicesInUse []DevSettings // Estrutura de Devices em uso

// default:Logs
var DefaultGwIsi40LogFilePath string = "./"
var DefaultGwIsi40LogFileName string = "GwIsi40"
var DefaultGwIsi40LogLevel int = 5 // DebugLevel
var DeleteExistingLogFiles bool = true
var DeviceConnected []string

// Informações do MQTT Broker Local
var MqttAddress string = "localhost"
var MqttPort string = "1883"
var MqttUser string = ""
var MqttPassword string = ""

// Códigos de erro do protocolo CIP
var CodeCIPerrors = [...]string{
	"Success",
	"Connection failure",
	"Resource unavailable",
	"Invalid parameter value",
	"Path segment error",
	"Path destination unknow",
	"Only part of the expected data was transferred",
	"The messaging connection was lost",
	"Service not supported",
	"Invalid attribute data detected",
	"Attribute list error",
	"Already in requested mode/state",
	"Object state conflict",
	"Object already exists",
	"Attribute not settable",
	"Privilege violation",
	"Device state conflict",
	"The data to be transmissed in the response buffer is larger than the allocated response buffer",
	"Fragmentation of a primitive value",
	"Not enough data",
	"Attribute not supported",
	"Too much data",
	"Object does not exist",
	"Sevice fragmentation sequence not in progress",
	"No stored attribute data",
	"Store operation failure",
	"Routing failure. request packet too large",
	"Routing failure. response packet too large",
	"Missing attribute list entiy data",
	"lnvalid attribute value list",
	"Embedded service error",
	"Vendor specific error",
	"A parameter associated with the request was invalid",
	"Write-once value or medium already written",
	"lnvalid Reply Received",
	"Buffer Overflow",
	"Message Format Error",
	"Message Format Error",
	"Path Size Invalid",
	"Unexpected attribute in list",
	"Invalid Member ID",
	"Member not settable",
	"Group 2 only server General failure",
	"Unknown Modbus error",
	"Attribute not gettable"}

type ConfigMQTT struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Formato da mensagem de entrada esperado pelo CallBack
// TODO: Deixar informação em repositório externo
//type InputMQTTMsg struct {
//	Data struct {
//		Action string `json:"action"`
//	}
//}

// Formato da mensagem de status do protocolo
type MessageStatusProtocol struct {
	MessageType string            `json:"messageType"`
	Data        MessageDataStatus `json:"data"`
}

type MessageDataStatus struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

//Fim

// Formato da mensagem de erro do device
type MessageDeviceError struct {
	MessageType string                 `json:"messageType"`
	Data        MessageDataDeviceError `json:"data"`
}

// Formato da mensagem de status do device
type MessageDeviceStatus struct {
	MessageType string                  `json:"messageType"`
	Data        MessageDataDeviceStatus `json:"data"`
}

type MessageDataDeviceError struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Source   string `json:"source"`
}
type MessageDataDeviceStatus struct {
	Protocol string `json:"protocol"`
	Device   string `json:"device"`
	Status   string `json:"status"`
}

type GwIsi40Settings struct {
	LogFilePath               string `json:"LogFilePath"`               // "../../logs"
	LogFileName               string `json:"LogFileName"`               //  "Gwisi40"
	LogLevel                  string `json:"LogLevel"`                  // "5=DebugLevel"
	DelLogFiles               bool   `json:"DelLogFiles"`               // "true"
	ConfigPath                string `json:"ConfigPath"`                //"../../utils/config"
	ConfigMqttFileName        string `json:"ConfigMqttFileName"`        // "mqttConfig.json"
	ConfigPlcFileName         string `json:"ConfigPlcFileName"`         // "plcConfig.json"
	BackendHttpAddr           string `json:"backendHttpAddr"`           // "127.0.0.1"
	BackendHttpPort           string `json:"BackendHttpPort"`           // "8585"
	BackendApiGwIsi40EndPoint string `json:"BackendApiGwIsi40EndPoint"` // "api/v1/gwisi40",
	BackendApiMqttEndPoint    string `json:"BackendApiMqttEndPoint"`    // "api/v1/mqtt"
	BackendApiPlcEndPoint     string `json:"BackendApiPlcEndPoint"`     // "api/v1/configurations"
}

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

type SettingsMQTT struct {
	Server       string `json:"server"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	InputInfo    string `json:"inputInfo"`
	ErrorInfo    string `json:"errorInfo"`
	StatusDevice string `json:"statusDevice"`
}

type ExitPayloadMsg struct {
	Address       string
	Port          string
	Name          string
	Others        string
	Protocol      string
	ReadTimeStamp string
	Topics        []string
	BitMemories   map[string]interface{}
	WordMemories  map[string]interface{}
}

// Decode da mensagem de entrada
type MessageInput struct {
	MessageType string      `json:"messageType"`
	Data        interface{} `json:"data"`
}

var MqttClient mqtt.Client
var StatusProtocol bool = true

// var LoopEthip bool = true
var Status bool = true

type MqttMsgStruct struct {
	Address       Address `json:"address"`
	ReadTimeStamp string  `json:"readTimeStamp"`
	Protocol      string  `json:"protocol"`
	Data          Data    `json:"data"`
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Name    string `json:"name"`
	Others  string `json:"others"`
}

type Data struct {
	BitMemories  map[string]interface{} `json:"bitMemories"`
	WordMemories map[string]interface{} `json:"wordMemories"`
}

// Send Status Protocol
//
//	type StatusProtocol struct {
//		Operation Info `json:"operation"`
//	}
type Info struct {
	Status   string `json:"status"`
	Protocol string `json:"protocol"`
}

// Funções do channel de controle das Go Routines
var (
	StopChan          chan struct{} // Canal de sinalização para encerrar as goroutines
	UpdateDevChan     chan struct{} // Canal de sinalização para iniciar um update
	Wg                sync.WaitGroup
	NumberDevice      int // Número de devices ativos
	UpdateStatusCount int // Indica quantos devices falta executar update
)

func CreateChannel() {
	StopChan = make(chan struct{})
	UpdateDevChan = make(chan struct{})
}
