package newprotocol

import (
	/// TODO: change newprotocol with your protocolname
	"fmt"
	"hash/crc32"
	"net"
	utilsPkg "newprotocol/utils"
	golog "newprotocol/utils/gologtofile"
	"strconv"
)

// Digital & Analog Inouts and Outputs Local Buffer
var DigitalOutputs [32768]byte // Max 32k Bytes = 257.344 Digital Outputs
var DigitalInputs [32768]byte  // Max 32k Bytes = 257.344 Digital Inputs
var AnalogOutputs [4096]uint16 // Max  4k Butes =   2.048 Analog  Outputs
var AnalogInputs [4096]uint16  // Max  4k Bytes =   2.048 Analog  Inputs

// New Protocol Constants
// ===================================newprotocol: Packet Format =============================
// RquMsg >> STX(1) CmdCode(1) RquPayLoadLen(1) <-RquPayLoadData(n)-------------> ETX(1) CRC16(2)
// RspMsg << STX(1) CmdCode(1) RspPayLoadLen(1) <-RspPayLoadData(n)-------------> ETX(1) CRC16(2)
// -------------------------------------------- <-ResultCode(1), ResultData(n)-->
// ===========================================================================================
const (
	// newprotocol Control Characters
	CNOP = 0x00
	CSTX = 0x02
	CETX = 0x03
	// NewProotocol CmdCodes
	CMD_NOP   = 0x00
	CMD_RD_DI = 0x01
	CMD_RD_DO = 0x02
	CMD_RD_AI = 0x03
	CMD_RD_AO = 0x04
	CMD_WR_DO = 0x05
	CMD_WR_AO = 0x06
	// newprotocol ResultCodes
	OK = 0x00 // Success
	// newprotocol Device Package ErrorCodes
	ERR_BAD_PACKET_FORMAT        = 0x81 // Package with Wrong Format
	ERR_BAD_CRC16                = 0x82 // Package with bad CRC16
	ERR_INVALID_CMD              = 0x83 // Package with Invelid Command
	ERR_INVALID_PARAMETER        = 0x84 // Package with Invalid parameter
	ERR_INVALID_VALUE            = 0x85 // Invelid value
	ERR_DIO_ADDR_ABOVE_MAX_LIMIT = 0x86 // Digital IO Above max address limit
	ERR_AIO_ADDR_ABOVE_MAX_LIMIT = 0x91 // Analog  IO Above max address limit
	// newprotocol Communication ErroCodes
	ERR_NET_INTERFACES    = 0x92 // net could not find any Avalable Interface
	ERR_NET_IFC_ADDRS     = 0x93 // net with invalid interface address
	ERR_NET_PARSE_CIDR    = 0x94 // net with invalid CIDR
	ERR_NET_LISTEN        = 0x95 // net fail to listen
	ERR_NET_LISTEN_ACCEPT = 0x96 // net fail to accept connection
	ERR_COMM_WRITE        = 0xA1 // net fail to write(Send Packet)
	ERR_COMM_READ         = 0xA2 // net fail to  Read(Recv Packet)
)

// ProtocolConnect make a connection with a Device
// / TODO: Replace it to your specific protocol struct
func ProtocolConnect(device utilsPkg.DevSettings) (bool, net.Conn) {
	conn, err := TCPDeviceConnect(device.Address, device.Port)
	if err != nil {
		golog.CtsLog.Error("ProtocolConnect:FAIL Device[%s:%s]Disconnected", device.Address, device.Port)
		utilsPkg.AnyDeviceIsConnected = false
		return false, nil
	}
	golog.CtsLog.Debug("ProtocolConnect:OK Device[%s:%s]Connected", device.Address, device.Port)
	utilsPkg.AnyDeviceIsConnected = true
	return true, conn
}

// Computing CRC16
// / TODO: Replace it to your specific protocol struct
func CalcCRC16(data []byte) uint16 {
	// Create a CRC16 table using the CRC32 package
	crcTable := crc32.MakeTable(crc32.IEEE)
	// Calculate the CRC16 checksum
	crc := crc32.Checksum(data, crcTable)
	// fmt.Printf("CRC16 Checksum: 0x%04X\n", crc)
	return uint16(crc)
}

// / TODO: Replace it to your specific protocol struct
func HexStringToUint16(hexStr string) (uint16, error) {
	// Remove the "0x" prefix if present
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Parse the hexadecimal string to uint64
	val, err := strconv.ParseUint(hexStr, 16, 16)
	if err != nil {
		golog.CtsLog.Error("HexStringToUint16: hexStr[%s] val[0x%X(%d)]", hexStr, val, val)
		return 0, err
	}

	// Convert the uint64 to uint16
	result := uint16(val)

	golog.CtsLog.Error("HexStringToUint16: hexStr[%s] result[0x%X(%d)]", hexStr, result, result)
	return result, nil
}

