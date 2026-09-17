package mqtt

import (
	"encoding/json"
	ethernetipPkg "ethernetip/internal/ethip"
	utilsPkg "ethernetip/utils"
	startEth "ethernetip/utils/config"
	golog "ethernetip/utils/gologtofile"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

var MqttBrokerInfo utilsPkg.ConfigMQTT
var Hostname string

func SetMqttBroker() {
	fmt.Println("Starting assign values to public variables.")
	MqttBrokerInfo = utilsPkg.ConfigMQTT{
		Server:   utilsPkg.MqttAddress,
		Port:     utilsPkg.MqttPort,
		Username: utilsPkg.MqttUser,
		Password: utilsPkg.MqttPassword,
	}

	Hostname, _ = os.Hostname()
}

func MqttCommunication() error {

	if MqttBrokerInfo.Server == "" || MqttBrokerInfo.Port == "" {
		golog.CtsLog.Info("mqttInit:MqttComunication: Server[%s] or Port[%s] is empty. \n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)
	} else {
		fmt.Printf("mqttInit:MqttCommunication: Starting communication with Broker MQTT[%s:%s]\n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			inMsg, err := MsgMQTTInput(msg, utilsPkg.MqttClient)
			if err != nil {
				golog.CtsLog.Info("mqttInit:MqttCommunication: Receiving MQTT message failed. Device:[%s:%s] Error:[%s]\n",
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
			//	if MqttToken = c.Subscribe("ethernetip", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
			if MqttToken = c.Subscribe("general", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
				golog.CtsLog.Info("%v", MqttToken.Error())
			}
		}

		utilsPkg.MqttClient = mqtt.NewClient(MqttOptions)
		MqttToken = utilsPkg.MqttClient.Connect()

		if MqttToken.Wait() && MqttToken.Error() != nil {
			golog.CtsLog.Info("mqttInit:MqttCommunication: Fail to connect with Broker Ip:Port(%s:%s)\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
			return MqttToken.Error()
			//TODO: Verificar como ficar tentando a conexão por um tempo
		} else {
			golog.CtsLog.Info("mqttInit:MqttCommunication: MQTT Client is Connected to MQTT Broker Ip:Port[%s:%s]\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
		}

		return nil
	}
	return fmt.Errorf("Missing_Informations")

}

func DecodeMsg(msg utilsPkg.MessageInput) {
	if msg.MessageType == "action" {
		golog.CtsLog.Info("mqttInit:DecodeMsg: Process received message.\n")

		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Process received message.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg:",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")
	} else {
		//golog.CtsLog.Info("Comando desconhecido!\n")
		return
	}
	action, err := msg.Data.(map[string]interface{})["action"].(string)
	if !err {
		golog.CtsLog.Error("mqttInit:RunAction: Unable to read message format. Error [%v]\n",
			err)

		msgLog := utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Unable to read input message format from MQTT. Error [%v]\n", err),
				Code:     "",
				Source:   "mqttInit:RunAction:",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")

		return
	}

	protocol, err := msg.Data.(map[string]interface{})["protocol"].(string)
	if !err {
		golog.CtsLog.Error("mqttInit:RunAction: Unable to confirm protocol. Error [%v]\n",
			err)

		msgLog := utilsPkg.MessageLog{
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
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")

		return
	}
	if protocol == "ethernetip" {
		switch action {
		case "start":
			StartProtocol()
		case "update":
			UpdateDevicesProtocol()
		case "restart":
			RestartProtocol()
		case "stop":
			StopProtocol()
		default:
			fmt.Println("Wrong Command")
		}
	}
}

func StartingLocalBroker() error {
	SetMqttBroker()
	err := MqttCommunication()
	return err
}

func RestartProtocol() {

	golog.CtsLog.Info("All devices stopped reading memories.\n")
	msgLog := utilsPkg.MessageLog{
		MessageType: "info",
		Data: utilsPkg.DataLog{
			Address:  "",
			Port:     "",
			Name:     "",
			Protocol: "ethernetip",
			Time:     time.Now().Format("2006-01-02 15:04:05"),
			Message:  "All devices will be restarted reading memories.",
			Code:     "",
			Source:   "mqttInit:DecodeMsg",
		},
	}
	_ = ethernetipPkg.SendMsgLog(msgLog, "log")
	StopProtocol()
	time.Sleep(2 * time.Second)
	StartProtocol()
}

func StopProtocol() {
	if utilsPkg.StatusProtocol {
		close(utilsPkg.StopChan)
		utilsPkg.Wg.Wait()
		golog.CtsLog.Info("All devices stopped reading memories.\n")
		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "ethernetip",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "All devices stopped reading memories.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")
	} else {
		golog.CtsLog.Info("All devices already closed\n")
	}
}

func StartProtocol() {
	if !utilsPkg.StatusProtocol {
		golog.CtsLog.Info("Start Ethernetip\n")
		startEth.StartEthernetip()
	} else {
		golog.CtsLog.Info("The devices is already running\n")
		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "ethernetip",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "The Device is already running.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")
	}
}

func UpdateDevicesProtocol() {
	if utilsPkg.StatusProtocol {
		if startEth.CheckNewDevices(utilsPkg.DevicesInUse) {
			StopProtocol()
			StartProtocol()
			return
		}
		utilsPkg.UpdateStatusCount = utilsPkg.NumberDevice
		close(utilsPkg.UpdateDevChan)
		for utilsPkg.UpdateStatusCount != 0 {
			time.Sleep(10 * time.Millisecond)
		}
		golog.CtsLog.Info("Updating devices Ethernetip\n")
		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "ethernetip",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "The Device is already updated.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")
		utilsPkg.UpdateDevChan = make(chan struct{})
		//startEth.StartEthernetip()
	} else {
		golog.CtsLog.Info("No devices Ethernetip to update.\n")
		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  "",
				Port:     "",
				Name:     "",
				Protocol: "ethernetip",
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "No devices Ethernetip to update.",
				Code:     "",
				Source:   "mqttInit:DecodeMsg",
			},
		}
		_ = ethernetipPkg.SendMsgLog(msgLog, "log")
	}
}

func MsgMQTTInput(msg mqtt.Message, client mqtt.Client) (utilsPkg.MessageInput, error) {

	var msgInf utilsPkg.MessageInput
	payload := msg.Payload()
	err := json.Unmarshal(payload, &msgInf)

	if err != nil {
		fmt.Printf("mqttInput:MsgMQTTInput: Fail to decode JSON Input Msg MQTT - Data[%s] err[%v]\n",
			msg.Payload(), err)
		return msgInf, err
	}
	//fmt.Printf("mqttInput:MsgMQTTInput: msgEthernetip.Protocol    [%s]\n", msgInf)
	return msgInf, err
}
