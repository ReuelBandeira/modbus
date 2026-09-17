package canopen

import (
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type ProtocolInfo struct {
	Name    string `json:"name"`
	Address uint16 `json:"address"`
}

type CanOpenEMCYInfo struct {
	Name       string `json:"name"`
	CANOPEN_ID string `json:"canid"`
	NODE_ID    uint8  `json:"nodeid"`
}

type CanOpenSDOInfo struct {
	Name       string `json:"name"`
	CANOPEN_ID string `json:"canid"`
	NODE_ID    uint8  `json:"nodeid"`
	Cmd        uint8  `json:"cmd"`
	ObjIndex   uint32 `json:"objIndex"`
	SubIndex   uint8  `json:"subIndex"`
}

type CanOpenPDOMapping struct {
	ObjIndex uint32 `json:"objIndex"`
	SubIndex uint8  `json:"subIndex"`
}

type CanOpenPDOInfo struct {
	Name       string              `json:"name"`
	CANOPEN_ID string              `json:"canid"`
	NODE_ID    uint8               `json:"nodeid"`
	MAPPING    []CanOpenPDOMapping `json:"mapping"`
}
type CanOpenNMTInfo struct {
	Name       string `json:"name"`
	CANOPEN_ID string `json:"canid"`
	NODE_ID    uint8  `json:"nodeid"`
	Command    uint8  `json:"command"`
}
type CanOpenSYNCInfo struct {
	Name       string `json:"name"`
	CANOPEN_ID string `json:"canid"`
	NODE_ID    uint8  `json:"nodeid"`
}

type CanOpenTIMEInfo struct {
	Name       string `json:"name"`
	CANOPEN_ID string `json:"canid"`
	NODE_ID    uint8  `json:"nodeid"`
}

type DeviceInfo struct {
	SDO  []CanOpenSDOInfo  `json:"SDO"`  // CANopen Protocol Only
	PDO  []CanOpenPDOInfo  `json:"PDO"`  // CANopen Protocol Only
	NMT  []CanOpenNMTInfo  `json:"NMT"`  // CANopen Protocol Only
	SYNC []CanOpenSYNCInfo `json:"SYNC"` // CANopen Protocol Only
	TIME []CanOpenTIMEInfo `json:"TIME"` // CANopen Protocol Only
	EMCY []CanOpenEMCYInfo `json:"EMCY"` // CANopen Protocol Only
}

var iCont = 0

// func ContinuousRead(dev utilsPkg.DevSettings, setMqtt utilsPkg.SettingsMQTT, client mqtt.Client, loop bool, mode string) {
func ContinuousRead(dev utilsPkg.DevSettings, bInfiniteLoop bool) {
	utilsPkg.AnyDevicelsRunning = true
	defer utilsPkg.WaitGroup.Done()

	var netConn net.Conn = nil
	var clientIsConnected bool = false
	var clientReadErrorReported bool = false

	var rc byte

	var resultsOldSDO []uint32  // CANopen Protocol Only
	var resultsOldPDO []uint32  // CANopen Protocol Only
	var resultsOldNMT []uint32  // CANopen Protocol Only
	var resultsOldSYNC []uint32 // CANopen Protocol Only
	var resultsOldTIME []uint32 // CANopen Protocol Only
	var resultsOldEMCY []uint32 // CANopen Protocol Only

	var curSDOValue uint32  // CANopen Protocol Only
	var curPDOValue uint32  // CANopen Protocol Only
	var curNMTValue uint32  // CANopen Protocol Only
	var curSYNCValue uint32 // CANopen Protocol Only
	var curTIMEValue uint32 // CANopen Protocol Only
	var curEMCYValue uint32 // CANopen Protocol Only

	readHour := time.Now()
	readTime := readHour.Format("2006-01-02 15:04:05")
	anyChangeDetected := false

	resultsCurSDO := make(map[string]interface{})  // CANopen Protocol Only
	resultsCurPDO := make(map[string]interface{})  // CANopen Protocol Only
	resultsCurNMT := make(map[string]interface{})  // CANopen Protocol Only
	resultsCurSYNC := make(map[string]interface{}) // CANopen Protocol Only
	resultsCurTIME := make(map[string]interface{}) // CANopen Protocol Only
	resultsCurEMCY := make(map[string]interface{}) // CANopen Protocol Only

	if !CanBusServicesRunning {
		CanBusServicesRunning = true
		//go RunCanBusSendService(dev)
		//go RunCanBusRecvService(dev)
	}
	// Save Device on Individual
	infoProtocol := changProtocolFormat(dev)

	resultsOldSDO = make([]uint32, len(infoProtocol.SDO))   // CANopen Protocol Only
	resultsOldPDO = make([]uint32, len(infoProtocol.PDO))   // CANopen Protocol Only
	resultsOldNMT = make([]uint32, len(infoProtocol.NMT))   // CANopen Protocol Only
	resultsOldSYNC = make([]uint32, len(infoProtocol.SYNC)) // CANopen Protocol Only
	resultsOldTIME = make([]uint32, len(infoProtocol.TIME)) // CANopen Protocol Only
	resultsOldEMCY = make([]uint32, len(infoProtocol.EMCY)) // CANopen Protocol Only

	// Trying to connect to CANopen TCP Server
	netConn, err := CANopenTcpClient_ConnectToServer(netConn, dev.Address, dev.Port)
	if err == nil {
		// Successfull Client
		clientIsConnected = true
		logPkg.CtsLog.Debug("ContinuousRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]Connected\n",
			dev.Name, dev.Protocol, dev.Address, dev.Port)
		// Sending Info Log Device is Connected
		MQTTSendGeneralDeviceStatus(dev, "connected")
		msgLog := utilsPkg.MqttLogMsg{
			Address:  dev.Address,
			Port:     dev.Port,
			Name:     dev.Name,
			Protocol: dev.Protocol,
			Id:       dev.Id,
			Time:     readTime,
			Message:  "connected",
			Code:     "OK",
			Source:   "ContinuousRead:",
		}
		_ = MQTTSendLogMessage(msgLog, "info")
	} else {
		//Fail to Connect with Device
		clientReadErrorReported = false
		clientIsConnected = false
		if netConn != nil {
			netConn.Close()
			netConn = nil
		}
		logPkg.CtsLog.Error("ContinuousRead: DeviceName[%s]\n\t Protocol[%s]\n\t  Ip:Port[%s:%s]Disconnected\n",
			dev.Name, dev.Protocol, dev.Address, dev.Port)
		// Sending Status that Device is connected
		MQTTSendGeneralDeviceStatus(dev, "disconnected")
		// Sending Error Log Device Disconnected
		msgLog := utilsPkg.MqttLogMsg{
			Address:  dev.Address,
			Port:     dev.Port,
			Name:     dev.Name,
			Protocol: dev.Protocol,
			Id:       dev.Id,
			Time:     readTime,
			Message:  "disconnected",
			Code:     "CANOPEN_400",
			Source:   "continuousRead:",
		}
		_ = MQTTSendLogMessage(msgLog, "error")
	}
	logPkg.CtsLog.Debug("ContinuousRead:device[%s] Dev[%s:%s] bInfiniteLoop[%v]Started...\n\t          dev[%v]\n\t infoProtocol[%v]\n\t       netConn[%v]\n", dev.Name, dev.Address, dev.Port, bInfiniteLoop, dev, infoProtocol, netConn)

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
			/// logPkg.CtsLog.Debug("routineRead: name[%s] Addr[%s:%s] RT[%v]ms\n", dev.Name, dev.Address, dev.Port, dev.ReadingTime)
			clientReadErrorReported = false
			if netConn == nil {
				netConn, err = CANopenTcpClient_ConnectToServer(netConn, dev.Address, dev.Port)
				if err == nil {
					// Successfull Client
					clientIsConnected = true
					logPkg.CtsLog.Debug("ContinuousRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]Connected\n",
						dev.Name, dev.Protocol, dev.Address, dev.Port)
					// Sending Info Log Device is Connected
					MQTTSendGeneralDeviceStatus(dev, "connected")
					msgLog := utilsPkg.MqttLogMsg{
						Address:  dev.Address,
						Port:     dev.Port,
						Name:     dev.Name,
						Protocol: dev.Protocol,
						Id:       dev.Id,
						Time:     readTime,
						Message:  "connected",
						Code:     "OK",
						Source:   "ContinuousRead:",
					}
					_ = MQTTSendLogMessage(msgLog, "info")
				} else {
					if clientIsConnected {
						clientIsConnected = false
						logPkg.CtsLog.Error("ContinuousRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]Disconnected\n",
							dev.Name, dev.Protocol, dev.Address, dev.Port)
						// Sending Info Log Device is Connected
						MQTTSendGeneralDeviceStatus(dev, "disconnected")
						msgLog := utilsPkg.MqttLogMsg{
							Address:  dev.Address,
							Port:     dev.Port,
							Name:     dev.Name,
							Protocol: dev.Protocol,
							Id:       dev.Id,
							Time:     readTime,
							Message:  "disconnected",
							Code:     "CANOPEN_200",
							Source:   "ContinuousRead:Fail to Connect",
						}
						_ = MQTTSendLogMessage(msgLog, "error")
					}
				}
			}
			/// ===========================================
			/// CANopen: Service Data Object (SDO) protocol
			/// ===========================================
			if len(infoProtocol.SDO) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Analog Input
				if netConn == nil || !clientIsConnected || clientReadErrorReported {
					continue
				}
				for idx, itemMemories := range infoProtocol.SDO {
					/// TODO  read curSDOValue from CANopenTcpClient
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					err = CANopenTcpClient_Proc_SDO_Download_Msg(netConn, 0x08, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 201
					}
					if rc != 0 {
						// FAIL to read CANopen call MQTTSendLogMessage
						readAnaCurInput := false
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						if utilsPkg.AnyDeviceIsConnected {
							MQTTSendGeneralDeviceStatus(dev, "disconnected")
							utilsPkg.AnyDeviceIsConnected = false
							if !clientReadErrorReported {
								clientReadErrorReported = true
								msgLog := utilsPkg.MqttLogMsg{
									Address:  dev.Address,
									Port:     dev.Port,
									Name:     dev.Name,
									Protocol: dev.Protocol,
									Id:       dev.Id,
									Time:     readTime,
									Message:  "",
									Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
									Source:   fmt.Sprint("ContinuousRead: erroor SDO", dev.Address, " CANopen ", itemMemories.CANOPEN_ID, " Value = ", readAnaCurInput),
								}
								_ = MQTTSendLogMessage(msgLog, "error")
								logPkg.CtsLog.Error("Error SDO", dev.Address, " CANopen ", itemMemories.CANOPEN_ID, " Value = ", curSDOValue)
							}
						}
					} else {
						clientReadErrorReported = false
						curSDOValue++ /// TODO  Remove it
						resultsCurSDO["name"] = itemMemories.Name
						resultsCurSDO["canid"] = itemMemories.CANOPEN_ID
						resultsCurSDO["nodeid"] = itemMemories.NODE_ID
						resultsCurSDO["cmd"] = itemMemories.Cmd
						resultsCurSDO["objIndex"] = itemMemories.ObjIndex
						resultsCurSDO["subIndex"] = itemMemories.SubIndex
						resultsCurSDO[itemMemories.CANOPEN_ID] = curSDOValue
						if resultsOldSDO[idx] != curSDOValue {
							//curValue is not equal Saved OldSDOValue
							anyChangeDetected = true
							resultsOldSDO[idx] = curSDOValue //fmt.Sprintf("%v", curSDOValue)
						}
					}

				}
			}

			/// ===========================================
			/// CANopen: Process Data Object (PDO) protocol
			/// ===========================================
			if len(infoProtocol.PDO) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Analog Input
				for idx, itemMemories := range infoProtocol.PDO {
					/// TODO  read curSDOValue from CANopenTcpClient
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					_, err = CANopenTcpClient_Proc_PDO_Msg(netConn, 0x01, 0x08, 0x1800, 0x01, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 202
					}
					if rc != 0 {
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						// FSDOL to read NewPLC call MQTTSendLogMessage
						if !clientReadErrorReported {
							clientReadErrorReported = true
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "",
								Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("ContinuousRead: error PDO", dev.Address, " CANopen ", itemMemories.Name, " CANOPEN_ID = ", itemMemories.CANOPEN_ID),
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("Error PDO", dev.Address, " CANopen ", itemMemories.Name, " Value = ", curPDOValue)
						}
					} else {
						clientReadErrorReported = false
						curPDOValue++ /// TODO  Remove it
						resultsCurPDO["name"] = itemMemories.Name
						resultsCurPDO["canid"] = itemMemories.CANOPEN_ID
						resultsCurPDO["nodeid"] = itemMemories.NODE_ID
						resultsCurPDO["mapping"] = itemMemories.MAPPING
						///resultsCurPDO["objIndex"] = itemMemories.MAPPING[idx].ObjIndex
						///resultsCurPDO["subIndex"] = itemMemories.MAPPING[idx].SubIndex
						resultsCurPDO[itemMemories.Name] = curPDOValue
						if resultsOldPDO[idx] != curPDOValue {
							//curValue is not equal Saved OldAOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldPDO[idx],
							///	curPDOValue)
							resultsOldPDO[idx] = curPDOValue
						}
					}
				}
			}

			/// ===========================================
			/// CANopen: Network management (NMT) protocols
			/// ===========================================
			if len(infoProtocol.NMT) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Digigal Input
				for idx, itemMemories := range infoProtocol.NMT {
					/// Dev is Connected and has Digital Input
					/// logPkg.CtsLog.Debug("routineRead: idx[%d] itemMemories.Address[%d] itemMemories.Name[%s]\n", idx, itemMemories.Address, itemMemories.Name)
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					_, err = CANopenTcpClient_Proc_PDO_Msg(netConn, 0x01, 0x08, 0x1800, 0x01, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 203
					}
					if rc != 0 {
						// FSDOL to read NewPLC call MQTTSendLogMessage
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						if !clientReadErrorReported {
							clientReadErrorReported = true
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "",
								Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("ContinuousRead: error NMT", dev.Address, " CANopen ", itemMemories.Name),
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("Error NMT", dev.Address, " CANopen ", itemMemories.Name)
						}
					} else {
						clientReadErrorReported = false
						curNMTValue++
						resultsCurNMT["name"] = itemMemories.Name
						resultsCurNMT["canid"] = itemMemories.CANOPEN_ID
						resultsCurNMT["nodeid"] = itemMemories.NODE_ID
						resultsCurNMT["command"] = itemMemories.Command
						resultsCurNMT[itemMemories.Name] = curNMTValue
						if resultsOldNMT[idx] != curNMTValue {
							//curValue is not equal Saved OldDIValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldNMT[idx],
							///	curNMTValue)
							resultsOldNMT[idx] = curNMTValue
						}
					}
				}
			}

			/// ===============================================
			/// CANopen: Synchronization Object (SYNC) protocol
			/// ===============================================
			if len(infoProtocol.SYNC) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Digigal Output
				for idx, itemMemories := range infoProtocol.SYNC {
					/// Dev is Connected and has Digital Output
					/// logPkg.CtsLog.Debug("routineRead: idx[%d] itemMemories.Address[%d] itemMemories.Name[%s]\n", idx, itemMemories.Address, itemMemories.Name)
					/// TODO  Access NewPLC too read curSYNCValue
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					_, err = CANopenTcpClient_Proc_PDO_Msg(netConn, 0x01, 0x08, 0x1800, 0x01, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 204
					}
					if rc != 0 {
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						// FSDOL to read NewPLC call MQTTSendLogMessage
						if !clientReadErrorReported {
							clientReadErrorReported = true
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "",
								Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("ContinuousRead: error SYNC", dev.Address, " CANopen ", itemMemories.Name, " Name = ", itemMemories.Name),
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("Error SYNC", dev.Address, " CANopen ", itemMemories.Name, " Value = ", curSYNCValue)
						}
					} else {
						clientReadErrorReported = false
						curSYNCValue++
						resultsCurSYNC["name"] = itemMemories.Name
						resultsCurSYNC["canid"] = itemMemories.CANOPEN_ID
						resultsCurSYNC["nodeid"] = itemMemories.NODE_ID
						resultsCurSYNC[itemMemories.Name] = curSYNCValue
						if resultsOldSYNC[idx] != curSYNCValue {
							//curValue is not equal Saved OldDOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldSYNC[idx],
							///	curSYNCValue)
							resultsOldSYNC[idx] = curSYNCValue
						}
					}
				}
			}

			/// ==========================================
			/// CANopen: Time Stamp Object (TIME) protocol
			/// ==========================================
			if len(infoProtocol.TIME) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Digigal Output
				for idx, itemMemories := range infoProtocol.TIME {
					/// Dev is Connected and has Digital Output
					/// logPkg.CtsLog.Debug("routineRead: idx[%d] itemMemories.Address[%d] itemMemories.Name[%s]\n", idx, itemMemories.Address, itemMemories.Name)
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					_, err = CANopenTcpClient_Proc_PDO_Msg(netConn, 0x01, 0x08, 0x1800, 0x01, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 205
					}
					if rc != 0 {
						// FSDOL to read NewPLC call MQTTSendLogMessage
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						if !clientReadErrorReported {
							clientReadErrorReported = true
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "",
								Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("ContinuousRead: TIME", dev.Address, " CANopen ", itemMemories.Name, " CANOPEN_ID = ", itemMemories.CANOPEN_ID),
							}
							_ = MQTTSendLogMessage(msgLog, "error")
							logPkg.CtsLog.Error("Error TIME", dev.Address, " CANopen ", itemMemories.Name, " Value = ", curTIMEValue)
						}
					} else {
						clientReadErrorReported = false
						curTIMEValue++
						resultsCurTIME["name"] = itemMemories.Name
						resultsCurTIME["canid"] = itemMemories.CANOPEN_ID
						resultsCurTIME["nodeid"] = itemMemories.NODE_ID
						resultsCurTIME[itemMemories.Name] = curTIMEValue
						if resultsOldTIME[idx] != curTIMEValue {
							//curValue is not equal Saved OldDOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldTIME[idx],
							///	curTIMEValue)
							resultsOldTIME[idx] = curTIMEValue
						}
					}
				}
			}

			/// =========================================
			/// CANopen: Emergency Object (EMCY) protocol
			/// =========================================
			if len(infoProtocol.EMCY) != 0 && CheckConnection(netConn, dev.Address, dev.Port) {
				/// Dev is Connected and has Digigal Output
				for idx, itemMemories := range infoProtocol.EMCY {
					/// Dev is Connected and has Digital Output
					/// logPkg.CtsLog.Debug("routineRead: idx[%d] itemMemories.Address[%d] itemMemories.Name[%s]\n", idx, itemMemories.Address, itemMemories.Name)
					/// TODO  Access NewPLC too read curEMCYValue
					if netConn == nil || !clientIsConnected || clientReadErrorReported {
						continue
					}
					_, err = CANopenTcpClient_Proc_PDO_Msg(netConn, 0x01, 0x08, 0x1800, 0x01, nil)
					if err == nil {
						rc = 0
					} else {
						clientIsConnected = false
						rc = 206
					}
					if rc != 0 {
						// FSDOL to read NewPLC call MQTTSendLogMessage
						if netConn != nil {
							netConn.Close()
							netConn = nil
						}
						if !clientReadErrorReported {
							clientReadErrorReported = true
							msgLog := utilsPkg.MqttLogMsg{
								Address:  dev.Address,
								Port:     dev.Port,
								Name:     dev.Name,
								Protocol: dev.Protocol,
								Id:       dev.Id,
								Time:     readTime,
								Message:  "",
								Code:     "CANOPEN_" + fmt.Sprintf("%03d", rc),
								Source:   fmt.Sprint("ContinuousRead: EMCY", dev.Address, " CANopen ", itemMemories.Name, " CANOPEN_ID = ", itemMemories.CANOPEN_ID),
							}
							_ = MQTTSendLogMessage(msgLog, "error")

							logPkg.CtsLog.Error("Error EMCY", dev.Address, " CANopen ", itemMemories.Name, " Value = ", curEMCYValue)
						}
					} else {
						clientReadErrorReported = false
						curEMCYValue++
						resultsCurEMCY["name"] = itemMemories.Name
						resultsCurEMCY["canid"] = itemMemories.CANOPEN_ID
						resultsCurEMCY["nodeid"] = itemMemories.NODE_ID
						resultsCurEMCY[itemMemories.Name] = curEMCYValue
						if resultsOldEMCY[idx] != curEMCYValue {
							//curValue is not equal Saved OldDOValue
							anyChangeDetected = true
							///logPkg.CtsLog.Debug("Dev[%s] Changed ItemName[%s] OldValue[%d] NewValue[%d]\n",
							///	dev.Name,
							///	itemMemories.Name,
							///	resultsOldEMCY[idx],
							///	curEMCYValue)
							resultsOldEMCY[idx] = curEMCYValue
						}
					}
				}
			}

			if anyChangeDetected && clientIsConnected && netConn != nil {
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
					SDO:           resultsCurSDO,
					PDO:           resultsCurPDO,
					NMT:           resultsCurNMT,
					SYNC:          resultsCurSYNC,
					TIME:          resultsCurTIME,
					EMCY:          resultsCurEMCY,
				}
				// Enviar payload json da mensagem
				err := MQTTSendMessageAutomatic(msg)
				if err != nil {
					logPkg.CtsLog.Error("routineRead:FAIL dev name[%s]  Ip[%s:%s] Changes Occurred - MQTTSendMessageAutomatic\n err[%s]ms\n", dev.Name, dev.Address, dev.Port, err)
				}
				logPkg.CtsLog.Debug("routineRead: dev name[%s]  Ip[%s:%s] Changes Occurred - MQTTSendMessageAutomatic\n msg%s\n", dev.Name, dev.Address, dev.Port, msg)
			}

		}
		if netConn == nil {
			// Trying to connect to CANopen TCP Server
			netConn, err := CANopenTcpClient_ConnectToServer(netConn, dev.Address, dev.Port)
			if err == nil {
				// Successfull Client
				if !clientIsConnected {
					clientIsConnected = true
					logPkg.CtsLog.Warn("ContinuousRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]Connected\n",
						dev.Name, dev.Protocol, dev.Address, dev.Port)
					// Sending Info Log Device is Connected
					MQTTSendGeneralDeviceStatus(dev, "connected")
					msgLog := utilsPkg.MqttLogMsg{
						Address:  dev.Address,
						Port:     dev.Port,
						Name:     dev.Name,
						Protocol: dev.Protocol,
						Id:       dev.Id,
						Time:     readTime,
						Message:  "connected",
						Code:     "OK",
						Source:   "ContinuousRead:",
					}
					_ = MQTTSendLogMessage(msgLog, "info")
				}
			} else {
				//Fail to Connect with Device
				clientReadErrorReported = false
				if netConn != nil {
					netConn.Close()
					netConn = nil
				}
				if clientIsConnected {
					clientIsConnected = false
					logPkg.CtsLog.Error("DeviceName[%s]\n Protocol[%s]\n Ip:Port[%s:%s]Disconnected\n",
						dev.Name, dev.Protocol, dev.Address, dev.Port)
					// Sending Status that Device is connected
					MQTTSendGeneralDeviceStatus(dev, "disconnected")
					// Sending Error Log Device Disconnected
					msgLog := utilsPkg.MqttLogMsg{
						Address:  dev.Address,
						Port:     dev.Port,
						Name:     dev.Name,
						Protocol: dev.Protocol,
						Id:       dev.Id,
						Time:     readTime,
						Message:  "disconnected",
						Code:     "CANOPEN_400",
						Source:   "continuousRead:",
					}
					_ = MQTTSendLogMessage(msgLog, "error")
				}
				time.Sleep(5 * time.Second)
			}
		}
		anyChangeDetected = false
		if !bInfiniteLoop {
			// Single thread as WriteRequests
			break
		}
		time.Sleep(time.Duration(dev.ReadingTime) * time.Millisecond)
	} // End for {} infinite loop
	if netConn != nil {
		netConn.Close()
		netConn = nil
	}
	logPkg.CtsLog.Debug("routineRead: dev name[%s]  Address[%s%s] Exiting\n", dev.Name, dev.Address, dev.Port)
}

