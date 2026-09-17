package mqtt

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	packagesPkg "newprotocol/protocol" /// 14/10/2023 ADDED FROM HOLDEN NOT USED
	utilsPkg "newprotocol/utils"
	start "newprotocol/utils/config"
	logPkg "newprotocol/utils/gologtofile"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	mqttServer "github.com/mochi-mqtt/server/v2" //28/11/2023 Added to start a MQTT Brocker if it was not running
	"github.com/mochi-mqtt/server/v2/hooks/auth" //28/11/2023 Added to start a MQTT Brocker if it was not running
	"github.com/mochi-mqtt/server/v2/listeners"  //28/11/2023 Added to start a MQTT Brocker if it was not running

	"github.com/google/uuid"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

var MqttBrokerInfo utilsPkg.MqttBrockerStruct

func SetMqttBroker() {
	logPkg.CtsLog.Info("Starting assign values to public variables.")
	MqttBrokerInfo = utilsPkg.MqttBrockerStruct{
		Server:   utilsPkg.LocalMqttBrockerAddress,
		Port:     utilsPkg.LocalMqttBrockerPort,
		Username: utilsPkg.LocalMqttBrockerUserName,
		Password: utilsPkg.LocalMqttBrockerPassword,
	}
}

func ConnectToLocalMqttBrocker() error {
	SetMqttBroker()

	/// Set Local Mqtt Brocker Info from Globals
	MqttBrokerInfo = utilsPkg.MqttBrockerStruct{
		Server:   utilsPkg.LocalMqttBrockerAddress,
		Port:     utilsPkg.LocalMqttBrockerPort,
		Username: utilsPkg.LocalMqttBrockerUserName,
		Password: utilsPkg.LocalMqttBrockerPassword,
	}

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		inMsg, err := MsgMQTTInput(msg, utilsPkg.MqttClient)
		if err != nil {
			logPkg.CtsLog.Error("SetLocalMqttBrokerInfo:MsgMQTTInput\n          OS[%s:%s] %d bits\n    HostName[%s]\n     ip:port[%s:%s]\n     usr:pwd[%s:%s]\n         msg[%s]\n",
				runtime.GOOS,
				runtime.GOARCH,
				utilsPkg.OsBits,
				utilsPkg.Hostname,
				utilsPkg.LocalMqttBrockerAddress,
				utilsPkg.LocalMqttBrockerPort,
				utilsPkg.LocalMqttBrockerUserName,
				utilsPkg.LocalMqttBrockerPassword,
				msg)
		} else {
			//inMsg.MessageType = "data"
			logPkg.CtsLog.Debug("SetLocalMqttBrokerInfo:MsgMQTTInput\n          inMsg[%v]", inMsg)
			DecodeMsg(inMsg)
		}
	}

	MqttOptions = mqtt.NewClientOptions()
	MqttOptions.AddBroker(MqttBrokerInfo.Server + ":" + MqttBrokerInfo.Port)
	MqttOptions.SetClientID(utilsPkg.Hostname + uuid.New().String())
	MqttOptions.SetUsername(MqttBrokerInfo.Username)
	MqttOptions.SetPassword(MqttBrokerInfo.Password)

	MqttOptions.OnConnect = func(c mqtt.Client) {
		//Subscribe in Topic: canopen
		//if MqttToken = c.Subscribe("canopen", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
		if MqttToken = c.Subscribe("general", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
			logPkg.CtsLog.Error("SetLocalMqttBrokerInfo:subscribe[canopen] MQTT Brocker[%s:%s]Not Subscribed\n",
				utilsPkg.LocalMqttBrockerAddress,
				utilsPkg.LocalMqttBrockerPort)
		} else {
			logPkg.CtsLog.Warn("SetLocalMqttBrokerInfo:subscribe[canopen] MQTT Brocker[%s:%s]Subscribed\n",
				utilsPkg.LocalMqttBrockerAddress,
				utilsPkg.LocalMqttBrockerPort)
		}
	}

	// Connecting to Internal MQTT Brocker
	utilsPkg.MqttClient = mqtt.NewClient(MqttOptions)
	MqttToken = utilsPkg.MqttClient.Connect()

	if MqttToken.Wait() && MqttToken.Error() != nil {
		logPkg.CtsLog.Error("SetLocalMqttBrokerInfo:MQTT Client Fail to connect to Internal MQTT Brocker[%s:%s]\n",
			utilsPkg.LocalMqttBrockerAddress,
			utilsPkg.LocalMqttBrockerPort)
		return MqttToken.Error()
	}
	logPkg.CtsLog.Warn("SetLocalMqttBrokerInfo:MQTT Client Connected Internal MQTT Brocker[%s:%s]Connected!!!\n",
		utilsPkg.LocalMqttBrockerAddress,
		utilsPkg.LocalMqttBrockerPort)
	return nil
}

