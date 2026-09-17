package mqtt

import (
	"encoding/json"
	"gateway/globals"
	"gateway/messages"
	"gateway/protocols"
	"gateway/services"
	"log"

	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectGlobalBroker() {

	// Gets all MQTT configurations from database
	mymqtt, err := services.GetMQTTConfig()

	log.Printf("External MQTT: %s:%s\n", mymqtt.Server, mymqtt.Port)
	if err == nil && mymqtt.Server != "" && mymqtt.Port != "" {
		if globals.GlobalClient == nil || !globals.GlobalClient.IsConnected() {

			// Configure Global MQTT Broker
			optsg := mqtt.NewClientOptions().AddBroker("tcp://" + mymqtt.Server + ":" + mymqtt.Port)
			globals.GlobalClient = mqtt.NewClient(optsg)
			log.Println("Connecting to external MQTT Broker on " + mymqtt.Server + ":" + mymqtt.Port)

			if token := globals.GlobalClient.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("Can not connect to MQTT: %v\n", token.Error())
				time.Sleep(30 * time.Second)
			}

			// Subscribe to MQTT topics
			SubscribeMQTTGlobalMessages()

		} else {
			log.Println("Can not connect to MQTT: ", err)
			time.Sleep(30 * time.Second)
		}
	}
}

func SubscribeMQTTGlobalMessages() {

	if globals.GlobalClient.IsConnected() {
		// Receives MQTT general messages from Global MQTT
		// Replicates to local MQTT
		globals.GlobalClient.Subscribe("general", 0, func(client mqtt.Client, msg mqtt.Message) {
			if globals.LocalClient.IsConnected() {

				message, err := messages.DecodeMessage(msg.Payload())
				if err != nil {
					log.Printf("Error decoding message: %s\n", err)
				}

				if message.MessageType == "action" {
					token := globals.LocalClient.Publish(message.Protocol, 0, false, msg.Payload())
					token.Wait()
					if token.Error() != nil {
						log.Printf("Can not publish to Local MQTT: %v\n", token.Error())
					}
				}
			} else {
				log.Println("Can not publish to Local MQTT: not connected")
			}
		})
	}
}

func ConnectLocalBroker() bool {
	optsl := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	globals.LocalClient = mqtt.NewClient(optsl)
	log.Println("Connecting to local MQTT Broker on tcp://127.0.0.1:1883")
	return ConnectBroker(globals.LocalClient)
}

func ConnectBroker(client mqtt.Client) bool {
	// Connects to MQTT broker. Retries 10 times if connection fails.
	for retries := 0; retries < 10; retries++ {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			time.Sleep(5 * time.Second)
			if retries > 5 {
				log.Printf("Can not connect to MQTT: %v\n", token.Error())
				return false
			}
			log.Println("MQTT connection failed. Retrying " + strconv.Itoa(retries))
		}
	}

	return true
}

func SendMsg(client mqtt.Client, topic string, payload string) {
	if client.IsConnected() {
		// Publishes to MQTT topics
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Can not publish to MQTT: %v\n", token.Error())
		}
	} else {
		log.Println("Can not publish to MQTT: client is not connected")
	}
}

func MQTTProtocolMessages(protocol string) {

	// Receives MQTT log messages from local MQTT
	// Saves log messages to log file
	if globals.LocalClient.IsConnected() {
		globals.LocalClient.Subscribe("log", 0, func(client mqtt.Client, msg mqtt.Message) {
			log.Printf("%v\n", msg.Payload())
		})
		globals.LocalClient.Subscribe("gateway", 0, func(client mqtt.Client, msg mqtt.Message) {
			type msgBackend struct {
				Protocol string `json:"protocol"`
			}
			var m msgBackend
			_ = json.Unmarshal(msg.Payload(), &m)

			if m.Protocol == "mqtt" {
				log.Printf("MQTT configuration changed on backend\n")
				services.GetMQTTConfig()
				ConnectGlobalBroker()
			} else {
				log.Printf("Devices for protocol %s changed on backend\n", m.Protocol)
				protocols.GetBackendDevSettings(m.Protocol)
			}
		})
	} else {
		log.Println("Can not receive MQTT log messages: not connected")
	}

	// Receives MQTT protocol messages from local MQTT
	// Replicates to Global MQTT ==> will be captured by Grafana, for instance
	globals.LocalClient.Subscribe(protocol, 0, func(client mqtt.Client, msg mqtt.Message) {
		if globals.GlobalClient.IsConnected() {
			// Publishes to MQTT topics
			token := globals.GlobalClient.Publish(protocol, 0, false, msg.Payload())
			token.Wait()
			if token.Error() != nil {
				log.Printf("Can not publish to Global MQTT: %v\n", token.Error())
			}
		} else {
			log.Println("Can not publish to Global MQTT: not connected")
		}
	})
}
