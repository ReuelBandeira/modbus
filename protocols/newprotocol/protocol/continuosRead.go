package newprotocol

import (
	/// TODO: change newprotocol with your protocolname
	"encoding/json"
	"fmt"
	"net"
	utilsPkg "newprotocol/utils"
	logPkg "newprotocol/utils/gologtofile"
	"os"
	"strings"
	"time"
)

type ProtocolInfo struct {
	Name    string `json:"name"`
	Address uint16 `json:"address"`
}

// / TODO: Replace it struct to your specific protocol struct
type DeviceInfo struct {
	AI []ProtocolInfo `json:"AI"`
	AO []ProtocolInfo `json:"AO"`
	DI []ProtocolInfo `json:"DI"`
	DO []ProtocolInfo `json:"DO"`
}

// func ContinuousRead(dev utilsPkg.DevSettings, setMqtt utilsPkg.SettingsMQTT, client mqtt.Client, loop bool, mode string) {
func ContinuousRead(dev utilsPkg.DevSettings, bInfiniteLoop bool) {
	utilsPkg.AnyDevicelsRunning = true
	defer utilsPkg.WaitGroup.Done()

	var rc byte = 0x80
	var readFailReported bool = false
	var deviceIsConnected bool = false
	var conn net.Conn = nil

	// TODO: According with specific protocol
	var resultsOldAI []uint32 // AI
	var resultsOldAO []uint32 // AO
	var resultsOldDI []uint32 // DI
	var resultsOldDO []uint32 // DO

	// TODO: According with specific protocol
	var curAIValue uint32 // AI
	var curAOValue uint32 // AO
	var curDIValue uint32 // DI
	var curDOValue uint32 // DO

	// TODO: According with specific protocol
	resultsCurAI := make(map[string]interface{}) // AI
	resultsCurAO := make(map[string]interface{}) // AO
	resultsCurDI := make(map[string]interface{}) // DI
	resultsCurDO := make(map[string]interface{}) // DO

	// Save Device on Individual
	infoProtocol := changProtocolFormat(dev)

	// TODO: According with specific protocol
	resultsOldAI = make([]uint32, len(infoProtocol.AI)) // AI
	resultsOldAO = make([]uint32, len(infoProtocol.AO)) // AO
	resultsOldDI = make([]uint32, len(infoProtocol.DI)) // DI
	resultsOldDO = make([]uint32, len(infoProtocol.DO)) // DO

	readHour := time.Now()
	readTime := readHour.Format("2006-01-02 15:04:05")
	anyChangeDetected := false
	if !strings.Contains(dev.Address, ".") || dev.Address == "localhost" || dev.Port != "50000" {
		logPkg.CtsLog.Debug("ContinuousRead:device[%s] Dev[%s:%s] DevAddress Must be nnn.nnn.nnn.nnn:4000 or localhost:40000\n", dev.Name, dev.Address, dev.Port)
		if dev.Address != "localhost" {
			dev.Address = "192.168.18.2" // Forced
		}
		dev.Port = "50000" // Forced
	}
	conn, err := TCPDeviceConnect(dev.Address, dev.Port)
	if err != nil {
		logPkg.CtsLog.Error("ContinuousRead: Unable to Connect with device[%s] Dev[%s:%s]\n",
			dev.Name, dev.Address, dev.Port)
		utilsPkg.AnyDeviceIsConnected = false
		MQTTSendGeneralDeviceStatus(dev, "disconnected")
		msgLog := utilsPkg.MqttLogMsg{
			Address:  dev.Address,
			Port:     dev.Port,
			Name:     dev.Name,
			Protocol: dev.Protocol,
			Id:       dev.Id,
			Time:     readTime,
			Message:  "Disconnected",
			Code:     "200",
			Source:   "ContinuousRead:ContinuousRead",
		}
		_ = MQTTSendLogMessage(msgLog, "error")
	} else {
		deviceIsConnected = true
		logPkg.CtsLog.Warn("ContinuousRead:Connect with device[%s] Dev[%s:%s] conn[%v]\n", dev.Name, dev.Address, dev.Port, conn)
		MQTTSendGeneralDeviceStatus(dev, "connected")
		defer conn.Close()
	}

	/// =======================
	/// GoRotine: Infinite Loop
	/// =======================
	for {
		select {
		//
		// Received: MQTT Client [Stop]
		//
		case <-utilsPkg.StopChan:
			if utilsPkg.AnyDevicelsRunning {
				utilsPkg.AnyDevicelsRunning = false
				logPkg.CtsLog.Warn("utilsPkg.StopChan[%v]true > false Stopping GoRoutine[%s] \n", utilsPkg.AnyDevicelsRunning, dev.Name)
				// Generating MQTT Message and device status stopped
				MQTTSendGeneralDeviceStatus(dev, "stopped")
				msgLog := utilsPkg.MqttLogMsg{
					Address:  dev.Address,
					Port:     dev.Port,
					Name:     dev.Name,
					Protocol: dev.Protocol,
					Id:       dev.Id,
					Time:     readTime,
					Message:  "Device Stopped",
					Code:     "OK",
					Source:   "continuousREad:ContinuousRead",
				}
				_ = MQTTSendLogMessage(msgLog, "info")
			} else {
				logPkg.CtsLog.Warn("utilsPkg.StopChan[%v]false GoRoutine[%s]Already Stopped\n", utilsPkg.AnyDevicelsRunning, dev.Name)
			}
			return

		//
		// Received: MQTT Client [Update]
		//
		case <-utilsPkg.UpdateDevChan:
			devConfig, memChanged, err := UpdateDevConfig(dev.Name, infoProtocol)
			if err != nil {
				logPkg.CtsLog.Error("UpdateDevConfig dev.Name[%s]Could not read deviceConfig.json File \ninfoProtocol[%v] \nerr[%v]", dev.Name, err)
			} else {
				if devConfig.Name == "" {
					if utilsPkg.AnyDevicelsRunning {
						utilsPkg.AnyDevicelsRunning = false
						logPkg.CtsLog.Debug("utilsPkg.UpdateDevChan[%v] dev.Name[", dev.Name, "]Empty\n", utilsPkg.UpdateDevChan)
					} else {
						logPkg.CtsLog.Debug("utilsPkg.UpdateDevChan[%v] dev.Name[", dev.Name, "]\n", utilsPkg.UpdateDevChan)
					}
					return
				}
				// Updating dev.Name
				logPkg.CtsLog.Debug("utilsPkg.UpdateDevChan[%v]  dev.Name[", dev.Name, "]\n", utilsPkg.UpdateDevChan)
				infoProtocol = changProtocolFormat(devConfig) // Reformat infoProtocol
				if len(memChanged) != 0 {
					logPkg.CtsLog.Debug("utilsPkg.UpdateDevChan[%v]  dev.Name[", dev.Name, "]\n infoProtocol[%v]", utilsPkg.UpdateDevChan, infoProtocol)
					// MemChanged
					for i := 0; i < len(memChanged); i++ {
						logPkg.CtsLog.Debug("memChanged[%d]=[%s]\n", i, memChanged[i])
					}
				} else {
					logPkg.CtsLog.Debug("utilsPkg.UpdateDevChan[%v]  dev.Name[", dev.Name, "]\n infoProtocol[%v]", utilsPkg.UpdateDevChan, infoProtocol)
				}
			}

		default:
			/// TODO: Replace it to your specific protocol struct
			/// logPkg.CtsLog.Debug("routineRead: name[%s] Addr[%s:%s] RT[%v]ms\n", dev.Name, dev.Address, dev.Port, dev.ReadingTime)
			/// logPkg.CtsLog.Debug("routineRead: AI  [%v]\n", infoProtocol.AI)
			/// logPkg.CtsLog.Debug("routineRead: AO [%v]\n", infoProtocol.AO)
			/// logPkg.CtsLog.Debug("routineRead: DI [%v]\n", infoProtocol.DI)
			/// logPkg.CtsLog.Debug("routineRead: DO[%v]\n", infoProtocol.DO)
			utilsPkg.AnyDevicelsRunning = true

			/// TODO: Replace it to your specific protocol struct
			if len(infoProtocol.DI) != 0 && CheckConnection(conn, dev.Address, dev.Port) {
				/// =========================================
				/// TODO: Specific protocol: Digital Input(DI)
				/// =========================================
				/// Dev is Connected and has Digigal Input
				for idx, itemMemories := range infoProtocol.DI {
					/// ================================
					/// TODO: Specific Protocol ReadData
					/// ================================
					if conn == nil || !deviceIsConnected {
						continue
					}
					Address := itemMemories.Address
					logPkg.CtsLog.Debug("ReadDigitalInput:(DI) Address[0x%X(%d)]", Address, Address)
					rcvdByteValue, err := ReadDigitalInput(conn, Address)
					curDIValue = uint32(rcvdByteValue)
					if err != nil {
						if conn != nil {
							logPkg.CtsLog.Error("ReadDigitalInput:(DI) err[%s] conn[%v] Disconnected", err, conn)
							MQTTSendGeneralDeviceStatus(dev, "disconnected")
							conn.Close()
							conn = nil
							deviceIsConnected = false
						}
						if !readFailReported {
							readFailReported = true
							rc = 201
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "Fail to Read",
								Code:     "" + fmt.Sprintf("%03d", rc),
								Source:   "continuouRsead:ReadDigitalInput:(DI)",
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("continuouRsead:ReadDigitalInput:(DI) ", dev.Address, " MenName ", " Address = ", itemMemories.Address, itemMemories.Name, " curDIValue = ", curDIValue)
							break
						}
					} else {
						readFailReported = false
						rc = 0
						logPkg.CtsLog.Debug("continuouRsead:ReadDigitalInput:(DI)OK rcvdByteValue[%X] curDIValue[%X]", rcvdByteValue, curDIValue)
						resultsCurDI[itemMemories.Name] = curDIValue
						if resultsOldDI[idx] != curDIValue {
							//curValue is not equal Saved OldDIValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldDI[idx],
							///	curDIValue)
							resultsOldDI[idx] = curDIValue
						}
					}
				}
			}

			/// TODO: Replace it to your specific protocol struct
			if len(infoProtocol.DO) != 0 && CheckConnection(conn, dev.Address, dev.Port) {
				/// =========================================
				/// TODO: Specific protocol: Digital Output(DO)
				/// =========================================
				/// Dev is Connected and has Digigal Output
				for idx, itemMemories := range infoProtocol.DO {
					/// ================================
					/// TODO: Specific Protocol ReadData
					/// ================================
					if conn == nil || !deviceIsConnected {
						continue
					}
					Address := itemMemories.Address
					logPkg.CtsLog.Debug("ReadDigitalOutput:(DO) Address[0x%X(%d)]", Address, Address)
					rcvdByteValue, err := ReadDigitalOutput(conn, Address)
					curDOValue = uint32(rcvdByteValue)
					if err != nil {
						if conn != nil {
							logPkg.CtsLog.Error("ReadDigitalOutput:(DO)  err[%s] conn[%v] Disconnected", err, conn)
							MQTTSendGeneralDeviceStatus(dev, "disconnected")
							conn.Close()
							conn = nil
						}
						if !readFailReported {
							readFailReported = true
							rc = 202
							// FAIL to read NewPLC call MQTTSendLogMessage
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "", // TODO: Convert rc to Text
								Code:     "" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("continuouRsead:ReadDigitalOutput:(DO)", dev.Address, " MenName ", itemMemories.Name, " Address = ", itemMemories.Address),
							}
							_ = MQTTSendLogMessage(msgLog, "error")

							logPkg.CtsLog.Error("continuouRsead:ReadDigitalOutput:(DO)", dev.Address, " MemName ", itemMemories.Name, " Address = ", itemMemories.Address, " curDOValue = ", curDOValue)
							break
						}

					} else {
						rc = 0
						logPkg.CtsLog.Debug("continuouRsead:ReadDigitalOutput:(DO):OK rcvdByteValue[%X] curDOValue[%X]", rcvdByteValue, curDOValue)
						resultsCurDO[itemMemories.Name] = curDOValue
						if resultsOldDO[idx] != curDOValue {
							//curValue is not equal Saved OldDOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldDO[idx],
							///	curDOValue)
							resultsOldDO[idx] = curDOValue
						}
					}
				}
			}

			/// =========================================
			/// TODO: Specific protocol: Analog Input(AI)
			/// =========================================
			/// TODO: Replace it to your specific protocol struct
			if len(infoProtocol.AI) != 0 && CheckConnection(conn, dev.Address, dev.Port) {
				/// Dev is Connected and has Analog Input
				for idx, itemMemories := range infoProtocol.AI {
					/// ================================
					/// TODO: Specific Protocol ReadData
					/// ================================
					if conn == nil || !deviceIsConnected {
						continue
					}
					Address := itemMemories.Address
					logPkg.CtsLog.Debug("continuouRsead:ReadAnalogInput:(AI) Address[0x%X(%d)]", Address, Address)
					rcvdWordValue, err := ReadAnalogInput(conn, Address)
					curAIValue = uint32(rcvdWordValue)
					if err != nil {
						if conn != nil {
							logPkg.CtsLog.Error("continuouRsead:ReadAnalogInput:(AI) err[%s] conn[%v] Disconnected", err, conn)
							MQTTSendGeneralDeviceStatus(dev, "disconnected")
							conn.Close()
							conn = nil
						}
						if !readFailReported {
							readFailReported = true
							rc = 203
							// FAIL to read NewPLC call MQTTSendLogMessage
							logPkg.CtsLog.Error("continuouRsead:ReadAnalogInput:(AI) err[%s] conn[%v]", err, conn)
							// FAIL to read NewPLC call MQTTSendLogMessage
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "", // TODO: Convert rc to Text
								Code:     "" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("continuouRsead:ReadAnalogInput:(AI)", dev.Address, " PLC - Memory ", itemMemories.Name),
							}
							_ = MQTTSendLogMessage(msgLog, "error")

							logPkg.CtsLog.Error("continuouRsead:ReadAnalogInput:(AI) ", dev.Address, " MemName ", itemMemories.Name, " Address = ", itemMemories.Address, " Value = ", curAIValue)
						}
						break
					} else {
						rc = 0
						logPkg.CtsLog.Debug("continuouRsead:ReadAnalogInput:(AI) rcvdWordValue[%X] curAIValue[%X]", rcvdWordValue, curAIValue)
						resultsCurAI[itemMemories.Name] = curAIValue
						if resultsOldAI[idx] != curAIValue {
							//curValue is not equal Saved OldAIValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldAI[idx],
							///	curAIValue)
							resultsOldAI[idx] = curAIValue //fmt.Sprintf("%v", curAIValue)
						}
					}

				}
			}

			/// =========================================
			/// TODO: Specific protocol: Analog Output(AO)
			/// =========================================
			/// TODO: Replace it to your specific protocol struct
			if len(infoProtocol.AO) != 0 && CheckConnection(conn, dev.Address, dev.Port) {
				/// Dev is Connected and has Analog Input
				for idx, itemMemories := range infoProtocol.AO {
					/// ================================
					/// TODO: Specific Protocol ReadData
					/// ================================
					if conn == nil || !deviceIsConnected {
						continue
					}
					Address := itemMemories.Address
					logPkg.CtsLog.Debug("continuousRead:ReadAnalogOutput:(AO) Address[0x%X(%d)]", Address, Address)
					rcvdWordValue, err := ReadAnalogOutput(conn, Address)
					curAOValue = uint32(rcvdWordValue)
					if err != nil {
						if conn != nil {
							logPkg.CtsLog.Error("continuousRead:ReadAnalogOutput:(AO) err[%s] conn[%v] Disconnected", err, conn)
							MQTTSendGeneralDeviceStatus(dev, "disconnected")
							conn.Close()
							conn = nil
						}
						if !readFailReported {
							readFailReported = true
							rc = 204
							// FAIL to read NewPLC call MQTTSendLogMessage
							logPkg.CtsLog.Error("continuousRead:ReadAnalogOutput:(AO) err[%s] conn[%v]", err, conn)
							// FAIL to read NewPLC call MQTTSendLogMessage
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "", // TODO: Convert rc to Text
								Code:     "" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("continuousRead:ReadAnalogOutput:(AO) ", dev.Address, " MemName", itemMemories.Name, " Address = ", itemMemories.Address),
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("continuousRead:ReadAnalogOutput:(AO) ", dev.Address, " MemName ", itemMemories.Name, " Address = ", itemMemories.Address, " curAOValue = ", curAOValue)
						}
						break

					} else {
						rc = 0
						logPkg.CtsLog.Debug("continuousRead:ReadAnalogOutput:(AO) rcvdByteValue[%X] curAOValue[%X]", rcvdWordValue, curAOValue)
						resultsCurAO[itemMemories.Name] = curAOValue
						if resultsOldAO[idx] != curAOValue {
							//curValue is not equal Saved OldAOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldAO[idx],
							///	curAOValue)
							resultsOldAO[idx] = curAOValue
						}
					}

				}
			}

			if anyChangeDetected && !readFailReported {
				// Houve alguma alteração gerar payload JSON da resposta
				msg := utilsPkg.NewDeviceExitPayloadJsonMsg{
					Address:       dev.Address,
					Port:          dev.Port,
					Name:          dev.Name,
					Others:        "",
					Protocol:      dev.Protocol,
					Id:            dev.Id,
					ReadTimeStamp: readTime, ///In Milliseconds
					Topics:        dev.Topics,
					AI:            resultsCurAI,
					AO:            resultsCurAO,
					DI:            resultsCurDI,
					DO:            resultsCurDO,
				}
				// Enviar payload json da mensagem
				err := MQTTSendMessageAutomatic(msg)
				if err != nil {
					logPkg.CtsLog.Error("routineRead:FAIL dev name[%s]  Ip[%s:%s] Changes Occurred - MQTTSendMessageAutomatic\n err[%s]ms\n", dev.Name, dev.Address, dev.Port, err)
				}
				logPkg.CtsLog.Warn("routineRead: dev name[%s]  Ip[%s:%s] Changes Occurred - msg\n%v\n", dev.Name, dev.Address, dev.Port, msg)
			}
			if conn == nil {
				// Device was not connected trying to connect nows
				conn, err = TCPDeviceConnect(dev.Address, dev.Port)
				if err != nil {
					if deviceIsConnected {
						deviceIsConnected = false
						logPkg.CtsLog.Error("ContinuousRead: Unable to Connect with device[%s] Dev[%s:%s]\n",
							dev.Name, dev.Address, dev.Port)
						MQTTSendGeneralDeviceStatus(dev, "disconnected")
						msgLog := utilsPkg.MqttLogMsg{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
							Id:       dev.Id,
							Time:     readTime,
							Message:  "Disconnected",
							Code:     "200",
							Source:   "ContinuousRead:ContinuousRead",
						}
						_ = MQTTSendLogMessage(msgLog, "error")
					}
				} else {
					if !deviceIsConnected {
						deviceIsConnected = true
						MQTTSendGeneralDeviceStatus(dev, "connected")
						msgLog := utilsPkg.MqttLogMsg{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
							Id:       dev.Id,
							Time:     readTime,
							Message:  "Device Connected",
							Code:     "OK",
							Source:   "continuousRead:ContinuousRead",
						}
						_ = MQTTSendLogMessage(msgLog, "info")
					}
				}
			}

			if utilsPkg.AnyDeviceIsConnected {
				// If connected wait deviceConfig.json ReadingTime
				time.Sleep(time.Duration(dev.ReadingTime) * time.Millisecond)
				if utilsPkg.DevicelUpdateRequest {
					utilsPkg.DevicelUpdateRequest = false
					if anyChangeDetected {
						MQTTSendGeneralDeviceStatus(dev, "updated")
						msgLog := utilsPkg.MqttLogMsg{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
							Id:       dev.Id,
							Time:     readTime,
							Message:  "Device Updated",
							Code:     "OK",
							Source:   "continuousREad:ContinuousRead",
						}
						_ = MQTTSendLogMessage(msgLog, "info")
					}
				}
				anyChangeDetected = false
			} else {
				// If Disconnected wait 20 seconds
				time.Sleep(time.Duration(20) * time.Millisecond)
			}
		} // end of default:
		if !bInfiniteLoop {
			break
		}

	} // End of GoRotine: Infinite Loop
	logPkg.CtsLog.Debug("routineRead: dev name[%s]  Ip[%s:%s] bInfiniteLoop[%v] Exiting GoRoutine...\n", dev.Name, dev.Address, dev.Port, bInfiniteLoop)
}

