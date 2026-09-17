package messages

import (
	"encoding/json"
	"errors"
)

// ============================================================
// Actions
type MQTTActionMsg struct {
	MessageType string         `json:"messageType"`
	Data        MQTTActionData `json:"data"`
}

type MQTTActionData struct {
	Action   string `json:"action"`
	Protocol string `json:"protocol"`
}

// ============================================================
// Status
type MQTTStatusMsg struct {
	MessageType string         `json:"messageType"`
	Data        MQTTStatusData `json:"data"`
}

type MQTTStatusData struct {
	Protocol string `json:"protocol"`
	Device   string `json:"device"`
	Id       uint   `json:"id"`
	Status   string `json:"status"`
}

// ============================================================
// Generic messages from MQTT
type Message struct {
	MessageType string
	Status      string
	Name        string
	Protocol    string
	Action      string
	Device      string
	Id          uint
}

func DecodeMessage(jsonText []byte) (Message, error) {
	var msg Message
	var parsed any

	err := json.Unmarshal(jsonText, &parsed)
	if err != nil {
		return msg, err
	}

	switch val := parsed.(type) {

	case nil:
		return msg, errors.New("json specifies null")

	case map[string]any:
		msg.MessageType = val["messageType"].(string)
		if dataVal, ok := val["data"].(map[string]interface{}); ok {
			if protocolVal, ok := dataVal["protocol"].(string); ok {
				msg.Protocol = protocolVal
			}
			if nameVal, ok := dataVal["name"].(string); ok {
				msg.Name = nameVal
			}
			if idVal, ok := dataVal["id"].(float64); ok {
				msg.Id = uint(idVal)
			}
			if statusVal, ok := dataVal["status"].(string); ok {
				msg.Status = statusVal
			}
			if actionVal, ok := dataVal["action"].(string); ok {
				msg.Action = actionVal
			}
			if deviceVal, ok := dataVal["device"].(string); ok {
				msg.Device = deviceVal
			}
		}

	default:
		return msg, errors.New("unexpected type on json")
	}

	return msg, nil
}
