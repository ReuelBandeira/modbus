package newprotocol

import (
	utilsPkg "newprotocol/utils"
	logPkg "newprotocol/utils/gologtofile"
	"os"
	"runtime"
	"testing"
)

func Test01_proccessClientRquPacket_NullPacket(t *testing.T) {

	utilsPkg.Hostname, _ = os.Hostname()
	// Initialize CtsLogs with default parameters
	rc, err := logPkg.InitCtsLogs(
		"",
		"newprototocol_test",
		false,
		3,
		500.000,
	)
	if (err != nil) || (rc < 0) {
		/// Fail to setup Logs
		/// TODO: Channel to Send Fail Information to Backend
		logPkg.CtsLog.Error("main:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n         err[%s]\n\n",
			rc,
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			"",
			"newprototocol_test",
			3,
			err)
		return
	}
	logPkg.CtsLog.Info("main:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n\n",
		runtime.GOOS,
		runtime.GOARCH,
		utilsPkg.OsBits,
		utilsPkg.Hostname,
		"",
		"newprototocol_test",
		3)

	// Force Warning Log Levels Onpy
	logPkg.CtsLog.SetLogLevel(3)

	retui16val, err := HexStringToUint16("0x600")
	if err != nil {
		/// FAIL to Connect to Internal MQTT Brocker
		logPkg.CtsLog.Error("main:Test01.1 FAIL HexStringToUint16(\"0x600\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	} else {
		logPkg.CtsLog.Info("main:Test01.1 PASS HexStringToUint16(\"0x600\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	}
	retui16val, err = HexStringToUint16("600")
	if err != nil {
		/// FAIL to Connect to Internal MQTT Brocker
		logPkg.CtsLog.Error("main:Test01.2 FAIL HexStringToUint16(\"600\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	} else {
		logPkg.CtsLog.Info("main:Test01.1 PASS HexStringToUint16(\"600\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	}
	retui16val, err = HexStringToUint16("FEDC")
	if err != nil {
		/// FAIL to Connect to Internal MQTT Brocker
		logPkg.CtsLog.Error("main:Test01.3 FAIL HexStringToUint16(\"FEDC\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	} else {
		logPkg.CtsLog.Info("main:Test01.3 PASS HexStringToUint16(\"FEDC\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	}
	retui16val, err = HexStringToUint16("fedc")
	if err != nil {
		/// FAIL to Connect to Internal MQTT Brocker
		logPkg.CtsLog.Error("main:Test01.4 FAIL HexStringToUint16(\"fedc\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	} else {
		logPkg.CtsLog.Info("main:Test01.4 PASS HexStringToUint16(\"fedc\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	}
	retui16val, err = HexStringToUint16("joao")
	if err != nil {
		/// FAIL to Connect to Internal MQTT Brocker
		logPkg.CtsLog.Error("main:Test01.5 PASS HexStringToUint16(\"joao\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	} else {
		logPkg.CtsLog.Info("main:Test01.5 FAIL HexStringToUint16(\"joao\") retui16val[0x%04X]\n err[%v]\n", int(retui16val), err)
	}

	resultCode, clientRspPkt, err := proccessClientRquPacket(nil)
	if err == nil {
		t.Errorf("resultCode[0x%X]FAIL expected err because ClientRquPacket = nil", resultCode)
	} else {
		t.Logf("resultCode[0x%X]PASS clientRquPacket=nil clientRspPkt[%X]\nerr[%s]", resultCode, clientRspPkt, err)
	}
}

func Test02_proccessClientRquPacket_ShortPacket(t *testing.T) {
	ClientRquPacket := []byte{0x01, 0x02}
	resultCode, clientRspPkt, err := proccessClientRquPacket(ClientRquPacket)
	if err == nil {
		t.Errorf("resultCode[0x%X]FAIL expected err because ClientRquPacket is short", resultCode)
	} else {
		t.Logf("resultCode[0x%X]PASS clientRquPacket=nil clientRspPkt[%X]\nerr[%s]", resultCode, clientRspPkt, err)
	}
}

func Test03_proccessClientRquPacket_InvalidPacketLen(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	resultCode, clientRspPkt, err := proccessClientRquPacket(ClientRquPacket)
	if err == nil {
		t.Errorf("resultCode[0x%X]FAIL expected err because ClientRquPacket has invalid len", resultCode)
	} else {
		t.Logf("resultCode[0x%X]PASS clientRspPkt[%X] Expected InvalidPacket \nerr[%s]", resultCode, clientRspPkt, err)
	}
}

func Test04_proccessClientRquPacket_MissingSTX(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	resultCode, clientRspPkt, err := proccessClientRquPacket(ClientRquPacket)
	if err == nil {
		t.Errorf("resultCode[0x%X] FAIL clientRquPacket=nil ClientRquPacket[0][0x%02X] Expected[0x02]NOP\nerr[%s]", resultCode, clientRspPkt, err)
	} else {
		t.Logf("resultCode[0x%X]PASS clientRquPacket=nil ClientRquPacket[0][0x%02X] Expected[0x00]NOP\nerr[%s]", resultCode, clientRspPkt, err)
	}
}

func Test05_proccessClientRquPacket_MissingETX(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00}
	resultCode, clientRspPkt, err := proccessClientRquPacket(ClientRquPacket)
	if err == nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[int(ClientRquPacket[2])+3]=[0x%02X] Expected[0x00]ETX clientRspPkt[%x]\nerr[%s]", resultCode, ClientRquPacket[int(ClientRquPacket[2])+3], clientRspPkt, err)
	} else {
		t.Logf("resultCode[0x%X]PASS ClientRquPacket[int(ClientRquPacket[2])+3]=[0x%02X] Expected[0x00]", resultCode, ClientRquPacket[int(ClientRquPacket[2])+3])
	}
}

func Test06_proccessClientRquPacket_GoodSTX_and_ETX_CRC16_MISSMATCH(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x00, 0x00, 0x03, 0x00, 0x00}
	resultCode, clientRspPkt, err := proccessClientRquPacket(ClientRquPacket)
	if err == nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[int(ClientRquPacket[2])+3]=[0x%02X] Expected[0x03]ETX clientRspPkt[%x]\nerr[%s]", resultCode, ClientRquPacket[int(ClientRquPacket[2])+3], clientRspPkt, err)
	} else {
		t.Logf("resultCode[0x%X]PASS ClientRquPacket[int(ClientRquPacket[2])+3]=[0x%02X] Expected[0x03]ETX", resultCode, ClientRquPacket[int(ClientRquPacket[2])+3])
	}
}

func Test07_proccessClientRquPacket_GoodSTX_and_ETX_CRC16_OK(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x00, 0x00, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != 0x00 {
				resultCode = clientRspPacket[3]
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x00]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]OK:ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test08_proccessClientRquPacket_invalidCmd(t *testing.T) {
	//                        STX   Cmd   Len   ETX   CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0xAA, 0x00, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != ERR_INVALID_CMD {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if resultCode != ERR_INVALID_CMD {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test09_proccessClientRquPacket_CmdRead_DI_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL  ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x01, 0x02, 0x05, 0x67, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test10_proccessClientRquPacket_CmdRead_DO_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL  ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x02, 0x02, 0x05, 0x02, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test11_proccessClientRquPacket_CmdRead_AI_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL  ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x03, 0x02, 0x07, 0xEF, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test12_proccessClientRquPacket_CmdRead_AO_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL  ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x04, 0x02, 0x07, 0xE0, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test13_proccessClientRquPacket_CmdWrite_DO_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL Data   ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x05, 0x03, 0x07, 0xE0, 0x0F, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

func Test14_proccessClientRquPacket_CmdWrite_AO_OK(t *testing.T) {
	//                        STX   Cmd   Len   AddrH ADDrL DataH DAtaL  ETX  CRCH  CRCL
	ClientRquPacket := []byte{0x02, 0x06, 0x04, 0x07, 0xE2, 0x0F, 0x05, 0x03, 0x00, 0x00}
	clientRquPacketWoCRC16 := ClientRquPacket[:len(ClientRquPacket)-2]
	ClientRqhCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	ClientRquPacket[len(ClientRquPacket)-2] = uint8(ClientRqhCrc16Calc >> 8)
	ClientRquPacket[len(ClientRquPacket)-1] = uint8(ClientRqhCrc16Calc & 0xFF)
	resultCode, clientRspPacket, err := proccessClientRquPacket(ClientRquPacket)
	if err != nil {
		t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%X] ClientRqhCrc16Calc[0x%04X] clientRspPkt[%X] err[%s]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, err)
	} else {
		clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
		var ClientRspCrc16Rcvd uint16 = 0
		ClientRspCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
		ClientRspCrc16Rcvd = uint16(clientRspPacket[len(clientRspPacket)-2]) << 8
		ClientRspCrc16Rcvd += uint16(clientRspPacket[len(clientRspPacket)-1])
		if ClientRspCrc16Calc != ClientRspCrc16Rcvd {
			resultCode = ERR_BAD_CRC16
			t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]ER:MISSMATCH ClientRspCrc16Rcvd[0x%04X]", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd)
		} else {
			if clientRspPacket[3] != OK {
				t.Errorf("resultCode[0x%X]FAIL ClientRquPacket[%0X] crc16Calc[0x%04X]ER clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ER expected[0x01]Invalid Command", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
			} else {
				t.Logf("resultCode[0x%X]PASS ClientRquPacket[%0X] crc16Calc[0x%04X]OK clientRspPacket[%X] ClientRspCrc16Calc[0x%04X]OK:MATCH ClientRspCrc16Rcvd[0x%04X] clientRspPacket[3]=[0x%02X]ResultCode=0x00(Success)", resultCode, ClientRquPacket, ClientRqhCrc16Calc, clientRspPacket, ClientRspCrc16Calc, ClientRspCrc16Rcvd, clientRspPacket[3])
				resultCode, clientRspData, err := unpackClientRspPacket(clientRspPacket)
				if err != nil {
					t.Errorf("resultCode[0x%X]FAIL\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]\n                err[%s]", resultCode, ClientRquPacket, clientRspPacket, clientRspData, err)
				} else {
					t.Logf("resultCode[0x%X]PASS\n >> ClientRquPacket[%X]\n << clientRspPacket[%X]\n <<   clientRspData[%X]", resultCode, ClientRquPacket, clientRspPacket, clientRspData)
				}
			}
		}
	}
}

/*
// BEGIN: TABLE DRIVEN TEST EXAMPLE
func Test_TableDriven(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name     string // FIELD 1: TestName
		input    int    // FIELD 2: Input Value to be testesd
		expected string // FIELD 3: Desired Result when function is called
	}{
		// the table itself
		{"9 should be OK", 9, "OK"},
		{"3 should be OK", 3, "OK"},
		{"1 is not  OK", 1, "NOK"},
		{"0 should be OK", 0, "OK"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultGot := FunctionToBeTested(tt.input)
			if resultGot != tt.expected {
				t.Errorf("tt.input[%d] resultGot %s, Expected %s", tt.input, resultGot, tt.expected)
			}
		})
	}
}

// END: TABLE DRIVEN TEST EXAMPLE

// BEGIN: FUNCTION TO BE TESTED
func FunctionToBeTested(input int) string {
	var dataout string = "NONE"
	switch input {
	case 9:
		dataout = "OK"
	case 3:
		dataout = "OK"
	case 0:
		dataout = "OK"
	case 1:
		dataout = "NOK"
	default:
	}
	return dataout
}

//    END: FUNCTION TO BE TESTED
*/