// func CheckConnection(ipAddress string, port int) bool {
func CheckConnection(conn net.Conn, ipAddress, port string) bool {

	tcpAddress := ipAddress + ":" + port

	if !isTCPConnected(conn, ipAddress, port) {
		if utilsPkg.AnyDeviceIsConnected {
			utilsPkg.AnyDeviceIsConnected = false
			logPkg.CtsLog.Error("CheckConnection: FAIL - tcpAddress[%s]Disconnected", tcpAddress)
		}
		return false
	}
	// logPkg.CtsLog.Warn("CheckConnection: PASS - tcpAddress[%s] Connected", tcpAddress)
	iCont++
	if iCont == 60 {
		// Every 60 good connections one will fail
		iCont = 0
		return false
	}
	return true
}

// func isTCPConnected(host string, port int) bool {
func isTCPConnected(conn net.Conn, host, port string) bool {
	tcpAddress := host + ":" + port
	if conn != nil {
		if !utilsPkg.AnyDeviceIsConnected {
			utilsPkg.AnyDeviceIsConnected = true
			logPkg.CtsLog.Debug("isTCPConnected:WARN tcpAddress[%s] Already Connected conn[%v]", tcpAddress, conn)
		}
		return true
	}
	if utilsPkg.AnyDeviceIsConnected {
		utilsPkg.AnyDeviceIsConnected = false
		logPkg.CtsLog.Debug("isTCPConnected:WARN tcpAddress[%s] Disconnected conn[%v]", tcpAddress, conn)
	}
	return false
}

