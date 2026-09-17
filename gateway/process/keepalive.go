package process

import (
	"encoding/json"
	"gateway/globals"
	"gateway/mqtt"
	"time"
)

func SendProtocolsKeepAlive() {
	type MessageRebootProtocol struct {
		MessageType string `json:"messageType"`
		Data        struct {
			Action string `json:"action"`
			Name   string `json:"protocol"`
		} `json:"data"`
	}

	for {

		for _, protocol := range globals.GatewayProtocols {
			var data MessageRebootProtocol

			data.MessageType = "action"
			data.Data.Action = "start"

			data.Data.Name = protocol.Protocol

			jsonPayload, err := json.Marshal(data)
			if err != nil {
				panic(err) //Check - não pode parar o programa
			}

			mqtt.SendMsg(globals.LocalClient, protocol.Protocol, string(jsonPayload))
		}

		// Sleeps for the interval specified in the configuration file
		time.Sleep(time.Duration(globals.Conf.KeepAlive.Interval) * time.Millisecond)
	}
}