func CheckConnection(conn net.Conn, ipAddress string, port string) bool {

	tcpAddress := ipAddress + ":" + port

	if !isTCPConnected(conn, ipAddress, port) {
		if utilsPkg.AnyDeviceIsConnected {
			utilsPkg.AnyDeviceIsConnected = false
			logPkg.CtsLog.Error("CheckConnection: FAIL - tcpAddress[%s] Disconnected", tcpAddress)
		}
		return false
	}
	// logPkg.CtsLog.Warn("CheckConnection: PASS - tcpAddress[%s] Connected", tcpAddress)
	return true

}

func TCPDeviceConnect(host, port string) (net.Conn, error) {
	tcpAddress := host + ":" + port
	conn, err := net.DialTimeout("tcp", tcpAddress, 5*time.Second)
	if err != nil {
		// Fail to Connect
		if utilsPkg.AnyDeviceIsConnected {
			logPkg.CtsLog.Error("isTCPConnected:FAIL net.DialTimeout tcpAddress[%s] Disconnected", tcpAddress)
			utilsPkg.AnyDeviceIsConnected = false
		}
		return nil, err
	}
	// Connected
	if !utilsPkg.AnyDeviceIsConnected {
		logPkg.CtsLog.Warn("isTCPConnected:PASS tcpAddress[%s] Connected conn[%v]", tcpAddress, conn)
	}
	utilsPkg.AnyDeviceIsConnected = true
	return conn, nil
}

