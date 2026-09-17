package main

import (
	utilsPkg "emuethernetip/utils"
	logPkg "emuethernetip/utils/gologtofile"
	"encoding/binary"
	"encoding/hex"
	"os"
	"runtime"
	"testing"
)

var bLogStarted = false

func Test01_procCipEipClientRquPacket_NullPacket(t *testing.T) {
	startLog()
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(nil)
	if err == nil {
		logPkg.CtsLog.Error("Test01:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL expected err because cipEipClientRspPkt = nil", resultCode)
	} else {
		logPkg.CtsLog.Info("Test01:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n clientRquPacket=nil\n cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
		t.Logf("resultCode[0x%02X]PASS clientRquPacket=nil cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test02_procCipEipClientRquPacket_ShortPacket(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 6500 0000 00000000 00000000
	cipEipRquData := []byte{} // Your request data

	// CIP Header Ethernet/IP Cmd
	cipEipCmd := uint16(0x0065)
	// CIP Header Ethernet/IP RquMsgLen
	cipEipRquDataLen := uint16(len(cipEipRquData))
	// CIP Header Ethernet/IP SenderHandler
	cipEipSenderHdl := uint32(0x00000000)
	// CIP Header Ethernet/IP Status
	cipEipStatus := uint32(0x00000000)
	// CIP Header Ethernet/IP SenderContext
	// cipEipSenderContext := uint64(0x0000000000000000)
	// CIP Header Ethernet/IP Option
	// cipEipOptions := uint32(0x00000000)

	// CIP Packet Ethernet/IP = CIP Header + CIP Ethernet/IP Request Data
	cipEipClientRquPkt := make([]byte, 16)

	// CIP: EthernetIP Header
	binary.LittleEndian.PutUint16(cipEipClientRquPkt[0:2], cipEipCmd)        //  0-1  cipEipCmd
	binary.LittleEndian.PutUint16(cipEipClientRquPkt[2:4], cipEipRquDataLen) //  2-3  cipEipRquDataLen
	binary.LittleEndian.PutUint32(cipEipClientRquPkt[4:8], cipEipSenderHdl)  //  4-7  cipEipSenderHdl
	binary.LittleEndian.PutUint32(cipEipClientRquPkt[8:12], cipEipStatus)    //  8-11 cipEipStatus
	// binary.LittleEndian.PutUint64(cipEipClientRquPkt[12:20], cipEipSenderContext) // 12-19 cipEipSenderContext
	//binary.LittleEndian.PutUint32(cipEipClientRquPkt[20:24], cipEipOptions)       // 20-23 cipEipSenderContext
	// CIP: EthernetIP Data
	// copy(cipEipClientRquPkt[24:], cipEipRquData)
	logPkg.CtsLog.Info("Test02: >> cipEipClientRquPkt[%d:%X]ShortPacket", len(cipEipClientRquPkt), cipEipClientRquPkt)

	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err == nil {
		logPkg.CtsLog.Error("Test02:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test02:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n clientRquPacket=nil\n cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test03_procCipEipClientRquPacket_InvalidDataLen(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 6500 0500 00000000 00000000 0000000000000000 00000000 01000000
	cipEipRquData := []byte{0x01, 0x00, 0x00, 0x00} // Your request data

	// CIP Header Ethernet/IP Cmd
	cipEipCmd := uint16(0x0065)
	// CIP Header Ethernet/IP RquMsgLen
	cipEipRquDataLen := uint16(len(cipEipRquData) + 1)
	// CIP Header Ethernet/IP SenderHandler
	cipEipSenderHdl := uint32(0x00000000)
	// CIP Header Ethernet/IP Status
	cipEipStatus := uint32(0x00000000)
	// CIP Header Ethernet/IP SenderContext
	cipEipSenderContext := uint64(0x0000000000000000)
	// CIP Header Ethernet/IP Option
	cipEipOptions := uint32(0x00000000)

	// CIP Packet Ethernet/IP = CIP Header + CIP Ethernet/IP Request Data
	cipEipClientRquPkt := make([]byte, 24+len(cipEipRquData))

	// CIP: EthernetIP Header
	binary.LittleEndian.PutUint16(cipEipClientRquPkt[0:2], cipEipCmd)             //  0-1  cipEipCmd
	binary.LittleEndian.PutUint16(cipEipClientRquPkt[2:4], cipEipRquDataLen)      //  2-3  cipEipRquDataLen
	binary.LittleEndian.PutUint32(cipEipClientRquPkt[4:8], cipEipSenderHdl)       //  4-7  cipEipSenderHdl
	binary.LittleEndian.PutUint32(cipEipClientRquPkt[8:12], cipEipStatus)         //  8-11 cipEipStatus
	binary.LittleEndian.PutUint64(cipEipClientRquPkt[12:20], cipEipSenderContext) // 12-19 cipEipSenderContext
	binary.LittleEndian.PutUint32(cipEipClientRquPkt[20:24], cipEipOptions)       // 20-23 cipEipSenderContext
	// CIP: EthernetIP Data
	copy(cipEipClientRquPkt[24:], cipEipRquData)
	logPkg.CtsLog.Info("Test03: >> cipEipClientRquPkt[%d:%X]InvalidDataLen", len(cipEipClientRquPkt), cipEipClientRquPkt)

	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err == nil {
		logPkg.CtsLog.Error("Test03:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test03:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n clientRquPacket=nil\n cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test04_procCipEipClientRquPacket_InvalidCmd(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 0000 1A00 00000000 00000000 FFFFFFFF00000000 00000000 000000000A00020000000000B2000A000E042100510324023001
	strCipEipClientRquPkt := "00001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324023001"

	// CIP: EthernetIP Conv strCipEipClientRquPkt > []byte cipEipClientRquPkt
	cipEipClientRquPkt, err := hex.DecodeString(strCipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test04:conversion of Hex str to hex []byte(%s)fail", hex.DecodeString)
		t.Errorf("conversion of Hex str to hex []byte(%s)fail", strCipEipClientRquPkt)
	}
	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err == nil {
		logPkg.CtsLog.Error("Test04:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test04:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n >> cipEipClientRquPkt[%X]\n << cipEipClientRspPkt[%X]\n", resultCode, cipEipClientRquPkt, cipEipClientRspPkt)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test05_procCipEipClientRquPacket_RegisterSessionOK(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 6500 0400 00000000 00000000 0000000000000000 00000000 01000000
	strCipEipClientRquPkt := "65000400000000000000000000000000000000000000000001000000"

	// CIP: EthernetIP Conv strCipEipClientRquPkt > []byte cipEipClientRquPkt
	cipEipClientRquPkt, err := hex.DecodeString(strCipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test05:conversion of Hex str to hex []byte(%s)fail", hex.DecodeString)
		t.Errorf("conversion of Hex str to hex []byte(%s)fail", strCipEipClientRquPkt)
	}
	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test05:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test05:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n >> cipEipClientRquPkt[%X]\n << cipEipClientRspPkt[%X]\n", resultCode, cipEipClientRquPkt, cipEipClientRspPkt)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test06_procCipEipClientRquPacket_SendRRData_Cls848_Ins01_Att01OK(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 6F00 1A00 00000000 00000000 FFFFFFFF00000000 00000000 000000000A00020000000000B2000A000E042100500324013001
	strCipEipClientRquPkt := "6F001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100500324013001"

	// CIP: EthernetIP Conv strCipEipClientRquPkt > []byte cipEipClientRquPkt
	cipEipClientRquPkt, err := hex.DecodeString(strCipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test06:conversion of Hex str to hex []byte(%s)fail", hex.DecodeString)
		t.Errorf("conversion of Hex str to hex []byte(%s)fail", strCipEipClientRquPkt)
	}
	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test06:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test06:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n >> cipEipClientRquPkt[%X]\n << cipEipClientRspPkt[%X]\n", resultCode, cipEipClientRquPkt, cipEipClientRspPkt)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func Test07_procCipEipClientRquPacket_SendRRData_Cls849_Ins02_Att01OK(t *testing.T) {
	startLog()
	// 0 1  2 3  4 5 6 7  8 9 0 1  2 3 4 5 6 7 8 9  0 1 2 3  4 5 6 7
	// Cmd  Len  SesHandl Status    SenderContext   Options  Data
	// 6F00 1A00 00000000 00000000 FFFFFFFF00000000 00000000 000000000A00020000000000B2000A000E042100510324023001
	strCipEipClientRquPkt := "6F001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324023001"

	// CIP: EthernetIP Conv strCipEipClientRquPkt > []byte cipEipClientRquPkt
	cipEipClientRquPkt, err := hex.DecodeString(strCipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test07:conversion of Hex str to hex []byte(%s)fail", hex.DecodeString)
		t.Errorf("conversion of Hex str to hex []byte(%s)fail", strCipEipClientRquPkt)
	}
	// CIP: EthernetIP Proccessing Client Request Packet
	resultCode, cipEipClientRspPkt, err := procCipEipClientRquPacket(cipEipClientRquPkt)
	if err != nil {
		logPkg.CtsLog.Error("Test07:FAIL procCipEipClientRquPacket\n resultCode[0x%02X]FAIL\n Expected err\nReceived cipEipClientRspPkt err:%s", resultCode, err)
		t.Errorf("resultCode[0x%02X]FAIL", resultCode)
	} else {
		logPkg.CtsLog.Info("Test07:PASS procCipEipClientRquPacket\n resultCode[0x%02X]\n >> cipEipClientRquPkt[%X]\n << cipEipClientRspPkt[%X]\n", resultCode, cipEipClientRquPkt, cipEipClientRspPkt)
		t.Logf("resultCode[0x%02X]PASS cipEipClientRspPkt[%X]\nerr:%s", resultCode, cipEipClientRspPkt, err)
	}
}

func startLog() {
	if !bLogStarted {
		utilsPkg.Hostname, _ = os.Hostname()
		// Initialize CtsLogs with default parameters
		rc, err := logPkg.InitCtsLogs(
			utilsPkg.LogFilePath,
			"emuethernetip_test",
			utilsPkg.DeleteExistingLogFiles,
			utilsPkg.LogLevel,
			utilsPkg.LogFileMaxSize,
		)
		if (err != nil) || (rc < 0) {
			/// Fail to setup Logs
			logPkg.CtsLog.Error("Test01:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d] MaxSize[%d]\n         err[%s]\n\n",
				rc,
				runtime.GOOS,
				runtime.GOARCH,
				utilsPkg.OsBits,
				utilsPkg.Hostname,
				utilsPkg.LogFilePath,
				utilsPkg.LogFileName,
				utilsPkg.LogLevel,
				utilsPkg.LogFileMaxSize,
				err)
			return
		}
		logPkg.CtsLog.Info("Test01:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s] MaxSize[%d]\n    LogLevel[%d]\n\n",
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.LogFilePath,
			utilsPkg.LogFileName,
			utilsPkg.LogLevel,
			utilsPkg.LogFileMaxSize)
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