// RD DO
// / TODO: Replace it to your specific protocol struct
func ReadDigitalOutput(conn net.Conn, Address uint16) (byte, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x02
	var RspData []byte
	var resultCode byte = 0x80
	if conn == nil {
		return 0xFF, fmt.Errorf("ReadDigitalOutput: not Connected conn[%v]", conn)
	} else if !utilsPkg.AnyDeviceIsConnected {
		return 0xFF, fmt.Errorf("ReadDigitalOutput: Disconnected conn[%v] ", conn)
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)               //STX
	ClientRquPacket = append(ClientRquPacket, CMD_RD_DO)          //RD DO
	ClientRquPacket = append(ClientRquPacket, payLoadLen)         // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))   // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF)) // ADDRL
	ClientRquPacket = append(ClientRquPacket, CETX)               // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalOutput: Sending data err[%s]", err)
		return 0xFF, err
	}
	golog.CtsLog.Debug("ReadDigitalOutput:\n >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalOutput: Receiving data err[%s]", err)
		return 0xFF, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("ReadDigitalOutput:\n  << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalOutput: Error receiving data: resultCode[%02X] err[%s]", resultCode, err)
		return 0xFF, err
	}
	golog.CtsLog.Debug("ReadDigitalInput:Address[0x%04X]\n resultCode[%02X]\n RspData[%X]\n RspData[2]=[%02X]", Address, resultCode, RspData, RspData[2])
	return RspData[2], nil
}

// RD DIIsAnyDeviceIsConnected
// / TODO: Replace it to your specific protocol struct
func ReadDigitalInput(conn net.Conn, Address uint16) (byte, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x02
	var RspData []byte
	var resultCode byte = 0x80
	if conn == nil {
		return 0xFF, fmt.Errorf("ReadDigitalInput: not Connected conn[%v]", conn)
	} else if !utilsPkg.AnyDeviceIsConnected {
		return 0xFF, fmt.Errorf("ReadDigitalInput: Disconnected conn[%v]", conn)
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)               //STX
	ClientRquPacket = append(ClientRquPacket, CMD_RD_DI)          //RD DI
	ClientRquPacket = append(ClientRquPacket, payLoadLen)         // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))   // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF)) // ADDRL
	ClientRquPacket = append(ClientRquPacket, CETX)               // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalInput: Sending data err[%s]", err)
		return 0xFF, err
	}
	golog.CtsLog.Debug("ReadDigitalInput:\n >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalInput: Receiving data err[%s]", err)
		return 0xFF, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("ReadDigitalInput:\n << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadDigitalInput: Error receiving data: resultCode[%02X] err[%s]", resultCode, err)
		return 0xFF, err
	}
	golog.CtsLog.Debug("ReadDigitalInput:Address[0x%04X]\n resultCode[%02X]\n RspData[%X]\n RspData[2]=[%02X]", Address, resultCode, RspData, RspData[2])
	return RspData[2], nil
}

// RD AO
// / TODO: Replace it to your specific protocol struct
func ReadAnalogOutput(conn net.Conn, Address uint16) (uint16, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x02
	var RspData []byte
	var resultCode byte = 0x80
	var AnalogOutputValue uint16 = 0xFFFF
	if conn == nil {
		return AnalogOutputValue, fmt.Errorf("ReadAnalogOutput: not Connected conn[%v]", conn)
	} else if !utilsPkg.AnyDeviceIsConnected {
		return AnalogOutputValue, fmt.Errorf("ReadAnalogOutput: Disconnected conn[%v]", conn)
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)               //STX
	ClientRquPacket = append(ClientRquPacket, CMD_RD_AO)          //RD AO
	ClientRquPacket = append(ClientRquPacket, payLoadLen)         // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))   // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF)) // ADDRL
	ClientRquPacket = append(ClientRquPacket, CETX)               // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogOutput: Sending data err[%s]", err)
		return AnalogOutputValue, err
	}
	golog.CtsLog.Debug("ReadAnalogOutput:\n >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogOutput: Receiving data err[%s]", err)
		return AnalogOutputValue, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("ReadAnalogOutput:\n << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogOutput: Error receiving data: resultCode[%02X] err[%s]", resultCode, err)
		return AnalogOutputValue, err
	}
	AnalogOutputValue = uint16(RspData[2]) << 8
	AnalogOutputValue += uint16(RspData[3])

	golog.CtsLog.Debug("ReadAnalogOutput:Address[0x%04X]\n resultCode[%02X]\n RspData[%X]\n AnalogOutputValue=[%04X]", Address, resultCode, RspData, AnalogOutputValue)
	return AnalogOutputValue, nil
}

// RD AI
// / TODO: Replace it to your specific protocol struct
func ReadAnalogInput(conn net.Conn, Address uint16) (uint16, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x02
	var RspData []byte
	var resultCode byte = 0x80
	var AnalogInputValue uint16 = 0xFFFF
	if conn == nil {
		return AnalogInputValue, fmt.Errorf("ReadAnalogInput: not Connected conn[%v]", conn)
	} else if !utilsPkg.AnyDeviceIsConnected {
		return AnalogInputValue, fmt.Errorf("ReadAnalogInput: Disconnected conn[%v]", conn)
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)               //STX
	ClientRquPacket = append(ClientRquPacket, CMD_RD_AI)          //RD AI
	ClientRquPacket = append(ClientRquPacket, payLoadLen)         // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))   // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF)) // ADDRL
	ClientRquPacket = append(ClientRquPacket, CETX)               // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogInput: Sending data err[%s]", err)
		return AnalogInputValue, err
	}
	golog.CtsLog.Debug("ReadAnalogInput: >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogInput: Receiving data err[%s]", err)
		return AnalogInputValue, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("ReadAnalogInput:\n << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("ReadAnalogInput: Error receiving data: resultCode[%02X] err[%s]", resultCode, err)
		return AnalogInputValue, err
	}
	AnalogInputValue = uint16(RspData[2]) << 8
	AnalogInputValue += uint16(RspData[3])

	golog.CtsLog.Debug("ReadAnalogInput:Address[0x%04X]\n resultCode[%02X]\n RspData[%X]\n AnalogInputValue=[%04X]", Address, resultCode, RspData, AnalogInputValue)
	return AnalogInputValue, nil
}