// func isTCPConnected(host string, port int) bool {
func isTCPConnected(conn net.Conn, host, port string) bool {
	tcpAddress := host + ":" + port
	if conn != nil {
		if !utilsPkg.AnyDeviceIsConnected {
			logPkg.CtsLog.Warn("isTCPConnected:WARN tcpAddress[%s] conn[%v]AlConnected", tcpAddress, conn)
			utilsPkg.AnyDeviceIsConnected = true
		}
		return true
	}
	utilsPkg.AnyDeviceIsConnected = false
	return false
}

// func changProtocolFormat(dev utilsPkg.DevSettings, mode string) DeviceInfo {
// TODO: Specific Protocol >>> protocol_config.json <<<
func changProtocolFormat(dev utilsPkg.DevSettings) DeviceInfo {

	/// TODO: Replace it to your specific protocol struct
	var AI []ProtocolInfo // TODO: Specific Protocol
	var AO []ProtocolInfo // TODO: Specific Protocol
	var DI []ProtocolInfo // TODO: Specific Protocol
	var DO []ProtocolInfo // TODO: Specific Protocol

	for _, data := range dev.Data {
		/// logPkg.CtsLog.Debug("changProtocolFormat:\ndata[%v]\n", data)

		/// TODO: Replace it to your specific protocol struct
		analogInput, _ := data["AI"].([]interface{})   // TODO: Specific Protocol
		analogOutput, _ := data["AO"].([]interface{})  // TODO: Specific Protocol
		digitalInput, _ := data["DI"].([]interface{})  // TODO: Specific Protocol
		digitalOutput, _ := data["DO"].([]interface{}) // TODO: Specific Protocol

		for _, item := range analogInput {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)         // TODO: Specific Protocol
			addrress, _ := dataMap["address"].(float64) // TODO: Specific Protocol
			memory := ProtocolInfo{                     // TODO: Specific Protocol
				Name:    name,
				Address: uint16(addrress),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:utem[%s] name[%s] address[%s]\n", item,name,address)
			AI = append(AI, memory) // TODO: Specific Protocol
		}

		for _, item := range analogOutput {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)         // TODO: Specific Protocol
			addrress, _ := dataMap["address"].(float64) // TODO: Specific Protocol
			memory := ProtocolInfo{                     // TODO: Specific Protocol
				Name:    name,             // TODO: Specific Protocol
				Address: uint16(addrress), // TODO: Specific Protocol
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:utem[%s] name[%s] address[%s]\n", item,name,address)
			AO = append(AO, memory)
		}

		for _, item := range digitalInput {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)         // TODO: Specific Protocol
			addrress, _ := dataMap["address"].(float64) // TODO: Specific Protocol
			memory := ProtocolInfo{                     // TODO: Specific Protocol
				Name:    name,             // TODO: Specific Protocol
				Address: uint16(addrress), // TODO: Specific Protocol
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:utem[%s] name[%s] address[%s]\n", item,name,address)
			DI = append(DI, memory) // TODO: Specific Protocol
		}

		for _, item := range digitalOutput {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)        // TODO: Specific Protocol
			address, _ := dataMap["address"].(float64) // TODO: Specific Protocol
			memory := ProtocolInfo{                    // TODO: Specific Protocol
				Name:    name,            // TODO: Specific Protocol
				Address: uint16(address), // TODO: Specific Protocol
			}
			DO = append(DO, memory) // TODO: Specific Protocol
			/// logPkg.CtsLog.Debug("changProtocolFormat:utem[%s] name[%s] address[%s]\n", item,name,address)
		}
	}

	/// TODO: Replace it to your specific protocol struct
	deviceNewInfo := DeviceInfo{ // TODO: Specific Protocol
		AI: AI, // TODO: Specific Protocol
		AO: AO, // TODO: Specific Protocol
		DI: DI, // TODO: Specific Protocol
		DO: DO, // TODO: Specific Protocol
	}
	return deviceNewInfo // TODO: Specific Protocol
}

