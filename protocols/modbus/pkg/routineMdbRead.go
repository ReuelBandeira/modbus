package pkg

import (
	//"sync"
	"fmt"
	"sync"
	"time"

	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"

	modbusLib "github.com/goburrow/modbus"
	"github.com/holdenoffmenn/functions"
)

var bitMemoriesMutex sync.Mutex
var wordMemoriesMutex sync.Mutex

// ReadInfoMdbs will read informations from device
func ReadInfoMdbs(device utilsPkg.Device, identifier rune) {
	//Get device data from device
	//var deviceData utilsPkg.Data
	deviceData := device.Data[0]

	wg := utilsPkg.WaitGroups[identifier]
	defer wg.Done()

	//Start Connection
	conn := modbusLib.NewTCPClientHandler(device.Address + ":" + device.Port)
	conn.Timeout = 5 * time.Second
	conn.SlaveId = deviceData.SlaveID

	err := conn.Connect()
	defer conn.Close()

	//Armazena resultados
	resultsGetBitMemories := make(map[string]bool)
	resultsGetWordMemories := make(map[string]uint32)
	lastBitMemories := make(map[string]bool)
	lastWordMemories := make(map[string]uint32)

	if err != nil {
		logger.Log.Errorf("routineMdbRead:ReadInfoMdbs: Unable to Connect with Modbus device [%s] error [%s]",
			device.Name, err)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Unable to Connect with Modbus device. Error [%v]", err),
				Code:     "",
				Source:   "routineMdbRead:ReadInfoMdbs",
			},
		}
		Sender(utilsPkg.MsgLog, utilsPkg.Log)
		return
	} else {
		logger.Log.Infof("routineMdbRead:ReadInfoMdbs: Starting device reading [%v]",
			device.ID)

		msgLog := utilsPkg.MessageLog{
			MessageType: "info",
			Data: utilsPkg.DataLog{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  "Starting device communication with modbus.",
				Code:     "",
				Source:   "routineMdbRead:ReadInfoMdbs",
			},
		}
		Sender(msgLog, utilsPkg.Log)

		clientModbus := modbusLib.NewClient(conn)
		readTime := time.Now().Format("2006-01-02 15:04:05") //readHour.Format("2006-01-02 15:04:05")

		//loop infinito de verificação, repetindo conforme o tempo informado no arquivo config.json
		for {
			select {
			case <-utilsPkg.StopChannels[identifier]:
				utilsPkg.DeviceWorkStatus[identifier] = "stopped"
				logger.Log.Infof("routineMdbRead:ReadInfoMdbs: Stopped Modbus Device Name[%v] ID[%v].",
					device.Name, identifier)
				utilsPkg.MsgLog = utilsPkg.MessageLog{
					MessageType: "info",
					Data: utilsPkg.DataLog{
						Address:  device.Address,
						Port:     device.Port,
						Name:     device.Name,
						Protocol: device.Protocol,
						Time:     time.Now().Format("2006-01-02 15:04:05"),
						Message:  "Stop Modbus Device.",
						Code:     "",
						Source:   "routineMdbRead:ReadInfoMdbs",
					},
				}
				Sender(utilsPkg.MsgLog, utilsPkg.Log)

				return

			default:
				utilsPkg.DeviceWorkStatus[identifier] = "running"
				if len(deviceData.BitMemories) != 0 {
					// Itera items configurados no arquivo config.json
					for _, itemMemories := range deviceData.BitMemories {
						resultBit, err := readBitMemory(uint8(deviceData.SlaveID), clientModbus, itemMemories.Address,
							device)
						if err != nil {
							//Forçar o fechamento da conexão
							conn.Close()
							//Adiciona FAIL como valor na chave que houve falha de leitura
							bitMemoriesMutex.Lock()
							resultsGetBitMemories[itemMemories.Name] = resultBit
							bitMemoriesMutex.Unlock()
							//Chama o loop de verificação de conexão
							_ = ConnModbus(device) // connCheck(conn)							
						} else {
							// Salva o resultado encontrado
							bitMemoriesMutex.Lock()
							resultsGetBitMemories[itemMemories.Name] = resultBit
							bitMemoriesMutex.Unlock()
						}

					}

				}
				if len(deviceData.WordMemories) != 0 {
					// Varre items configurados no arquivo config.json
					for _, itemMemories := range deviceData.WordMemories {
						resultWord, err := readWordMemory(
							deviceData.SlaveID,
							clientModbus,
							itemMemories.Address,
							itemMemories.Format,
							device,
						)

						if err != nil {
							//Forçar o fechamento e reabertura da conexão
							conn.Close()
							wordMemoriesMutex.Lock()
							resultsGetWordMemories[itemMemories.Name] = resultWord
							wordMemoriesMutex.Unlock()
							_ = ConnModbus(device) //(conn)							
						} else {
							wordMemoriesMutex.Lock()
							resultsGetWordMemories[itemMemories.Name] = resultWord
							wordMemoriesMutex.Unlock()
						}
					}

				}

				changesBit := functions.CompareMapsStrBool(lastBitMemories, resultsGetBitMemories)
				changesWord := functions.CompareMapsStrUint32(lastWordMemories, resultsGetWordMemories)

				if !changesBit && !changesWord {
					// Não houve alterações
					logger.Log.Debugf("routineMdbRead:ReadInfoMdbs: Device Name[%s]  IP[%s:%s] No  Changes.",
						device.Name, device.Address, device.Port)
				} else {
					logger.Log.Debugf("routineMdbRead:ReadInfoMdbs: Device Name[%s]  IP[%s:%s] Changes Found.",
						device.Name, device.Address, device.Port)

					// Holve alguma alteração gerar payload JSON da resposta
					dataToSend := utilsPkg.MqttMsgStruct{
						Address: utilsPkg.Address{
							Address: device.Address,
							Port:    device.Port,
							Name:    device.Name,
							ID:      device.ID,
							Others:  "",
						},
						Topics:        device.Topics,
						ReadTimeStamp: readTime, //.Format(time.RFC3339),
						Protocol:      device.Protocol,
						Data: utilsPkg.DataExit{
							BitMemories:  resultsGetBitMemories,
							WordMemories: resultsGetWordMemories,
						},
					}

					// Enviar payload json da mensagem
					Sender(dataToSend, device.Topics)
				}

				// Atualiza ultimos valores com os valores correntes
				lastBitMemories = functions.CopyMapStringBool(resultsGetBitMemories)
				lastWordMemories = functions.CopyMapStringUint32(resultsGetWordMemories)

				// Dorme conforme tempo configurado no arquivo json.config
				time.Sleep(time.Duration(device.ReadingTime) * time.Millisecond)

			}

		}

	}
}

