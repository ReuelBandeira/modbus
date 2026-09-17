package ethip

import (
	"encoding/json"
	"fmt"

	utilsPkg "ethernetip/utils"
	golog "ethernetip/utils/gologtofile"
	"net"
	"os"
	"time"
)

type MemoryEthInfo struct {
	Name      string `json:"name"`
	Class     uint16 `json:"class"`
	Instance  uint16 `json:"instance"`
	Attribute uint16 `json:"attribute"`
	Value     uint32 `json:"value"`
}

type PlcEthInfo struct {
	BitMemories  []MemoryEthInfo `json:"bitMemories"`
	WordMemories []MemoryEthInfo `json:"wordMemories"`
}

var PlcConnected bool = true

// func ReadInfoEth(dev utilsPkg.DevSettings, setMqtt utilsPkg.SettingsMQTT, client mqtt.Client, loop bool, mode string) {
func ReadInfoEth(dev utilsPkg.DevSettings) {
	utilsPkg.StatusProtocol = true
	defer utilsPkg.Wg.Done()

	plcInfoEth := changeEthFormat(dev)

	var resultsMemoryBit []byte
	var resultsMemoryWord []uint32
	var readBitMemory byte
	var readMemoryWord uint32
	var status byte
	var firstLoop bool = true
	var resultsGetBitMemories map[string]interface{}
	var resultsGetWordMemories map[string]interface{}

	resultsMemoryBit = make([]byte, len(plcInfoEth.BitMemories))
	resultsMemoryWord = make([]uint32, len(plcInfoEth.WordMemories))

	readHour := time.Now()
	readTime := readHour.Format("2006-01-02 15:04:05")
	sleepTime := dev.ReadingTime
	changeMemory := false

	p := &EipPlc{}

	if !isTCPConnected(dev.Address, dev.Port) {
		golog.CtsLog.Error("routineEthRead:ReadInfoEth: Unable to Connect with Ethernetip device [%s]\n",
			dev.Name)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  dev.Address,
				Port:     dev.Port,
				Name:     dev.Name,
				Protocol: dev.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Unable to Connect with Ethernetip device.",
				Code:     "",
				Source:   "routineEthRead:ReadInfoEth",
			},
		}
		SendMsgLog(utilsPkg.MsgLog, "log")
	} else {
		golog.CtsLog.Info("routineEthRead:ReadInfoEth: Starting device reading [%v]",
			dev.Name)

		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  dev.Address,
				Port:     dev.Port,
				Name:     dev.Name,
				Protocol: dev.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Starting device communication with ethernetip.",
				Code:     "",
				Source:   "routineEthRead:ReadInfoEth",
			},
		}
		SendMsgLog(msgLog, "log")
	}
	for {
		if isTCPConnected(dev.Address, dev.Port) {
			break
		} else {
			golog.CtsLog.Info("routineEthRead: dev name[%s]  Ip[%s:%s] Disconnected - Waiting 2 Seconds...\n", dev.Name, dev.Address, dev.Port)
			time.Sleep(2 * time.Second)
		}
		select {
		case <-utilsPkg.StopChan:
			utilsPkg.StatusProtocol = false
			return
		case <-utilsPkg.UpdateDevChan:
			utilsPkg.UpdateStatusCount--
		default:
		}

	}
	p, _ = ConnEthernetIP(dev.Address, dev.Port)
	for {
		select {
		case <-utilsPkg.StopChan:
			utilsPkg.StatusProtocol = false
			golog.CtsLog.Info("Leitura do IP [%s] interrompida. \n", dev.Name)

			golog.CtsLog.Info("routineEthRead:ReadInfoEth: Stop EtherneitIP Device Name[%v].",
				dev.Name)

			utilsPkg.MsgLog = utilsPkg.MessageLog{
				MessageType: "info",
				Data: utilsPkg.DataLog{
					Address:  dev.Address,
					Port:     dev.Port,
					Name:     dev.Name,
					Protocol: dev.Protocol,
					Time:     time.Now().Format("2006-01-02 15:04:05"),
					Message:  fmt.Sprintf("Stop EthernetIP Device. %s", dev.Name),
					Code:     "",
					Source:   "routineEthRead:ReadInfoEth",
				},
			}
			SendMsgLog(utilsPkg.MsgLog, "log")
			SendStatusDevice(dev, "disconnected")
			if p.tcpConn != nil {
				_ = p.tcpConn.Close()
				p.tcpConn = nil
			}
			return

		case <-utilsPkg.UpdateDevChan:
			devConfig, memChanged, err := UpdateDevConfig(dev.Name, plcInfoEth)
			if err != nil {
				golog.CtsLog.Error("Fail to read JSON File: ", err)
			} else {
				plcInfoEth = changeEthFormat(devConfig)
				//resultsGetBitMemories = make(map[string]interface{})
				//resultsGetWordMemories = make(map[string]interface{})
				resultsMemoryBit = make([]byte, len(plcInfoEth.BitMemories))
				resultsMemoryWord = make([]uint32, len(plcInfoEth.WordMemories))
				golog.CtsLog.Info("Status update das memorias do %s \n", devConfig.Name)
				if len(memChanged) != 0 {
					for i := 0; i < len(memChanged); i++ {
						golog.CtsLog.Info(" %s\n", memChanged[i])
					}
				}
				utilsPkg.UpdateStatusCount--
				time.Sleep(3 * time.Second)
			}

		default:

			utilsPkg.StatusProtocol = true
			resultsGetBitMemories = make(map[string]interface{})
			resultsGetWordMemories = make(map[string]interface{})
			if len(plcInfoEth.BitMemories) != 0 && CheckConnection(p, dev) {
				for chave, itemMemories := range plcInfoEth.BitMemories {
					/*
						if mode == "write" {
							status = p.WriteBitMemory(uint32(itemMemories.Class), uint32(itemMemories.Instance), uint32(itemMemories.Attribute), byte(itemMemories.Value))
							if status != 0 {
								// if status not equal zero, sinalize that one error was reported
								// CIP protocol report error. Publish error using mqtt
								msgError := utilsPkg.MqttErrorMsg{
									Address:  dev.Address,
									Port:     dev.Port,
									Name:     dev.Name,
									Protocol: dev.Protocol,
									Time:     readTime,
									Message:  utilsPkg.CodeCIPerrors[status],
									Code:     "EIP_0" + fmt.Sprintf("%v", status),
									Source:   fmt.Sprint("Error writing memory bit ", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", itemMemories.Value),
								}
								_ = MQTTSendEthErrorMessage(msgError, setMqtt.ErrorInfo)
								golog.CtsLog.Error("Error writing memory bit", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", readMemoryWord)
							}
						}
					*/
					readBitMemory, status = p.ReadBitMemory(uint32(itemMemories.Class), uint32(itemMemories.Instance), uint32(itemMemories.Attribute))
					readBit := false
					if readBitMemory == 1 {
						readBit = true
					}
					if status != 0 {
						// if status not equal zero, sinalize that one error was reported
						// CIP protocol report error. Publish error using mqtt

						utilsPkg.MsgLog = utilsPkg.MessageLog{
							MessageType: "error",
							Data: utilsPkg.DataLog{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Time:     readTime,
								Message:  fmt.Sprintf(utilsPkg.CodeCIPerrors[status]),
								Code:     "EIP_0" + fmt.Sprintf("%v", status),
								Source:   "routineEthRead:ReadInfoEth",
							},
						}
						SendMsgLog(utilsPkg.MsgLog, "log")
						golog.CtsLog.Error("Error reading memory bit", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", readBitMemory)
					}
					if resultsMemoryBit[chave] != readBitMemory || firstLoop {

						//golog.CtsLog.Info(dev.Name, " - ", dev.Address, " Memory Bit", itemMemories.Name, " Value = ", readBit)
						resultsMemoryBit[chave] = readBitMemory
						resultsGetBitMemories[itemMemories.Name] = readBit //fmt.Sprintf("%v", readBitMemory)
						changeMemory = true
					}

				}

			}

			if len(plcInfoEth.WordMemories) != 0 && CheckConnection(p, dev) {
				for chave, itemMemories := range plcInfoEth.WordMemories {
					/*
						if mode == "write" {
							status = p.WriteWordMemory(uint32(itemMemories.Class), uint32(itemMemories.Instance), uint32(itemMemories.Attribute), itemMemories.Value)
							if status != 0 {
								// if status not equal zero, sinalize that one error was reported
								// CIP protocol report error. Publish error using mqtt
								msgError := utilsPkg.MqttErrorMsg{
									Address:  dev.Address,
									Port:     dev.Port,
									Name:     dev.Name,
									Protocol: dev.Protocol,
									Time:     readTime,
									Message:  codeCIPerrors[status],
									Code:     "EIP_0" + fmt.Sprintf("%v", status),
									Source:   fmt.Sprint("Error writing memory word", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", itemMemories.Value),
								}
								_ = MQTTSendEthErrorMessage(msgError, setMqtt.ErrorInfo)
								golog.CtsLog.Error("Error writing memory word", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", readMemoryWord)
							}
						}
					*/
					//analogInputsRead := p.ReadAnalogInputs(uint32(itemMemories.Class), uint32(itemMemories.Instance), uint32(itemMemories.Attribute))

					readMemoryWord, status = p.ReadWordMemory(uint32(itemMemories.Class), uint32(itemMemories.Instance), uint32(itemMemories.Attribute))

					if status != 0 {
						// if status not equal zero, sinalize that one error was reported
						// CIP protocol report error. Publish error using mqtt
						utilsPkg.MsgLog = utilsPkg.MessageLog{
							MessageType: "error",
							Data: utilsPkg.DataLog{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Time:     readTime,
								Message:  fmt.Sprintf(utilsPkg.CodeCIPerrors[status]),
								Code:     "EIP_0" + fmt.Sprintf("%v", status),
								Source:   "routineEthRead:ReadInfoEth",
							},
						}
						SendMsgLog(utilsPkg.MsgLog, "log")
						golog.CtsLog.Error("Error reading memory word", dev.Address, " PLC - Memory ", itemMemories.Name, " Value = ", readMemoryWord)
					}

					if resultsMemoryWord[chave] != readMemoryWord || firstLoop {

						//golog.CtsLog.Info(dev.Name, " - ", dev.Address, " Memory Word", itemMemories.Name, " Value = ", readMemoryWord)
						resultsMemoryWord[chave] = readMemoryWord
						resultsGetWordMemories[itemMemories.Name] = readMemoryWord //fmt.Sprintf("%v", readMemoryWord)
						changeMemory = true
					}
				}

			}
			if changeMemory {
				// Houve alguma alteração gerar payload JSON da resposta
				msg := utilsPkg.ExitPayloadMsg{
					Address:       dev.Address,
					Port:          dev.Port,
					Name:          dev.Name,
					Others:        "",
					Protocol:      dev.Protocol,
					ReadTimeStamp: readTime,
					Topics:        dev.Topics,
					BitMemories:   resultsGetBitMemories,
					WordMemories:  resultsGetWordMemories,
				}
				// Enviar payload json da mensagem
				MQTTSendMessageAutomatic(msg)
				golog.CtsLog.Info("routineEthRead: dev name[%s]  Ip[%s:%s] Changes Occurred - Sleeping %v Milliseconds...\n", dev.Name, dev.Address, dev.Port, sleepTime)
			} else {
				if CheckConnection(p, dev) {
					golog.CtsLog.Info("routineEthRead: dev name[%s]  Ip[%s:%s] No  Changes Occurred - Sleeping %v Milliseconds...\n", dev.Name, dev.Address, dev.Port, sleepTime)
				} else {
					golog.CtsLog.Info("routineEthRead: dev name[%s]  Ip[%s:%s] Disconnected - Sleeping %v Milliseconds...\n", dev.Name, dev.Address, dev.Port, sleepTime)
				}

			}
			if isTCPConnected(dev.Address, dev.Port) {
				if !PlcConnected {
					PlcConnected = true
					SendStatusDevice(dev, "connected")
				}
			} else {
				if PlcConnected {
					PlcConnected = false
					SendStatusDevice(dev, "disconnected")
				}
			}
			changeMemory = false
			firstLoop = false
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
	}
}

func CheckConnection(p *EipPlc, device utilsPkg.DevSettings) bool {

	tcpAddress := device.Address + ":" + device.Port

	if !isTCPConnected(device.Address, device.Port) {
		if p.tcpConn != nil {
			SendStatusDevice(device, "disconnected")
			utilsPkg.MsgLog = utilsPkg.MessageLog{
				MessageType: "error",
				Data: utilsPkg.DataLog{
					Address:  device.Address,
					Port:     device.Port,
					Name:     device.Name,
					Protocol: device.Protocol,
					Time:     time.Now().Format("2006-01-02 15:04:05"),
					Message:  "The connection was terminated due lost communication.",
					Code:     "",
					Source:   "routineEthRead:ReadInfoEth",
				},
			}
			SendMsgLog(utilsPkg.MsgLog, "log")
			p.tcpConn = nil
			golog.CtsLog.Error("continuosReadEthIp: FAIL - PLC lost connection ", tcpAddress)
		}
		return false
	}
	if p.tcpConn == nil {
		mutex.Lock()
		p.Connect()
		mutex.Unlock()
		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Starting device communication with ethernetip.",
				Code:     "",
				Source:   "routineEthRead:ReadInfoEth",
			},
		}
		SendMsgLog(msgLog, "log")
	}
	return true

}