// WR AO
// / TODO: Replace it to your specific protocol struct
func WriteDigitalOutput(conn net.Conn, Address uint16, turnOnValue uint8) (uint8, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x03
	var RspData []byte
	var resultCode byte = 0x80
	var ReturnnedValue uint8 = 0x00 // OFF
	if conn == nil {
		return ReturnnedValue, fmt.Errorf("WriteDigitalOutput: not Connected conn=nil")
	} else if !utilsPkg.AnyDeviceIsConnected {
		return ReturnnedValue, fmt.Errorf("WriteDigitalOutput: not Connected")
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)               //STX
	ClientRquPacket = append(ClientRquPacket, CMD_WR_DO)          //WR DO
	ClientRquPacket = append(ClientRquPacket, payLoadLen)         // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))   // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF)) // ADDRL
	ClientRquPacket = append(ClientRquPacket, byte(turnOnValue))  // VALUE
	ClientRquPacket = append(ClientRquPacket, CETX)               // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("WriteDigitalOOut: Sending data err[%s]", err)
		return ReturnnedValue, err
	}
	golog.CtsLog.Debug("WriteDigitalOutput:\n >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("WriteDigitalOOut: Receiving data err[%s]", err)
		return ReturnnedValue, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("WriteDigitalOutput:\n << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("WriteDigitalOutput: Error ParseClientRspPacket: resultCode[%02X] err[%s]", resultCode, err)
		return ReturnnedValue, err
	}
	ReturnnedValue = RspData[2]

	golog.CtsLog.Warn("WriteDigitalOutput:Address[0x%04X]\n resultCode[%02X]\n RspData[%X]\n ReturnnedValue=[%02X]", Address, resultCode, RspData, ReturnnedValue)
	return ReturnnedValue, nil
}

// WR AO
// / TODO: Replace it to your specific protocol struct
func WriteAnalogOutput(conn net.Conn, Address uint16, ValueToWrite uint16) (uint16, error) {
	var err error = nil
	var ClientRquPacket []byte
	var payLoadLen byte = 0x04
	var RspData []byte
	var resultCode byte = 0x80
	var wroteAnalogOutputValue uint16 = 0x0000
	if conn == nil {
		return wroteAnalogOutputValue, fmt.Errorf("WriteAnalogOutput: not Connected conn=nil")
	} else if !utilsPkg.AnyDeviceIsConnected {
		return wroteAnalogOutputValue, fmt.Errorf("WriteAnalogOutput: not Connected")
	}
	// Creating Packet
	ClientRquPacket = append(ClientRquPacket, CSTX)                    //STX
	ClientRquPacket = append(ClientRquPacket, CMD_WR_AO)               //WR AO
	ClientRquPacket = append(ClientRquPacket, payLoadLen)              // LEN
	ClientRquPacket = append(ClientRquPacket, byte(Address>>8))        // ADDRH
	ClientRquPacket = append(ClientRquPacket, byte(Address&0xFF))      // ADDRL
	ClientRquPacket = append(ClientRquPacket, byte(ValueToWrite>>8))   // VALH
	ClientRquPacket = append(ClientRquPacket, byte(ValueToWrite&0xFF)) // VALL
	ClientRquPacket = append(ClientRquPacket, CETX)                    // ETX
	// Appending CRC16
	ClientRqhCrc16Calc := CalcCRC16(ClientRquPacket)
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc>>8))   // CRCH
	ClientRquPacket = append(ClientRquPacket, byte(ClientRqhCrc16Calc&0xFF)) // CRCL

	// Sending Request Packet
	_, err = conn.Write(ClientRquPacket)
	if err != nil {
		golog.CtsLog.Error("WriteAnalogOutput: Sending data err[%s]", err)
		return wroteAnalogOutputValue, err
	}
	golog.CtsLog.Debug("WriteAnalogOutput:\n >> ClientRquPacket[%X]", ClientRquPacket)

	// Receiving Response Packet the response from the server into a byte array
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("WriteAnalogOutput: Receiving data err[%s]", err)
		return wroteAnalogOutputValue, err
	}
	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	golog.CtsLog.Debug("WriteAnalogOutput:\n << ClientRspPacket[%X]\n", ClientRspPacket)
	// TODO Unpacking Received packet too get response data
	resultCode, RspData, err = ParseClientRspPacket(ClientRspPacket)
	if err != nil {
		golog.CtsLog.Error("WriteAnalogOutput: Error receiving data: resultCode[%02X] err[%s]", resultCode, err)
		return wroteAnalogOutputValue, err
	}
	wroteAnalogOutputValue = uint16(RspData[2]) << 8
	wroteAnalogOutputValue += uint16(RspData[3])

	golog.CtsLog.Warn("WriteAnalogOutput:Address[0x%04X]\n resultCode[0x%02X]\n RspData[%X]\n wroteAnalogOutputValue=[%04X]", Address, resultCode, RspData, wroteAnalogOutputValue)
	return wroteAnalogOutputValue, nil
}