// UpdateDevConfig will get all informations in file
func UpdateDevConfig(name string, dataMem DeviceInfo) (utilsPkg.DevSettings, []string, error) {
	var devInf []utilsPkg.Devices
	var itemInf utilsPkg.DevSettings
	var devConfig string = "deviceConfig.json"
	var memChanged []string

	file, err := os.ReadFile(devConfig)
	if err != nil {
		logPkg.CtsLog.Error("Fail to read JSON File: ", err)
		return itemInf, memChanged, err
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		logPkg.CtsLog.Error("FAIL to decode JSON file")
		logPkg.CtsLog.Error("%v", err)
		return itemInf, memChanged, err
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			if item.Name == name {
				newDeviceInfo := changProtocolFormat(item)
				memChanged = checkMemoryChange(dataMem, newDeviceInfo)
				return item, memChanged, nil
			}
		}
	}
	return itemInf, memChanged, err
}

// TODO: Specific Protocol
func checkMemoryChange(dataMem, newDeviceInfo DeviceInfo) []string {
	var resMemoryChanged []string

	// Verifica se ocorreu eliminação de alguma memória Bit ou Word
	/// TODO: Replace it to your specific protocol struct
	for _, mem := range dataMem.DI { // TODO: Specific Protocol
		notfound := true
		for _, newMem := range newDeviceInfo.DI { // TODO: Specific Protocol
			if mem.Name == newMem.Name { // TODO: Specific Protocol
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "DI "+mem.Name+" deleted")
		}
	}
	for _, mem := range dataMem.DO { // TODO: Specific Protocol
		notfound := true
		for _, newMem := range newDeviceInfo.DO { // TODO: Specific Protocol
			if mem.Name == newMem.Name { // TODO: Specific Protocol
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "DO "+mem.Name+" deleted")
		}
	}
	// Verifica se ocorreu adição de alguma memória Bit ou Word
	/// TODO: Replace it to your specific protocol struct
	for _, mem := range newDeviceInfo.DI { // TODO: Specific Protocol
		notfound := true
		for _, newMem := range dataMem.DI { // TODO: Specific Protocol
			if mem.Name == newMem.Name { // TODO: Specific Protocol
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "BitMemory "+mem.Name+" added")
		}
	}
	for _, mem := range newDeviceInfo.DO { // TODO: Specific Protocol{
		notfound := true
		for _, newMem := range dataMem.DO { // TODO: Specific Protocol
			if mem.Name == newMem.Name { // TODO: Specific Protocol
				notfound = false
			}
		}
		if notfound {
			resMemoryChanged = append(resMemoryChanged, "DO "+mem.Name+" added")
		}
	}
	return resMemoryChanged
}