// Decoding Rcvd MQTT JSON Nessage
func DecodeMsg(msg utilsPkg.MessageInput) {
	var status string = ""
	// Locating [action] on JSON Message
	action, err := msg.Data.(map[string]interface{})["action"].(string)
	if !err {
		// Locating [action] on JSON Message
		status, err = msg.Data.(map[string]interface{})["status"].(string)
		if !err {
			logPkg.CtsLog.Warn("msg.MessageType[%s] msg.Data[%v] err[%v]Expected find action or status on JSON MsgData\n", msg.MessageType, msg.Data, err)
		}
		// Status Found!!!
	}
	// Action or status Found!!!

	// Locating [protocol] on JSON Message
	protocol, err := msg.Data.(map[string]interface{})["protocol"].(string)
	if !err {
		if msg.MessageType != "status" {
			logPkg.CtsLog.Warn("msg.MessageType[%s] status[%s] action[%s] msg.Data[%v]Expected find protocol  JSON MsgData\n", msg.MessageType, status, action, msg.Data)
		} else {
			logPkg.CtsLog.Warn("msg.MessageType[%s] status[%s] action[%s] msg.Data[%v]Found status on JSON MsgData\n", msg.MessageType, status, action, msg.Data)
		}
	}
	if protocol != "canopen" {
		return
	}
	if msg.MessageType == "action" {
		logPkg.CtsLog.Warn("msg.MessageType[%s] msg.Data[%v]Processing...\n", msg.MessageType, msg.Data)
	} else {
		// testing other msg.MessageType [status or deviceStatus]
		if msg.MessageType == "status" {
			logPkg.CtsLog.Warn("msg.MessageType[%s] msg.Data[%v]OK\n", msg.MessageType, msg.Data)
		} else if msg.MessageType == "deviceStatus" {
			logPkg.CtsLog.Warn("msg.MessageType[%s] msg.Data[%v]OK\n", msg.MessageType, msg.Data)
		} else if msg.MessageType != "" {
			// msg.MessageType is not [action, status or deviceStatus]
			logPkg.CtsLog.Error("msg.MessageType[%s]Expected Action e, msg.Data[%v]\n", msg.MessageType, msg.Data)
		} else {
			logPkg.CtsLog.Warn("msg.MessageType[%s]<<<<< Expected action,status or deviceStatus >>>> msg.Data[%v]\n", msg.MessageType, msg.Data)
		}
		return
	}
	logPkg.CtsLog.Warn("msg.MessageType[%s] action[%s] msg.Data[%v]\n", msg.MessageType, action, msg.Data)

	// Action and Protocol Found!!!
	if protocol == "canopen" {
		// Action, Protocol Found and is correct protocol: Proccess Action now...!!!
		switch action {

		case "stop":
			// Action:Stop
			if utilsPkg.AnyDevicelsRunning {
				// Action:Stop when Go Routines are running
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Stopping...\n", msg.MessageType, protocol, action)
				start.StopDevicesProtocol()
			} else {
				// Action:Stop when already stopped - no Go Routines are running
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Already Stopped\n", msg.MessageType, protocol, action)
			}

		case "start":
			// Action:Start
			if !utilsPkg.AnyDevicelsRunning {
				// Action:Start when previously stopped
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Starting...\n", msg.MessageType, protocol, action)
				start.StartProtocol()
			} else {
				// Action:Start when Already Started: All GoRoutines are running: do nothing...
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Already Running!!!\n", msg.MessageType, protocol, action)
			}

		case "restart":
			// Action:Restart
			start.RestartDevicesProtocol()
			if utilsPkg.AnyDevicelsRunning {
				// Restart: When Go Routineas are running >> Stop, wait 1s second then Start
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Restart:Stop, wait 1s then Start\n", msg.MessageType, protocol, action)
				start.StopDevicesProtocol()         // STEP01: Stop
				time.Sleep(1000 * time.Millisecond) // STEP02: Wait 1s
				start.StartProtocol()               // STEP03: Start
			} else {
				// Restart: When previouly stopped  >> Just Start
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Restart:Previously Stopped just Start\n", msg.MessageType, protocol, action)
				// time.Sleep(1000 * time.Millisecond) // TODO: Check if it is required
				start.StartProtocol() // Starting
			}

		case "update":
			// Action:Update - Update >> Read deviceConfig.json >>  Restart
			logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Updating...\n", msg.MessageType, protocol, action)
			if !utilsPkg.DevicelUpdateRequest {
				utilsPkg.DevicelUpdateRequest = true
			}
			start.UpdateDevicesProtocol()
			start.RestartDevicesProtocol()
			if utilsPkg.AnyDevicelsRunning {
				// Restart: When Go Routineas are running >> Stop, wait 1s second then Start
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Restart:Stop, wait 1s then Start\n", msg.MessageType, protocol, action)
				start.StopDevicesProtocol()         // STEP01: Stop
				time.Sleep(1000 * time.Millisecond) // STEP02: Wait 1s
				start.StartProtocol()               // STEP03: Start
			} else {
				// Restart: When previouly stopped  >> Just Start
				logPkg.CtsLog.Warn("msg.MessageType[%s] protocol[%s] action[%s]Restart:Previously Stopped just Start\n", msg.MessageType, protocol, action)
				// time.Sleep(1000 * time.Millisecond) // TODO: Check if it is required
				start.StartProtocol() // Starting
			}

		default:
			// Action:Unknown
			logPkg.CtsLog.Debug("msg.MessageType[%s] protocol[%s] action[%s]UNKNOWN\n", msg.MessageType, protocol, action)
		}
	} else {
		// Action & Protocol found but is from other protocol - ignore
		logPkg.CtsLog.Debug("msg.MessageType[%s] action[%s] protocol[%s]Expected[canopen] msg.Data[%v]ignore...\n", msg.MessageType, action, protocol, msg.Data)
	}
}