// ======================================
// unpackClientRspPacket - On Successfull
// ======================================
// 0  Hdr STX
// 1  Hdr Cmd
// 2  Hdr PayLoadLen
// 3  Pld 0 ResultCode
// 4-n PayloadData
// n+1 ETX
// n+2-3 CRC16
// / TODO: Replace it to your specific protocol struct
func unpackClientRspPacket(clientRspPkt []byte) (byte, []byte, error) {

	golog.CtsLog.Debug(" unpackClientRspPacket\n >> clientRspPkt[%X] 1 >>>>>> \n\n", clientRspPkt)
	resultCode, clientRspData, err := ParseClientRspPacket(clientRspPkt)
	if err != nil || resultCode != OK {
		return resultCode, nil, fmt.Errorf("resultCode[0x%X]ER err[%v]", int(resultCode), err)
	} else if len(clientRspData) != int(clientRspPkt[2]-1) {
		return resultCode, nil, fmt.Errorf("resultCode[0x%X]ER len(clientRspPayload)[%d] <> int(clientRspPacket[2]-1=[%d]", int(resultCode), len(clientRspPkt), int(clientRspPkt[2]-1))
	}
	golog.CtsLog.Debug(" unpackClientRspPacket\n >> clientRspPkt[%X] 2 >>>>>> \n\n", clientRspPkt)

	// Get ClientRequestCmd & ClientRequestPayloadLen
	clientRspCmd := clientRspPkt[1]
	clientRspPayloadLen := clientRspPkt[2]

	golog.CtsLog.Debug(" unpackClientRspPacket\n >> clientRspPkt[%X] 3 >>>>>> \n\n", clientRspPkt)
	// Create Client Response Packet
	switch clientRspCmd {
	case CMD_NOP: // NOP
		golog.CtsLog.Warn(" unpackClientRspPacket\n >> clientRspPkt[%X] clientRequestCmd[0x%02X]NOP   clientRspPayloadLen[%d]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen)

	case CMD_RD_DI: // Rd DI   ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RO
		var wDigitalInAddr uint16 = 0
		wDigitalInAddr = uint16(clientRspData[0]) * 256
		wDigitalInAddr += uint16(clientRspData[1])
		if int(wDigitalInAddr) > len(DigitalInputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n >> clientRspPkt[%X] clientRspCmd[0x%02X]Rd DO clientRequestPayloadLen[%d] wAnalogOutAddr[0x%04X]  wDigitalInAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wDigitalInAddr, resultCode)
			break
		}
		clientRspData = append(clientRspData, byte(wDigitalInAddr>>8))   // ADDRH
		clientRspData = append(clientRspData, byte(wDigitalInAddr&0xFF)) // ADDRH
		// Reading Ditigal Input incremented
		DigitalInputs[int(wDigitalInAddr)]++
		var byDIValue = DigitalInputs[int(wDigitalInAddr)]
		clientRspData = append(clientRspData, byDIValue) // Value = 0x01(ON)
		golog.CtsLog.Warn(" unpackClientRspPacket\n >> clientRspPkt[%X] clientRspCmd[0x%02X]Rd DI clientRspPayloadLen[%d] wDigitalInAddr[0x%04X] value[0x%02X]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wDigitalInAddr, byDIValue)

	case CMD_RD_DO: // Rd DO ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RW
		var wDigitalOutAddr uint16 = 0
		wDigitalOutAddr = uint16(clientRspData[0]) * 256
		wDigitalOutAddr += uint16(clientRspData[1])
		if int(wDigitalOutAddr) > len(DigitalOutputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n >> clientRspPkt[%X] clientRspCmd[0x%02X]Rd DO clientRspPayloadLen[%d] wAnalogOutAddr[0x%04X]  wDigitalOutAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wDigitalOutAddr, resultCode)
			break
		}
		clientRspData = append(clientRspData, byte(wDigitalOutAddr>>8))   // ADDRH
		clientRspData = append(clientRspData, byte(wDigitalOutAddr&0xFF)) // ADDRH
		// Reading Ditigal Output
		var byDOValue = DigitalOutputs[int(wDigitalOutAddr)]
		clientRspData = append(clientRspData, byDOValue)
		golog.CtsLog.Warn(" unpackClientRspPacket\n << clientRspPkt[%X] clientRequestCmd[0x%02X]Rd DO clientRequestPayloadLen[%d] wDigitalOutAddr[0x%04X] byDOValue[0x%02X] clientRspData[%x]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wDigitalOutAddr, byDOValue, clientRspData)

	case CMD_RD_AI: // Rd AI ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RO
		var wAnalogInAddr uint16 = 0
		wAnalogInAddr = uint16(clientRspData[0]) * 256
		wAnalogInAddr += uint16(clientRspData[1])
		if int(wAnalogInAddr) > len(AnalogInputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Rd AI clientRspPayloadLen[%d] wAnalogInAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogInAddr, resultCode)
			break
		}
		clientRspData = append(clientRspData, byte(wAnalogInAddr>>8))   // ADDRH
		clientRspData = append(clientRspData, byte(wAnalogInAddr&0xFF)) // ADDRH
		// Reading Inccremented Analog Input Data
		AnalogInputs[int(wAnalogInAddr)]++
		var wAIValue = AnalogInputs[int(wAnalogInAddr)]
		clientRspData = append(clientRspData, byte(wAIValue>>8))
		clientRspData = append(clientRspData, byte(wAIValue&0xFF))
		golog.CtsLog.Warn(" unpackClientRspPacket\n << clientRspPkt[%X] clientRequestCmd[0x%02X]Rd AI clientRspPayloadLen[%d] wAnalogInAddr[0x%04X] value[0x%04X]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogInAddr, wAIValue)

	case CMD_RD_AO: // Rd AO ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RO
		var wAnalogOutAddr uint16 = 0
		wAnalogOutAddr = uint16(clientRspData[0]) * 256
		wAnalogOutAddr += uint16(clientRspData[1])
		if int(wAnalogOutAddr) > len(AnalogOutputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Rd AO clientRsptPayloadLen[%d] wAnalogInAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogOutAddr, resultCode)
			break
		}
		clientRspData = append(clientRspData, byte(wAnalogOutAddr>>8))   // ADDRH
		clientRspData = append(clientRspData, byte(wAnalogOutAddr&0xFF)) // ADDRH
		// Reading AnalogInputValue
		var wAOValue = AnalogOutputs[int(wAnalogOutAddr)]
		clientRspData = append(clientRspData, byte(wAOValue>>8))
		clientRspData = append(clientRspData, byte(wAOValue&0xFF))
		golog.CtsLog.Warn(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Rd AO clientRspPayloadLen[%d] wAnalogInAddr[0x%04X] wAOValue[0x%04X]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogOutAddr, wAOValue)

	case CMD_WR_DO: // Wr DO ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RO
		var byaDigitalOutAddr uint16 = 0
		byaDigitalOutAddr = uint16(clientRspData[0]) * 256 // ADDRH
		byaDigitalOutAddr += uint16(clientRspData[1])      // ADDRL
		if int(byaDigitalOutAddr) > len(DigitalOutputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Wr DO clientRspPayloadLen[%d] wAnalogOutAddr[0x%04X]  byaDigitalOutAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, byaDigitalOutAddr, resultCode)
			break
		}
		// Save AO wValue on AnalogOutputs
		var byValue byte = 0
		byValue = uint8(clientRspData[2]) //Value
		DigitalOutputs[byaDigitalOutAddr] = byValue
		golog.CtsLog.Warn(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Wr D0 clientRspPayloadLen[%d] memBitAdd[0x%04X] value[0x%02X]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, byaDigitalOutAddr, byValue)

	case 6: // Wr AO ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RW
		var wAnalogOutAddr uint16 = 0
		wAnalogOutAddr = uint16(clientRspData[0]) * 256 // ADDRH
		wAnalogOutAddr += uint16(clientRspData[1])      //ADDRL
		if int(wAnalogOutAddr) > len(AnalogOutputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]Wr AO clientRspPayloadLen[%d] wAnalogOutAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogOutAddr, resultCode)
			break
		}
		// Save AO wValue on AnalogOutputs
		var wValue uint16 = 0
		wValue = uint16(clientRspData[2]) * 256 // ValueH
		wValue += uint16(clientRspData[3])      // ValueL
		AnalogOutputs[wAnalogOutAddr] = wValue  // Save it value
		golog.CtsLog.Warn(" unpackClientRspPacket\n << clientRspPacket[%X] clientRspCmd[0x%02X]Wr AO clientRspPayloadLen[%d] memBitAdd[0x%04X] value[0x%04X]\n", clientRspPkt, clientRspCmd, clientRspPayloadLen, wAnalogOutAddr, wValue)

	default:
		resultCode = ERR_INVALID_CMD // Invalid Cmd Received
		golog.CtsLog.Error(" unpackClientRspPacket\n << clientRspPkt[%X] clientRspCmd[0x%02X]UNWN resultCode[0x%02X]Invalid Cmd clientRspData[%X]\n", clientRspPkt, clientRspCmd, resultCode, clientRspData)
	}
	golog.CtsLog.Debug(" unpackClientRspPacket\n >> clientRspPacket[%X] 4 >>>>>> \n\n", clientRspPkt)
	golog.CtsLog.Debug(" unpackClientRspPacket\n >> clientRspPacket[%X] clientRspData[%X]\n", clientRspPkt, clientRspData)
	return resultCode, clientRspData, nil
}

