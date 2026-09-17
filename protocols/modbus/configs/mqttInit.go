package configs

import (
	"fmt"
	"os"
	"time"

	"icts/modbus/configs/logger"
	packagesPkg "icts/modbus/pkg"
	utilsPkg "icts/modbus/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions
var MqttBrokerInfo utilsPkg.ConfigMQTT
var Hostname string

func SetMqttBroker() {
	logger.Log.Debug("mqttInit:SetMqttBroker: Starting configuration of global variables.")
	MqttBrokerInfo = utilsPkg.ConfigMQTT{
		Server:   utilsPkg.MqttAddress,
		Port:     utilsPkg.MqttPort,
		Username: utilsPkg.MqttUser,
		Password: utilsPkg.MqttPassword,
	}
	Hostname, _ = os.Hostname()

	utilsPkg.CreateChannel()

}

func MqttCommunication() error {

	if MqttBrokerInfo.Server == "" || MqttBrokerInfo.Port == "" {
		logger.Log.Errorf("mqttInit:MqttComunication: Server[%s] or Port[%s] is empty.",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)

	} else {
		logger.Log.Infof("mqttInit:MqttCommunication: Starting communication with Broker MQTT[%s:%s]",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			inMsg, err := packagesPkg.MsgMQTTInput(msg, utilsPkg.MqttClient)
			if err != nil {
				logger.Log.Errorf("mqttInit:MqttCommunication: Failed to receive MQTT message. Device:[%s:%s] Error:[%s]",
					MqttBrokerInfo.Server, MqttBrokerInfo.Port, err)
			} else {
				DecodeMsg(inMsg)
			}
		}

		MqttOptions = mqtt.NewClientOptions()
		MqttOptions.AddBroker(MqttBrokerInfo.Server + ":" + MqttBrokerInfo.Port)
		MqttOptions.SetClientID(Hostname + uuid.New().String())
		MqttOptions.SetUsername(MqttBrokerInfo.Username)
		MqttOptions.SetPassword(MqttBrokerInfo.Password)

		MqttOptions.OnConnect = func(c mqtt.Client) {
			//Subscriber in error and write channels
			if MqttToken = c.Subscribe("general", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
				logger.Log.Errorf("Error in MQTT Token [%v]", MqttToken.Error())
			}
			// if MqttToken = c.Subscribe("modbus", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
			// 	logger.Log.Errorf("Error in MQTT Token [%v]", MqttToken.Error())
			// }
		}

		utilsPkg.MqttClient = mqtt.NewClient(MqttOptions)
		MqttToken = utilsPkg.MqttClient.Connect()

		if MqttToken.Wait() && MqttToken.Error() != nil {
			logger.Log.Errorf("mqttInit:MqttCommunication: Fail to connect with Broker Ip[%s]:Port[%s]",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
			return MqttToken.Error()
			//TODO: Verificar como ficar tentando a conexão por um tempo
		} else {
			logger.Log.Infof("mqttInit:MqttCommunication: MQTT Client is Connected to MQTT Brocker Ip[%s]:Port[%s]",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
		}

		return nil
	}
	return fmt.Errorf("MissingInformations")

}

func RunAction(msg utilsPkg.MessageInput) {
	action, err := msg.Data.(map[string]interface{})["action"].(string)
	if !err {
		logger.Log.Errorf("mqttInit:RunAction: Unable to read message format. Error [%v]",
			err)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Unable to read input message format from MQTT. Error [%v]", err),
				Code:     "",
				Source:   "mqttInit:RunAction:",
			},
		}
		packagesPkg.Sender(utilsPkg.MsgLog, utilsPkg.Log)
		return
	}

	protocol, err := msg.Data.(map[string]interface{})["protocol"].(string)
	if !err {
		logger.Log.Errorf("mqttInit:RunAction: Unable to confirm protocol. Error [%v]",
			err)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf(" Unable to confirm protocol. Error [%v]", err),
				Code:     "",
				Source:   "mqttInit:RunAction:",
			},
		}
		packagesPkg.Sender(utilsPkg.MsgLog, utilsPkg.Log)
		return
	}

	if protocol == "modbus" {
		switch action {
		//Start the read if the s
		case "start":
			//packagesPkg.StartReadFile()
			for _, device := range utilsPkg.OrigiNalDataID {
				status := utilsPkg.DeviceWorkStatus[device.ID]
				if status == "stopped" {
					StartIndividualDevice(device)
				} else {
					logger.Log.Infof("mqttInit:RunAction: The Device ID[%v] is [%s].", device.ID,
								status)
					
					msgLog := utilsPkg.MessageLog{
						MessageType: "info",
						Data: utilsPkg.DataLog{
							Address:  device.Address,
							Port:     device.Port,
							Name:     device.Name,
							Protocol: device.Protocol,
							Time:     time.Now().Format("2006-01-02 15:04:05"),
							Message:  "The Device is already running.",
							Code:     "",
							Source:   "mqttInit:RunAction",
						},
					}
					_ = packagesPkg.Sender(msgLog, utilsPkg.Log)
				}
			}

		case "restart":
			packagesPkg.StartReadFile()
			logger.Log.Info("mqttInit:RunAction: Restart Modbus.")

			for _, device := range utilsPkg.OrigiNalDataID {
				if utilsPkg.DeviceWorkStatus[device.ID] != "stopped" {
					close(utilsPkg.StopChannels[device.ID])
					//delete(utilsPkg.StopChannels, device.ID)

					msgLog := utilsPkg.MessageLog{
						MessageType: "info",
						Data: utilsPkg.DataLog{
							Address:  device.Address,
							Port:     device.Port,
							Name:     device.Name,
							Protocol: device.Protocol,
							Time:     time.Now().Format("2006-01-02 15:04:05"),
							Message:  "Restarting device.",
							Code:     "",
							Source:   "mqttInit:RunAction",
						},
					}
					_ = packagesPkg.Sender(msgLog, utilsPkg.Log)
				}
			}

			//Wait all go routines stop
			for _, wg := range utilsPkg.WaitGroups {
				wg.Wait()
			}

			for _, device := range utilsPkg.OrigiNalDataID {
				delete(utilsPkg.StopChannels, device.ID)
			}			
			for _, device := range utilsPkg.OrigiNalDataID {
				if utilsPkg.DeviceWorkStatus[device.ID] == "stopped" {
					StartIndividualDevice(device)
				}
			}

		case "stop":
			logger.Log.Info("mqttInit:RunAction: Stop Message Receive")
			for _, device := range utilsPkg.OrigiNalDataID {
				if utilsPkg.DeviceWorkStatus[device.ID] != "stopped" {
					close(utilsPkg.StopChannels[device.ID])
					logger.Log.Infof("mqttInit:RunAction: Send Stop to Go routine [%v]", device.ID) 					
				} else {
					logger.Log.Infof("mqttInit:RunAction: The go routine [%v] has stopped.", device.Name)
				}
			}
			//Wait all go routines stop
			for _, wg := range utilsPkg.WaitGroups {
				wg.Wait()
			}
			//Delete all channels
			for _, device := range utilsPkg.OrigiNalDataID {
				delete(utilsPkg.StopChannels, device.ID)
			}
			time.Sleep(100 * time.Millisecond)

		case "update":
			changes := packagesPkg.DiffDevices()
			if len(changes) > 0 {
				for _, change := range changes {
					switch change.Type {
					case utilsPkg.New:
						logger.Log.Infof("mqttInit:RunAction: New device detected: ID [%d]", change.Device.ID)

						msgLog := utilsPkg.MessageLog{
							MessageType: "info",
							Data: utilsPkg.DataLog{
								Address:  change.Device.Address,
								Port:     change.Device.Port,
								Name:     change.Device.Name,
								Protocol: change.Device.Protocol,
								Time:     time.Now().Format("2006-01-02 15:04:05"),
								Message:  "New device detected.",
								Code:     "",
								Source:   "mqttInit:RunAction",
							},
						}
						_ = packagesPkg.Sender(msgLog, utilsPkg.Log)

						StartIndividualDevice(change.Device)

					case utilsPkg.Deleted:
						logger.Log.Infof("mqttInit:RunAction: Device deleted: ID [%d]", change.Device.ID)
						if utilsPkg.DeviceWorkStatus[change.Device.ID] != "stopped" {
							close(utilsPkg.StopChannels[change.Device.ID])		//Request channel close
							utilsPkg.WaitGroups[change.Device.ID].Wait()		//Wait this channel is closed
							delete(utilsPkg.StopChannels, change.Device.ID)		//Delete the channel

							msgLog := utilsPkg.MessageLog{
								MessageType: "info",
								Data: utilsPkg.DataLog{
									Address:  change.Device.Address,
									Port:     change.Device.Port,
									Name:     change.Device.Name,
									Protocol: change.Device.Protocol,
									Time:     time.Now().Format("2006-01-02 15:04:05"),
									Message:  "Device deleted.",
									Code:     "",
									Source:   "mqttInit:RunAction",
								},
							}
							_ = packagesPkg.Sender(msgLog, utilsPkg.Log)
						}

					case utilsPkg.Updated:
						logger.Log.Infof("mqttInit:RunAction: Device updated: ID [%d]", change.Device.ID)
						
						if utilsPkg.DeviceWorkStatus[change.Device.ID] != "stopped" {
							msgLog := utilsPkg.MessageLog{
								MessageType: "info",
								Data: utilsPkg.DataLog{
									Address:  change.Device.Address,
									Port:     change.Device.Port,
									Name:     change.Device.Name,
									Protocol: change.Device.Protocol,
									Time:     time.Now().Format("2006-01-02 15:04:05"),
									Message:  "Device updated.",
									Code:     "",
									Source:   "mqttInit:RunAction",
								},
							}
							_ = packagesPkg.Sender(msgLog, utilsPkg.Log)
							close(utilsPkg.StopChannels[change.Device.ID])		//Request channel close
							utilsPkg.WaitGroups[change.Device.ID].Wait()		//Wait this channel is closed
							delete(utilsPkg.StopChannels, change.Device.ID)		//Delete the channel
							StartIndividualDevice(change.Device)
						} else {
							if IsClosed(utilsPkg.StopChannels[change.Device.ID]){
								delete(utilsPkg.StopChannels, change.Device.ID)
								StartIndividualDevice(change.Device)
							} else {
								close(utilsPkg.StopChannels[change.Device.ID])
								utilsPkg.WaitGroups[change.Device.ID].Wait()
								delete(utilsPkg.StopChannels, change.Device.ID)
								StartIndividualDevice(change.Device)
							}
						}
					}
				}
			} else {
				logger.Log.Info("mqttInit:RunAction: No updates identified")
				msgLog := utilsPkg.MessageLog{
					MessageType: "info",
					Data: utilsPkg.DataLog{
						Protocol: "modbus",
						Time:     time.Now().Format("2006-01-02 15:04:05"),
						Message:  "No updates identified.",
						Code:     "",
						Source:   "mqttInit:RunAction",
					},
				}
				packagesPkg.Sender(msgLog, utilsPkg.Log)
			}
		
		default:
			logger.Log.Warn("mqttInit:RunAction: Wrong Command")

			msgLog := utilsPkg.MessageLog{
				MessageType: "error",
				Data: utilsPkg.DataLog{
					Address:  "",
					Port:     "",
					Name:     "",
					Protocol: "modbus",
					Time:     time.Now().Format("2006-01-02 15:04:05"),
					Message:  "Wrong comand frum MQTT action",
					Code:     "",
					Source:   "mqttInit:RunAction:",
				},
			}
			packagesPkg.Sender(msgLog, utilsPkg.Log)
		}
	}
}

func DecodeMsg(msg utilsPkg.MessageInput) {
	switch msg.MessageType {
	case "action":
		logger.Log.Info("mqttInit:DecodeMsg: Process received message.")

		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Protocol: "modbus",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Process received message.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg:",
			},
		}
		_ = packagesPkg.Sender(msgLog, utilsPkg.Log)
		//time.Sleep(500 * time.Millisecond)
		RunAction(msg)
	}
}

func IsClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		fmt.Println("CANAL ABERTO")
		return false
	}
}