// / 14/10/2023 BEGIN ADDED FROM HOLDEN - NOT USED
func RunAction(msg utilsPkg.MessageInput) {
	action, err := msg.Data.(map[string]interface{})["action"].(string)
	if !err {
		logPkg.CtsLog.Error("mqttInit:RunAction: Unable to read message format. Error [%v]",
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
		logPkg.CtsLog.Debug("mqttInit:RunAction: Unable to confirm protocol. Error [%v]",
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
			for _, dev := range utilsPkg.OrigiNalDataID {
				if !utilsPkg.DeviceWorkStatus[dev.ID] {
					StartIndividualDevice(dev)
				} else {
					logPkg.CtsLog.Debug("mqttInit:RunAction: The Device ID[%v] is already running.", dev.ID)

					msgLog := utilsPkg.MessageLog{
						MessageType: "info",
						Data: utilsPkg.DataLog{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
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
			logPkg.CtsLog.Debug("mqttInit:RunAction: Restart Modbus.")

			for _, dev := range utilsPkg.OrigiNalDataID {
				if utilsPkg.DeviceWorkStatus[dev.ID] {
					close(utilsPkg.StopChannels[dev.ID])
					//delete(utilsPkg.StopChannels, dev.ID)

					msgLog := utilsPkg.MessageLog{
						MessageType: "info",
						Data: utilsPkg.DataLog{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
							Time:     time.Now().Format("2006-01-02 15:04:05"),
							Message:  "Restarting dev.",
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
			for _, dev := range utilsPkg.OrigiNalDataID {
				if !utilsPkg.DeviceWorkStatus[dev.ID] {
					StartIndividualDevice(dev)
				}
			}

		case "stop":
			logPkg.CtsLog.Debug("mqttInit:RunAction: Stop Message Receive")
			for _, dev := range utilsPkg.OrigiNalDataID {
				if utilsPkg.DeviceWorkStatus[dev.ID] {
					close(utilsPkg.StopChannels[dev.ID])
					logPkg.CtsLog.Debug("mqttInit:RunAction: Send Stop to Go routine [%v]", dev.ID) // <- true
					//delete(utilsPkg.StopChannels, dev.ID)
				} else {
					logPkg.CtsLog.Debug("mqttInit:RunAction: The go routine [%v] has stopped.", dev.Name)
				}
			}
			//Wait all go routines stop
			for _, wg := range utilsPkg.WaitGroups {
				wg.Wait()
			}

			for _, dev := range utilsPkg.OrigiNalDataID {
				delete(utilsPkg.StopChannels, dev.ID)
			}

			time.Sleep(500 * time.Millisecond)

		case "update":
			changes := packagesPkg.DiffDevices()
			if len(changes) > 0 {
				for _, change := range changes {
					switch change.Type {
					case utilsPkg.New:
						logPkg.CtsLog.Debug("mqttInit:RunAction: New device detected: ID [%d]", change.Device.ID)

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
						logPkg.CtsLog.Debug("mqttInit:RunAction: Device deleted: ID [%d]", change.Device.ID)
						if utilsPkg.DeviceWorkStatus[change.Device.ID] {
							close(utilsPkg.StopChannels[change.Device.ID])
							delete(utilsPkg.StopChannels, change.Device.ID)

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
						logPkg.CtsLog.Debug("mqttInit:RunAction: Device updated: ID [%d]", change.Device.ID)
						if utilsPkg.DeviceWorkStatus[change.Device.ID] {

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
							close(utilsPkg.StopChannels[change.Device.ID])
							//delete(utilsPkg.StopChannels, change.dev.ID)
							utilsPkg.WaitGroups[change.Device.ID].Wait()
							StartIndividualDevice(change.Device)
						}
					}
				}
			} else {
				logPkg.CtsLog.Debug("mqttInit:RunAction: No updates identified")
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
			logPkg.CtsLog.Warn("mqttInit:RunAction: Wrong Command")

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

func StartIndividualDevice(dev utilsPkg.Device) {
	logPkg.CtsLog.Debug("modbusInit:StartRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]",
		dev.Name, dev.Protocol, dev.Address, dev.Port)

	wg := &sync.WaitGroup{} // CANopenTcpServer_Create a new WaitGroup for this routine
	utilsPkg.WaitGroups[dev.ID] = wg
	wg.Add(1)

	utilsPkg.StopChannels[dev.ID] = make(chan struct{})
	go func(device utilsPkg.Device) {
	}(dev)

	time.Sleep(100 * time.Millisecond)

}

/// 14/10/2023 END ADDED FROM HOLDEN - NOT USED

// Begin 28/11/2022 Added by Oliveira to Start Local MQTT Brocker if it is not running it to test protocol standalone
func StartLocalMQTTServer() error {
	logPkg.CtsLog.Warn("StartLocalMQTTServer: Bein")

	// Create the new MQTT Server.
	server := mqttServer.New(nil)

	// Allow all connections.
	err := server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		logPkg.CtsLog.Error("StartLocalMQTTServer: Local Mqtt Brocker err[%s]", err)
		return err
	}

	// Create a TCP listener on a standard port.
	localMqttBrockerTcp := listeners.NewTCP("t1", ":1883", nil)
	err = server.AddListener(localMqttBrockerTcp)
	if err != nil {
		logPkg.CtsLog.Error("StartLocalMQTTServer: Local Mqtt Brocker[%s] err[%s]", localMqttBrockerTcp.Address(), err)
		return err
	}

	go func() {
		err := server.Serve()
		if err != nil {
			logPkg.CtsLog.Error("StartLocalMQTTServer: Local Mqtt Brocker[%s] err[%s]", localMqttBrockerTcp.Address(), err)
		} else {
			logPkg.CtsLog.Warn("StartLocalMQTTServer: Local Mqtt Brocker[%s] Running...", localMqttBrockerTcp.Address())
		}
	}()
	return err
}

// End 28/11/2022 Added by Oliveira to Start Local MQTT Brocker if it is not running it to test protocol standalone
