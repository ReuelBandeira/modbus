package mqtt

import (
	"gwisi40server/database"
	"gwisi40server/messages"
	"gwisi40server/models"
	"gwisi40server/status"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	mqttServer "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

var (
	GlobalMQTTClient mqtt.Client
	LocalMQTTClient  mqtt.Client
)

func StartLocalMQTTServer() {
	// Create the new MQTT Server.
	server := mqttServer.New(nil)

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func SubscribeStatus() {
	LocalMQTTClient.Subscribe("general", 0, func(client mqtt.Client, msg mqtt.Message) {
		message, err := messages.DecodeMessage(msg.Payload())
		if err != nil {
			log.Printf("Error decoding message: %s\n", err)
		}

		if message.MessageType == "deviceStatus" {
			//  guardar em memória o status do device
			status.SetDeviceStatus(message.Id, message.Status != "not connected" && message.Status != "disconnected")
		}
	})
}

func SubscribeMQTTChanges() {
	LocalMQTTClient.Subscribe("gateway", 0, func(client mqtt.Client, msg mqtt.Message) {
		message, err := messages.DecodeMessage(msg.Payload())
		if err != nil {
			log.Printf("Error decoding message: %s\n", err)
		}

		if message.Protocol == "mqtt" {
			InitializeExternalBrokerConnection()
		}
	})
}

func InitializeExternalBrokerConnection() {

	var mymqtt models.MQTTSettings

	for {
		// get MQTTSettings record from database
		result := database.DB.Find(&mymqtt)

		// if there is a MQTTSettings record in the database
		if result.RowsAffected > 0 {
			// if not already connected, connects to external MQTT server
			if GlobalMQTTClient == nil || !GlobalMQTTClient.IsConnected() {
				optsl := mqtt.NewClientOptions().AddBroker("tcp://" + mymqtt.Server + ":" + mymqtt.Port)
				GlobalMQTTClient = mqtt.NewClient(optsl)
				log.Println("Connecting to global MQTT Broker on tcp://" + mymqtt.Server + ":" + mymqtt.Port)

				if token := GlobalMQTTClient.Connect(); token.Wait() && token.Error() != nil {
					log.Printf("Error connecting to global MQTT Broker on tcp://%s:%s: %v\n", mymqtt.Server, mymqtt.Port, token.Error())
				}

				time.Sleep(30 * time.Second)
			}
		} else {
			// if there is no MQTTSettings record in the database, tries again in 5 seconds
			log.Println("=====> WARNING: Could not find MQTTSettings record in database")
			time.Sleep(5 * time.Second)
		}
	}
}

func InitializeLocalBrokerConnection() bool {

	optsl := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	LocalMQTTClient = mqtt.NewClient(optsl)
	log.Println("Connecting to local MQTT Broker on tcp://127.0.0.1:1883")

	if token := LocalMQTTClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to local MQTT Broker on tcp://127.0.0.1:1883: %v\n", token.Error())
		return false
	}

	return true
}

func Publish(payload string) bool {
	token := LocalMQTTClient.Publish("gateway", 0, false, payload)
	token.Wait()
	return token.Error() == nil
}

func VerifyBrokerConnection() bool {
	return GlobalMQTTClient.IsConnected()
}