// readBitMemory will make a read to device modbus by Bit
func readBitMemory(slaveID uint8, client modbusLib.Client, addressMemory uint16, device utilsPkg.Device) (bool, error) {
	results, err := client.ReadCoils(addressMemory, 1)
	if err != nil {
		logger.Log.Errorf("routineMdbRead:readBitMemory: Unable to read Device [%v] - Memory Bit [%v] - Error [%v]",
			device.Name, addressMemory, err)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Unable to read memory [%v] - Error [%v]", addressMemory, err),
				Code:     "",
				Source:   "routineMdbRead:readBitMemory",
			},
		}
		Sender(utilsPkg.MsgLog, utilsPkg.Log)

		return false, err
	}

	//Convert a resposta de array de byte > int > string
	if results[0] == 1 {
		return true, nil
	}
	return false, nil
}

// readWordMemory will make a read to device modbus by Word
func readWordMemory(slaveID uint8, client modbusLib.Client, addressMemory uint16, typeMem uint8, device utilsPkg.Device) (uint32, error) {
	var resp uint32
	results, err := client.ReadHoldingRegisters(addressMemory, 1)
	if err != nil {
		logger.Log.Errorf("routineMdbRead:readWordMemory: Unable to read Device [%v] - Memory Word [%v] - Error [%v]",
			device.Name, addressMemory, err)

		utilsPkg.MsgLog = utilsPkg.MessageLog{
			MessageType: "error",
			Data: utilsPkg.DataLog{
				Address:  device.Address,
				Port:     device.Port,
				Name:     device.Name,
				Protocol: device.Protocol,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Message:  fmt.Sprintf("Unable to read memory [%v] - Error [%v]", addressMemory, err),
				Code:     "",
				Source:   "routineMdbRead:readWordMemory",
			},
		}
		Sender(utilsPkg.MsgLog, utilsPkg.Log)

		return resp, err
	}

	if len(results) == 4 {
		value32 := uint32(results[0])<<24 | uint32(results[1])<<16 | uint32(results[2])<<8 | uint32(results[3])
		//CtsLog.Info("Read 32-bit value: \n", value32)
		//resp = strconv.Itoa(int(value32))
		resp = value32
	} else if len(results) == 2 {
		value16 := uint32(results[0])<<8 | uint32(results[1])
		//result1 := uint32(results[0])
		//CtsLog.Info("Read 16-bit value: \n", value16)
		//resp = strconv.Itoa(int(value16))
		resp = value16
	}

	return resp, nil
}
