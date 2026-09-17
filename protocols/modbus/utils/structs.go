package utils

import (
	//"runtime"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Config Local MQTT
type ConfigMQTT struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Formato da mensagem de entrada esperado pelo CallBack
// TODO: Deixar informação em repositório externo
// type InputMQTTMsg struct {
// 	Action string `json:"action"`
// }

var MqttClient mqtt.Client
var MsgLog MessageLog

//var DataMsgLog DataLog

// var StatusProtocol bool = true
var DeviceWorkStatus = make(map[rune]string)

// Struct to send message to MQTT
type MqttMsgStruct struct {
	Address       Address  `json:"address"`
	ReadTimeStamp string   `json:"readTimeStamp"`
	Protocol      string   `json:"protocol"`
	Topics        []string `json:"topics"`
	Data          DataExit `json:"data"`
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Name    string `json:"name"`
	ID      rune   `json:"id"`
	Others  string `json:"others"`
}

type DataExit struct {
	BitMemories  map[string]bool   `json:"bitMemories"`
	WordMemories map[string]uint32 `json:"wordMemories"`
}

//End Message Device MQTT

// Decode da mensagem de entrada
type MessageInput struct {
	MessageType string      `json:"messageType"`
	Data        interface{} `json:"data"`
}

// End input message

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
	ID       rune   `json:"id"`
}

//Fim

// Funções do channel de controle das Go Routines
var (
	StopChannels map[rune]chan struct{}
	WaitGroups   map[rune]*sync.WaitGroup
)

func CreateChannel() {
	StopChannels = make(map[rune]chan struct{})
	WaitGroups = make(map[rune]*sync.WaitGroup)
}

//End Go Routines Control

// Add new structs
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

// Struct Message Error
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
