package utils

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HostName
var Hostname string

const OsBits = 32 << (^uint(0) >> 63)

var StatusProtocol bool = true

// var Status bool = true
var Status bool = true

// Setup MQTT Client and Local MQTT Broker info
var MqttClient mqtt.Client /// MQTT Client Handler
var LocalMqttBrockerAddress string = "localhost"
var LocalMqttBrockerPort string = "1883"
var LocalMqttBrockerUserName string = ""
var LocalMqttBrockerPassword string = ""

// TCP Server Namd/IpAddress + Port
var TcpSvrpAddress string = "localhost"
var TcpSvrPort uint16 = 502 // Modbus: Emulate CANopen TCP Server Port
// TCP Console Server Port
var TcpConsolePort uint16 = TcpSvrPort + 1
var TcpConnectionId int = 0

// TCP Console Cmd Flags
var StartServer = false
var StopServer = false
var ServerStopped = false
var RestartServer = false
var ServerRestarted = false
var BlockResponse = false

var ForceCrcError = false
var CrcErrorForced = false

type MqttBrockerStruct struct {
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

var StopChan chan bool // Canal de sinalização para encerrar as goroutines

var WaitGroupGoroutines sync.WaitGroup

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

// Formato da mensagem de status do device
type MessageDeviceStatus struct {
	MessageType string                  `json:"messageType"`
	Data        MessageDataDeviceStatus `json:"data"`
}

type MessageDataDeviceStatus struct {
	Protocol string `json:"protocol"`
	Device   string `json:"device"`
	Status   string `json:"status"`
}

type ProtocolSettings struct {
	LogFilePath               string `json:"LogFilePath"`               // "../../logs"
	LogFileName               string `json:"LogFileName"`               //  "Gwisi40"
	LogLevel                  string `json:"LogLevel"`                  // "5=DebugLevel"
	DelLogFiles               bool   `json:"DelLogFiles"`               // "true"
	ConfigPath                string `json:"ConfigPath"`                //"../../utils/config"
	ConfigMqttFileName        string `json:"ConfigMqttFileName"`        // "mqttConfig.json"
	ConfigDeviceFileName      string `json:"ConfigDeviceFileName"`      // "deviceConfig.json"
	BackendHttpAddr           string `json:"backendHttpAddr"`           // "127.0.0.1"
	BackendHttpPort           string `json:"BackendHttpPort"`           // "8585"
	BackendApiGwIsi40EndPoint string `json:"BackendApiGwIsi40EndPoint"` // "api/v1/gwisi40",
	BackendApiMqttEndPoint    string `json:"BackendApiMqttEndPoint"`    // "api/v1/mqtt"
	BackendApiDeviceEndPoint  string `json:"BackendApiDeviceEndPoint"`  // "api/v1/configurations"
}

type Devices struct {
	Devices []DevSettings `json:"devices"`
}

type DevSettings struct {
	Address     string                   `json:"address"`
	Port        string                   `json:"port"`
	Name        string                   `json:"name"`
	Protocol    string                   `json:"protocol"`
	Id          int                      `json:"id"` // Added 24/08/2023
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

type NewDeviceExitPayloadJsonMsg struct {
	Address       string
	Port          string
	Name          string
	Others        string
	Protocol      string
	ReadTimeStamp string
	Topics        []string
	AI            map[string]interface{}
	AO            map[string]interface{}
	DI            map[string]interface{}
	DO            map[string]interface{}
}

// Decode da mensagem de entrada
type MessageInput struct {
	MessageType string      `json:"messageType"`
	Data        interface{} `json:"data"`
}

type MessageDeviceLog struct { // Added 24/08/2023
	MessageType string     `json:"messageType"` // Added 24/08/2023
	Data        MqttLogMsg `json:"data"`        // Added 24/08/2023
}

// / Must be inside : MessageInput struct
type MqttLogMsg struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Protocol string `json:"prootocol"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Source   string `json:"source"`
}

type MqttMsgStruct struct {
	Address        Address        `json:"address"`
	ReadTimeStamp  string         `json:"readTimeStamp"`
	Protocol       string         `json:"protocol"`
	JsonDeviceData JsonDeviceData `json:"data"`
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Name    string `json:"name"`
	Others  string `json:"others"`
}

type JsonDeviceData struct {
	AI map[string]interface{} `json:"AI"`
	AO map[string]interface{} `json:"AO"`
	DI map[string]interface{} `json:"DI"`
	DO map[string]interface{} `json:"DO"`
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

// GoRotines Channel Control
var (
	DoneChan      chan bool
	UpdateDevChan chan bool
	SyncOnceChan  sync.Once
	SyncWaitGroup sync.WaitGroup
)

func CreateGoRotineChannel() {
	DoneChan = make(chan bool)
	UpdateDevChan = make(chan bool)
}