// func isTCPConnected(host string, port int) bool {
func isTCPConnected(host, port string) bool {

	Addr := host + ":" + port
	conn, err := net.DialTimeout("tcp", Addr, 5*time.Second)
	if err != nil {
		// Connection failed
		return false
	}

	// Close the connection
	defer conn.Close()

	// Connection successful
	return true
}

// func changeEthFormat(dev utilsPkg.DevSettings, mode string) PlcEthInfo {
func changeEthFormat(dev utilsPkg.DevSettings) PlcEthInfo {

	var bitMemories []MemoryEthInfo
	var wordMemories []MemoryEthInfo

	for _, data := range dev.Data {
		bit, _ := data["bitMemories"].([]interface{})
		word, _ := data["wordMemories"].([]interface{})

		for _, item := range bit {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			class, _ := dataMap["class"].(float64)
			instance, _ := dataMap["instance"].(float64)
			attribute, _ := dataMap["attribute"].(float64)
			/*
				if mode == "write" {
					value, _ := dataMap["value"].(float64)
					memory := MemoryEthInfo{
						Name:      name,
						Class:     uint16(class),
						Instance:  uint16(instance),
						Attribute: uint16(attribute),
						Value:     uint32(value),
					}
					bitMemories = append(bitMemories, memory)
				} else {
			*/
			memory := MemoryEthInfo{
				Name:      name,
				Class:     uint16(class),
				Instance:  uint16(instance),
				Attribute: uint16(attribute),
			}
			bitMemories = append(bitMemories, memory)
			//}
		}

		for _, item := range word {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			class, _ := dataMap["class"].(float64)
			instance, _ := dataMap["instance"].(float64)
			attribute, _ := dataMap["attribute"].(float64)
			/*
				if mode == "write" {
					value, _ := dataMap["value"].(float64)
					memory := MemoryEthInfo{
						Name:      name,
						Class:     uint16(class),
						Instance:  uint16(instance),
						Attribute: uint16(attribute),
						Value:     uint32(value),
					}
					wordMemories = append(wordMemories, memory)
				} else {
			*/
			memory := MemoryEthInfo{
				Name:      name,
				Class:     uint16(class),
				Instance:  uint16(instance),
				Attribute: uint16(attribute),
			}
			wordMemories = append(wordMemories, memory)
			//}
		}
	}

	plcInfo := PlcEthInfo{
		BitMemories:  bitMemories,
		WordMemories: wordMemories,
	}

	return plcInfo
}

