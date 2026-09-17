package utils

import (
	"os"
	"path/filepath"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/v3/process" // 28/11/2023 Added to check if a proccess is running
)

// ========================================
// Global vaiables commons for any protocol
// ========================================
var Hostname string

const OsBits = 32 << (^uint(0) >> 63) // 32 or 64 bits
var ProtocolIsRunning bool = false
var DevicelUpdateRequest bool = false
var AnyDevicelsRunning bool = true
var AnyDeviceIsConnected bool = false
var Status bool = true
var MsgLog MessageLog

// MQTT Client and Local MQTT Broker Global vars
var LocalMqttBrockerAddress string = "localhost"
var LocalMqttBrockerPort string = "1883"
var LocalMqttBrockerUserName string = ""
var LocalMqttBrockerPassword string = ""
var MqttClient mqtt.Client /// MQTT Client Handler

// default:Logs
var DefaultLogFilePath string = "."
var DefaultLogFileName string = "newprotocol"
var LogFileMaxSize int = 500000
var DefaultLogLevel int = 5 // DebugLevel

var DeleteExistingLogFiles bool = true

// default:Backend
var DefaultBackEndListenAddr string = "localhost"
var DefaultBackEndListenPort string = "8585"

// default:configPath
var DefaultConfigPath string = "."
var DefaultConfigFileName string = "deviceConfig.json"

// default:Devices Config EndPoint and file
var DefaultBackEndDevConfEndPoint string = "api/v1/configurations"

// 14/10/2023 - BEGIN FROM ROLDEN General & LOGS
var General []string = []string{"general"}
var Log []string = []string{"log"}

var FilePath string = "deviceConfig.json"

var DeviceWorkStatus = make(map[rune]bool)

type BitMemory struct {
	Address uint16 `json:"address"`
	Name    string `json:"name"`
}

type WordMemory struct {
	Address uint16 `json:"address"`
	Name    string `json:"name"`
	Format  uint8  `json:"format"`
}
type Data struct {
	BitMemories  []BitMemory  `json:"bitMemories"`
	SlaveID      uint8        `json:"slaveId"`
	WordMemories []WordMemory `json:"wordMemories"`
}

type Device struct {
	Address     string   `json:"address"`
	Port        string   `json:"port"`
	Name        string   `json:"name"`
	Protocol    string   `json:"protocol"`
	ReadingTime int      `json:"readingTime"`
	ID          rune     `json:"id"`
	Topics      []string `json:"topics"`
	Data        []Data   `json:"data"`
}

type DeviceStruct struct {
	Device []Device `json:"devices"`
}

//Create struct to save info

var (
	OrigiNalDataID (map[rune]Device)
	UpdatedDataID  (map[rune]Device)
)

type ChangeType int

const (
	New ChangeType = iota
	Deleted
	Updated
)

type DeviceChange struct {
	Device Device
	Type   ChangeType
}

var (
	StopChannels map[rune]chan struct{}
	WaitGroups   map[rune]*sync.WaitGroup
)

func CreateFromHoldenChannel() {
	StopChannels = make(map[rune]chan struct{})
	WaitGroups = make(map[rune]*sync.WaitGroup)
}

// 14/10/2023 - END FROM ROLDEN  General & LOGS
type MqttBrockerStruct struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct Message Error/Info
type MessageLog struct {
	MessageType string  `json:"messageType"`
	Data        DataLog `json:"data"`
}
type DataLog struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Source   string `json:"source"`
}

// Formato da mensagem de status do protocolo
type MessageStatusProtocol struct {
	MessageType string            `json:"messageType"`
	Data        MessageDataStatus `json:"data"`
}

type MessageDataStatus struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

// Formato da mensagem de status do device
type MessageDeviceStatus struct {
	MessageType string                  `json:"messageType"`
	Data        MessageDataDeviceStatus `json:"data"`
}

type MessageDataDeviceStatus struct {
	Protocol string `json:"protocol"`
	Device   string `json:"device"`
	Id       int    `json:"id"`
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
	Id       int    `json:"id"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Source   string `json:"source"`
}

type MqttMsgStruct struct {
	Address        Address        `json:"address"`
	ReadTimeStamp  string         `json:"readTimeStamp"`
	Protocol       string         `json:"protocol"`
	Id             int            `json:"id"`
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

type Info struct {
	Status   string `json:"status"`
	Protocol string `json:"protocol"`
}

// Funções do channel de controle das Go Routines
var (
	StopChan          chan struct{} // Canal de sinalização para encerrar as goroutines
	UpdateDevChan     chan struct{} // Canal de sinalização para iniciar um update
	WaitGroup         sync.WaitGroup
	NumberDevice      int // Número de devices ativos
	UpdateStatusCount int // Indica quantos devices falta executar update
)

func CreateChannel() {
	StopChan = make(chan struct{})
	UpdateDevChan = make(chan struct{})
}

// TODO: Specific Protocol
type NewDeviceExitPayloadJsonMsg struct {
	Address       string
	Port          string
	Name          string
	Others        string
	Protocol      string
	Id            int
	ReadTimeStamp string
	Topics        []string
	AI            map[string]interface{} // TODO: Specific Protocol
	AO            map[string]interface{} // TODO: Specific Protocol
	DI            map[string]interface{} // TODO: Specific Protocol
	DO            map[string]interface{} // TODO: Specific Protocol
}

// BEGIN 28/11/2023 Added to check if it proccess is already running
func Running() bool {

	processes := []string{}

	myself, _ := os.Executable()
	myself = filepath.Base(myself)

	v, _ := process.Processes()

	for _, p := range v {
		proc, err := p.Name()
		if err == nil && proc == myself {
			processes = append(processes, proc)
		}
	}

	return len(processes) > 1
}

// END 28/11/2023 Added to check if it proccess is already running