// func changProtocolFormat(dev utilsPkg.DevSettings, mode string) DeviceInfo {
// Specific Protocol >>> protocol_config.json <<<
func changProtocolFormat(dev utilsPkg.DevSettings) DeviceInfo {

	var SDO []CanOpenSDOInfo
	var PDO []CanOpenPDOInfo
	var NMT []CanOpenNMTInfo
	var SYNC []CanOpenSYNCInfo
	var TIME []CanOpenTIMEInfo
	var EMCY []CanOpenEMCYInfo

	for _, data := range dev.Data {

		sdo, _ := data["SDO"].([]interface{})
		pdo, _ := data["PDO"].([]interface{})
		nmt, _ := data["NMT"].([]interface{})
		sync, _ := data["SYNC"].([]interface{})
		time, _ := data["TIME"].([]interface{})
		emcy, _ := data["EMCY"].([]interface{})

		for _, item := range sdo {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			cmd, _ := dataMap["cmd"].(float64)
			objIndex, _ := dataMap["objIndex"].(float64)
			subIndex, _ := dataMap["subIndex"].(float64)
			canOpenInfo := CanOpenSDOInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
				Cmd:        uint8(cmd),
				ObjIndex:   uint32(objIndex),
				SubIndex:   uint8(subIndex),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:SDO name[%s] canid[%s] nodeid[%d] cmd[%d] objIndex[%d] subIndex[%d]\n", name, canid, nodeid, cmd, objIndex, subIndex)
			SDO = append(SDO, canOpenInfo)
		}

		for _, item := range pdo {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			mapping, _ := dataMap["mapping"].([]interface{})
			mapping1, _ := dataMap["mapping"].([]CanOpenPDOMapping)
			for _, itemmapping := range mapping {
				dataMaping, _ := itemmapping.(map[string]interface{})
				objindex, _ := dataMaping["objIndex"].(float64)
				subindex, _ := dataMaping["subIndex"].(float64)
				mappingPdo := CanOpenPDOMapping{
					ObjIndex: uint32(objindex),
					SubIndex: uint8(subindex),
				}
				mapping = append(mapping, mappingPdo)
				/// logPkg.CtsLog.Debug("changProtocolFormat:PDO objindex[%d] subindex[%d] mapping[%v]\n", objindex, subindex, mapping)
			}
			canOpenInfo := CanOpenPDOInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
				MAPPING:    mapping1,
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:PDO name[%s] canid[%s] mapping1[%v]\n", name, canid, mapping1)
			PDO = append(PDO, canOpenInfo)
		}

		for _, item := range nmt {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			cmd, _ := dataMap["command"].(float64)
			canOpenInfo := CanOpenNMTInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
				Command:    uint8(cmd),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:NMT utem[%s] name[%s] canid[%s] nodeid[0x%X] cmd[%d] NMT\n%v\n", item, name, canid, nodeid, cmd, NMT)
			NMT = append(NMT, canOpenInfo)
		}

		for _, item := range sync {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			canOpenInfo := CanOpenSYNCInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:SYNC utem[%s] name[%s] canid[%s] nodeid[0x%X] NMT\n%v\n", item, name, canid, nodeid, NMT)
			SYNC = append(SYNC, canOpenInfo)
		}

		for _, item := range time {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			canOpenInfo := CanOpenTIMEInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:SYNC utem[%s] name[%s] canid[%s] nodeid[0x%X] NMT\n%v\n", item, name, canid, nodeid, NMT)
			TIME = append(TIME, canOpenInfo)
		}

		for _, item := range emcy {
			dataMap, _ := item.(map[string]interface{})
			name, _ := dataMap["name"].(string)
			canid, _ := dataMap["canid"].(string)
			nodeid, _ := dataMap["nodeid"].(float64)
			canOpenInfo := CanOpenEMCYInfo{
				Name:       name,
				CANOPEN_ID: canid,
				NODE_ID:    uint8(nodeid),
			}
			/// logPkg.CtsLog.Debug("changProtocolFormat:EMCY utem[%s] name[%s] canid[%s] nodeid[0x%X] NMT\n%v\n", item, name, canid, nodeid, NMT)
			EMCY = append(EMCY, canOpenInfo)
		}

	}

	deviceNewInfo := DeviceInfo{
		SDO:  SDO,
		PDO:  PDO,
		NMT:  NMT,
		SYNC: SYNC,
		TIME: TIME,
		EMCY: EMCY,
	}
	/// logPkg.CtsLog.Debug("changProtocolFormat:deviceNewInfo:\n%v\n", deviceNewInfo)
	return deviceNewInfo
}