// ClientRequestPacket:Header
// ClientRequestPacket+0 = ClientRequestPacletHeader+0 = STX
// ClientRequestPacket+1 = ClientRequestPacletHeader+1 = Cmd
// ClientRequestPacket+2 = ClientRequestPacletHeader+2 = PayLoadLen(n)
// ClientRequestPacket:Payload
// ClientRequestPacket+3 = ClientRequestPayload+0 > PayLoadLen(n)
// ClientRequestPacket:Trealler
// ClientRequestPacket+n+3+0 = ClientPayLoadLen+n+1:ClientRequestTrealler+0 = ETX
// ClientRequestPacket+n+3+1 = ClientPayLoadLen+n+2:ClientRequestTrealler+1 = CRC16H
// ClientRequestPacket+n+3+2 = ClientPayLoadLen+n+3:ClientRequestTrealler+2 = CRC16L
// / TODO: Replace it to your specific protocol struct
func ParseClientRquPacket(clientRquPacket []byte) (byte, []byte, error) {
	if clientRquPacket == nil {
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("clientRquPacket = nil")
	} else if len(clientRquPacket) < 6 { // STX(1) + Cmd(1) + PayLoadLen(1) + ETX(1) + CRC16(2)
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("len(clientRquPacket)=%d expected >=6", len(clientRquPacket))
	} else if len(clientRquPacket) != int(clientRquPacket[2])+6 { // STX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("len(clientRquPacket)=%d expected clientRquPacket[2]=[0x%02X]+6", len(clientRquPacket), clientRquPacket[2])
	} else if clientRquPacket[0] != CSTX { // STX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("clientRquPacket[0]=[0x%02X] expected [02]STX", clientRquPacket[0])
	} else if clientRquPacket[int(clientRquPacket[2])+3] != CETX { // ETX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("clientRquPacket[int(clientRquPacket[2])+3]0x%02X expected [03]ETX", clientRquPacket[int(clientRquPacket[2])+3])
	}
	// Coomputing CRC16
	clientRquPacketWoCRC16 := clientRquPacket[:len(clientRquPacket)-2]
	var clientRquPacketCrc16Recvd uint16 = 0
	clientRquPacketCrc16Recvd = uint16(clientRquPacket[len(clientRquPacket)-2]) * 256
	clientRquPacketCrc16Recvd += uint16(clientRquPacket[len(clientRquPacket)-1])
	ClientRequPacketCrc16Calc := CalcCRC16(clientRquPacketWoCRC16)
	if ClientRequPacketCrc16Calc != clientRquPacketCrc16Recvd {
		return ERR_BAD_CRC16, nil, fmt.Errorf(" clientRquPacket[%X] crc16Calc[0x%04X] != crc16Recvd[0x%04X]Missmatch", clientRquPacket, ClientRequPacketCrc16Calc, clientRquPacketCrc16Recvd)
	}
	clientRquPacketPayload := clientRquPacket[3 : len(clientRquPacket)-3]
	golog.CtsLog.Debug("parseClientRquPacket\n << clientRquPacket[%X] clientRquPacketPayload[%0X] ClientRequPacketCrc16Calc[0x%04X]\n", clientRquPacket, clientRquPacketPayload, ClientRequPacketCrc16Calc)
	return OK, clientRquPacketPayload, nil
}