// UpdateDevConfig will get all informations in file

func UpdateDevConfig(name string, dataMem PlcEthInfo) (utilsPkg.DevSettings, []string, error) {

	var devInf []utilsPkg.Devices
	var itemInf utilsPkg.DevSettings
	var devConfig string = "deviceConfig.json"
	var memChanged []string

	file, err := os.ReadFile(devConfig)
	if err != nil {
		golog.CtsLog.Error("Fail to read JSON File: ", err)
		return itemInf, memChanged, err
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		golog.CtsLog.Error("FAIL to decode JSON file")
		golog.CtsLog.Error("%v", err)
		return itemInf, memChanged, err
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			if item.Name == name {
				newPlcInfo := changeEthFormat(item)
				memChanged = checkMemoryChange(dataMem, newPlcInfo)
				return item, memChanged, nil
			}
		}
	}
	return itemInf, memChanged, err
}

func checkMemoryChange(dataMem, newPlcInfo PlcEthInfo) []string {
	var resMemoryChanged []string

	// Verifica se ocorreu eliminação de alguma memória Bit ou Word
	for _, mem := range dataMem.BitMemories {
		notfound := true
		for _, newMem := range newPlcInfo.BitMemories {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "BitMemory "+mem.Name+" deleted")
		}
	}
	for _, mem := range dataMem.WordMemories {
		notfound := true
		for _, newMem := range newPlcInfo.WordMemories {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "WordMemory "+mem.Name+" deleted")
		}
	}
	// Verifica se ocorreu adição de alguma memória Bit ou Word
	for _, mem := range newPlcInfo.BitMemories {
		notfound := true
		for _, newMem := range dataMem.BitMemories {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "BitMemory "+mem.Name+" added")
		}
	}
	for _, mem := range newPlcInfo.WordMemories {
		notfound := true
		for _, newMem := range dataMem.WordMemories {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "WordMemory "+mem.Name+" added")
		}
	}
	return resMemoryChanged
}