// UpdateDevConfig will get all informations in file
func UpdateDevConfig(name string, dataMem DeviceInfo) (utilsPkg.DevSettings, []string, error) {
	var devInf []utilsPkg.Devices
	var itemInf utilsPkg.DevSettings
	var devConfig string = "deviceConfig.json"
	var memChanged []string

	file, err := os.ReadFile(devConfig)
	if err != nil {
		logPkg.CtsLog.Error("FSDOl to read JSON File: ", err)
		return itemInf, memChanged, err
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		logPkg.CtsLog.Error("FSDOL to decode JSON file")
		logPkg.CtsLog.Error("%v", err)
		return itemInf, memChanged, err
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			if item.Name == name {
				newDeviceInfo := changProtocolFormat(item)
				memChanged = checkcanOpenInfoChange(dataMem, newDeviceInfo)
				return item, memChanged, nil
			}
		}
	}
	return itemInf, memChanged, err
}

func checkcanOpenInfoChange(dataMem, newDeviceInfo DeviceInfo) []string {
	var rescanOpenInfoChanged []string

	// Verifica se ocorreu eliminação de alguma memória Bit ou Word
	for _, mem := range dataMem.NMT {
		notfound := true
		for _, newMem := range newDeviceInfo.NMT {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			rescanOpenInfoChanged = append(rescanOpenInfoChanged, "DI "+mem.Name+" deleted")
		}
	}
	for _, mem := range dataMem.SYNC {
		notfound := true
		for _, newMem := range newDeviceInfo.SYNC {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			rescanOpenInfoChanged = append(rescanOpenInfoChanged, "SYNC "+mem.Name+" deleted")
		}
	}
	// Verifica se ocorreu adição de alguma memória Bit ou Word
	for _, mem := range newDeviceInfo.NMT {
		notfound := true
		for _, newMem := range dataMem.NMT {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			rescanOpenInfoChanged = append(rescanOpenInfoChanged, "NMT "+mem.Name+" added")
		}
	}
	for _, mem := range newDeviceInfo.SYNC {
		notfound := true
		for _, newMem := range dataMem.SYNC {
			if mem.Name == newMem.Name {
				notfound = false
			}
		}
		if notfound {
			rescanOpenInfoChanged = append(rescanOpenInfoChanged, "NMT "+mem.Name+" added")
		}
	}
	return rescanOpenInfoChanged
}