// ============================================================================
// << ClientRspPacket - On Successfull Exit returns with Client Response Packet
// ============================================================================
// 0  STX
// 1  Cmd
// 2  PayLoadLen
// 3  ResultCode
// 4-n PayloadData
// n+1 ETX
// n+2-3 CRC16
// / TODO: Replace it to your specific protocol struct
func ParseClientRspPacket(clientRspPacket []byte) (byte, []byte, error) {
	if clientRspPacket == nil {
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("ParseClientRspPacketclientRsoPacket = nil")
		// ResultDataLen = PayLoadLen-1  PayLoadData[0]=resultCode[0x00=SUCCESS]
	} else if len(clientRspPacket) < 7 { // STX(1) + Cmd(1) + PayLoadLen(1) +ResultCode(1) + ResultData(n) + ETX(1) + CRC16(2)
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("ParseClientRspPacketlen(clientRspPacket)=%d expected >=7", len(clientRspPacket))
	} else if len(clientRspPacket) != int(clientRspPacket[2])+6 { // STX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("ParseClientRspPacket len(clientRspPacket)=%d expected clientRspPacket[2]=[0x%02X]+6=[%d]", len(clientRspPacket), clientRspPacket[2], int(clientRspPacket[2])+6)
	} else if clientRspPacket[0] != CSTX { // STX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("ParseClientRspPacket clientRspuPacket[0]=[0x%02X] expected [02]STX", clientRspPacket[0])
	} else if clientRspPacket[int(clientRspPacket[2])+3] != CETX { // ETX
		return ERR_BAD_PACKET_FORMAT, nil, fmt.Errorf("ParseClientRspPacket clientRspPacket[int(clientRspPacket[2])+3]0x%02X expected [03]ETX", clientRspPacket[int(clientRspPacket[2])+3])
	}
	// Coomputing CRC16
	clientRspPacketWoCRC16 := clientRspPacket[:len(clientRspPacket)-2]
	var clientRspPacketCrc16Recvd uint16 = 0
	clientRspPacketCrc16Recvd = uint16(clientRspPacket[len(clientRspPacket)-2]) * 256
	clientRspPacketCrc16Recvd += uint16(clientRspPacket[len(clientRspPacket)-1])
	ClientRspPacketCrc16Calc := CalcCRC16(clientRspPacketWoCRC16)
	if ClientRspPacketCrc16Calc != clientRspPacketCrc16Recvd {
		return ERR_BAD_CRC16, nil, fmt.Errorf("ParseClientRspPacket clientRspPacket[%X] crc16Calc[0x%04X] != crc16Recvd[0x%04X]Missmatch", clientRspPacket, ClientRspPacketCrc16Calc, clientRspPacketCrc16Recvd)
	}
	resultCode := clientRspPacket[3]
	// 0   1   2           3          4-n         n+1 n+2-3
	// STX Cmd PayLoadLen  ResultCode PayLoadData ETX CRC16
	//         Payload     <-------------------->
	//         PayLoadData            <--------->
	clientRspData := clientRspPacket[4 : len(clientRspPacket)-3]
	golog.CtsLog.Debug("ParseClientRspPacket\n >> clientRspPacket[%X]\n         ResultCode[%02X]\n      ClientRspData[%X]\n", clientRspPacket, resultCode, clientRspData)
	return resultCode, clientRspData, nil
}

// ============================================================
// << ClientRquPacket - On Entry receives Client Request Packet
// ============================================================
// 0  STX
// 1  Cmd
// 2  PayLoadLen
// 3-n PayloadData
// n+1 ETX
// n+2-3 CRC16
// ============================================================================
// << ClientRspPacket - On Successfull Exit returns with Client Response Packet
// ============================================================================
// 0  STX
// 1  Cmd
// 2  PayLoadLen
// 3  ResultCode
// 4-n PayloadData
// n+1 ETX
// n+2-3 CRC16
// / TODO: Replace it to your specific protocol struct
func proccessClientRquPacket(clientRquPacket []byte) (byte, []byte, error) {

	var resultCode byte = OK      // Success
	var ResultPayLoadlen byte = 1 // Result Payload Len = 1 Only ResultCode
	var clientRspPacket []byte

	// Parse/Decode Client Request Packet and gat clientRequestPayload
	resultCode, clientRquPayload, err := ParseClientRquPacket(clientRquPacket)
	if err != nil || resultCode != OK {
		return resultCode, nil, fmt.Errorf("resultCode[0x%X]ER err[%v]", int(resultCode), err)
	} else if len(clientRquPayload) != int(clientRquPacket[2]) {
		return resultCode, nil, fmt.Errorf("resultCode[0x%X]ER len(clientRquPayload)[%d] <> int(clientRquPacket[2]=[%d]", int(resultCode), len(clientRquPayload), int(clientRquPacket[2]))
	}

	// Get ClientRequestCmd & ClientRequestPayloadLen
	clientRequestCmd := clientRquPacket[1]
	clientRequestPayloadLen := clientRquPacket[2]

	clientRspPacket = append(clientRspPacket, CSTX)             // STX
	clientRspPacket = append(clientRspPacket, clientRequestCmd) // Cmd
	clientRspPacket = append(clientRspPacket, ResultPayLoadlen) // Len
	clientRspPacket = append(clientRspPacket, resultCode)       // ResultCode

	// Create Client Response Packet
	switch clientRequestCmd {
	case CMD_NOP: // NOP
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]NOP   clientRequestPayloadLen[%d]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen)

	case CMD_RD_DI: // Rd DI   ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RO
		ResultPayLoadlen += 3
		var wDigitalInAddr uint16 = 0
		wDigitalInAddr = uint16(clientRquPayload[0]) * 256
		wDigitalInAddr += uint16(clientRquPayload[1])
		if int(wDigitalInAddr) > len(DigitalInputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd DO clientRequestPayloadLen[%d] wAnalogOutAddr[0x%04X]  wDigitalInAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wDigitalInAddr, resultCode)
			break
		}
		// Reading Ditigal Input incremented
		DigitalInputs[int(wDigitalInAddr)]++ // TODO: Increment Value
		var byDIValue = DigitalInputs[int(wDigitalInAddr)]
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byDIValue)           // Value = 0x01(ON)
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd DI clientRequestPayloadLen[%d] wDigitalInAddr[0x%04X] value[0x%02X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wDigitalInAddr, byDIValue)

	case CMD_RD_DO: // Rd DO ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RW
		var wDigitalOutAddr uint16 = 0
		wDigitalOutAddr = uint16(clientRquPayload[0]) * 256
		wDigitalOutAddr += uint16(clientRquPayload[1])
		if int(wDigitalOutAddr) > len(DigitalOutputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd DO clientRequestPayloadLen[%d] wAnalogOutAddr[0x%04X]  wDigitalOutAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wDigitalOutAddr, resultCode)
			break
		}
		// Reading Ditigal Output
		ResultPayLoadlen += 3
		DigitalOutputs[int(wDigitalOutAddr)]++ //TODO: Incremented alue
		var byDOValue = DigitalOutputs[int(wDigitalOutAddr)]
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byDOValue)           // Value = 0x01(ON)
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd DO clientRequestPayloadLen[%d] wDigitalOutAddr[0x%04X] byDOValue[0x%02X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wDigitalOutAddr, byDOValue)

	case CMD_RD_AI: // Rd AI ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RO
		var wAnalogInAddr uint16 = 0
		wAnalogInAddr = uint16(clientRquPayload[0]) * 256
		wAnalogInAddr += uint16(clientRquPayload[1])
		if int(wAnalogInAddr) > len(AnalogInputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd AI clientRequestPayloadLen[%d] wAnalogInAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogInAddr, resultCode)
			break
		}
		// Reading Inccremented Analog Input Data
		ResultPayLoadlen += 4
		AnalogInputs[int(wAnalogInAddr)]++ //TODO: Incremented alue
		var wAIValue = AnalogInputs[int(wAnalogInAddr)]
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byte(wAIValue>>8))   // ValueH
		clientRspPacket = append(clientRspPacket, byte(wAIValue&0xFF)) // ValueL
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd AI clientRequestPayloadLen[%d] wAnalogInAddr[0x%04X] value[0x%04X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogInAddr, wAIValue)

	case CMD_RD_AO: // Rd AO ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RO
		var wAnalogInAddr uint16 = 0
		wAnalogInAddr = uint16(clientRquPayload[0]) * 256
		wAnalogInAddr += uint16(clientRquPayload[1])
		if int(wAnalogInAddr) > len(AnalogOutputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd AO clientRequestPayloadLen[%d] wAnalogInAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogInAddr, resultCode)
			break
		}
		// Reading AnalogInputValue
		ResultPayLoadlen += 4
		AnalogOutputs[int(wAnalogInAddr)]++ // TODO:Increment value
		var wAOValue = AnalogOutputs[int(wAnalogInAddr)]
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byte(wAOValue>>8))   // ValueH
		clientRspPacket = append(clientRspPacket, byte(wAOValue&0xFF)) // ValueL
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Rd AO clientRequestPayloadLen[%d] wAnalogInAddr[0x%04X] wAOValue[0x%04X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogInAddr, wAOValue)

	case CMD_WR_DO: // Wr DO ClientRspPayload = ResultCode(1) ADDR(2) Value(1) RO
		var byaDigitalOutAddr uint16 = 0
		byaDigitalOutAddr = uint16(clientRquPayload[0]) * 256 // ADDRH
		byaDigitalOutAddr += uint16(clientRquPayload[1])      // ADDRL
		if int(byaDigitalOutAddr) > len(DigitalOutputs) {
			resultCode = ERR_DIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Wr DO clientRequestPayloadLen[%d] wAnalogOutAddr[0x%04X]  byaDigitalOutAddr[0x%02X]ERR_DIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, byaDigitalOutAddr, resultCode)
			break
		}
		// Save DO byaDigitalOutAddr
		ResultPayLoadlen += 3
		var byValue byte = 0
		byValue = uint8(clientRquPayload[2])                           //Value
		DigitalOutputs[int(byaDigitalOutAddr)] = byValue               // TODO:Save Rcvd  value
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byValue)             // Value = 0x01(ON)
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Wr D0 clientRequestPayloadLen[%d] memBitAdd[0x%04X] value[0x%02X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, byaDigitalOutAddr, byValue)

	case 6: // Wr AO ClientRspPayload = ResultCode(1) ADDR(2) Value(2) RW
		var wAnalogOutAddr uint16 = 0
		wAnalogOutAddr = uint16(clientRquPayload[0]) * 256 // ADDRH
		wAnalogOutAddr += uint16(clientRquPayload[1])      //ADDRL
		if int(wAnalogOutAddr) > len(AnalogOutputs) {
			resultCode = ERR_AIO_ADDR_ABOVE_MAX_LIMIT // Above Max Limit
			golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Wr AO clientRequestPayloadLen[%d] wAnalogOutAddr[0x%04X]  resultCode[0x%02X]ERR_AIO_ADDR_ABOVE_MAX_LIMIT\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogOutAddr, resultCode)
			break
		}
		// Save AO wValue on AnalogOutputs
		ResultPayLoadlen += 4
		var wValue uint16 = 0
		wValue = uint16(clientRquPayload[2]) * 256                     // ValueH
		wValue += uint16(clientRquPayload[3])                          // ValueL
		AnalogOutputs[int(wAnalogOutAddr)] = wValue                    // Save it value
		clientRspPacket = append(clientRspPacket, clientRquPayload[0]) // ADDRH
		clientRspPacket = append(clientRspPacket, clientRquPayload[1]) // ADDRL
		clientRspPacket = append(clientRspPacket, byte(wValue>>8))     // ValueH
		clientRspPacket = append(clientRspPacket, byte(wValue&0xFF))   // ValueL
		golog.CtsLog.Warn(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]Wr AO clientRequestPayloadLen[%d] memBitAdd[0x%04X] value[0x%04X]\n", clientRquPacket, clientRequestCmd, clientRequestPayloadLen, wAnalogOutAddr, wValue)

	default:
		resultCode = ERR_INVALID_CMD // Invalid Cmd Received
		golog.CtsLog.Error(" parseClientRquPacket\n << clientRquPacket[%X] clientRequestCmd[0x%02X]UNWN resultCode[0x%02X]Invalid Cmd clientRequestPayloadLen[%d]\n", clientRquPacket, clientRequestCmd, resultCode, clientRequestPayloadLen)
	}

	clientRspPacket[2] = ResultPayLoadlen           // Update Client RspPacket Len
	clientRspPacket[3] = resultCode                 // Update Client RspPacket ResultCode
	clientRspPacket = append(clientRspPacket, CETX) // ETX
	clientRspPacketCrc16Calc := CalcCRC16(clientRspPacket)
	clientRspPacket = append(clientRspPacket, byte(clientRspPacketCrc16Calc/256))  // CRC16H
	clientRspPacket = append(clientRspPacket, byte(clientRspPacketCrc16Calc&0xff)) // CRC16L
	golog.CtsLog.Debug("proccessClientRspPacket\n >> clientRspPacket[%X] clientRspPacketCrc16Calc[0x%04X]\n", clientRspPacket, clientRspPacketCrc16Calc)
	return resultCode, clientRspPacket, nil
}
