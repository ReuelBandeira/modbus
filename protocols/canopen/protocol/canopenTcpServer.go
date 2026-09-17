package canopen

import (
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"net"
	"time"
)

// / ==============================================================================
// / CANopenTcpServer_HandleClientConnection - CANopen TCP(CIA309): Server (Device)
// / ==============================================================================
func CANopenTcpServer_HandleClientConnection(connId int, canOpenClientConn net.Conn) {
	if canOpenClientConn == nil {
		logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] canOpenClientConn=nil", connId)
		return
	}
	defer canOpenClientConn.Close()

	var resultCode byte = utilsPkg.OK
	// Get the client address and port
	remoteAddr := canOpenClientConn.RemoteAddr().String()
	// Get the client Network
	remoteNetwork := canOpenClientConn.RemoteAddr().Network()

	logPkg.CtsLog.Warn("handleCANopencanOpenClientConnection:connId[%d] remoteNetwork[%s] remoteAddr[%s] Connected ", connId, remoteNetwork, remoteAddr)

	buffer := make([]byte, 65536) // CANopen Max Buffer Size allowed is 64K
	startTime := time.Now()

	for {
		// Read CANopen Rqu Packet from the CANopen TCP Client(CIA309)
		n, err := canOpenClientConn.Read(buffer)
		if n == 0 {
			// elapsedTime := time.Since(startTime)
			if time.Since(startTime) >= 10*time.Second {
				// Timeout Occurred w/o read nothing for 10 seconds after connected
				logPkg.CtsLog.Debug("handleCANopencanOpenClientConnection:connId[%d]No Data for more than 10 Seconds n[%d] 4 err[%s]", connId, n, err)
				break
			}
			continue

		}
		/// Last Successfull read from this connection
		startTime = time.Now()

		/// Store Rcvd CANoopen TCP Client Packet into a buffer with it size
		rcvdCANopenTcpClientRquPacket := buffer[:n]
		logPkg.CtsLog.Warn("handleCANopencanOpenClientConnection:connId[%d]Rcvd Client Rqu\n << rcvdCANopenTcpClientRquPacket[%d:%X] from remoteAddr[%s]", connId, len(rcvdCANopenTcpClientRquPacket), rcvdCANopenTcpClientRquPacket, remoteAddr)

		/// Proccess Rcvd CANoopen TCP Client Packet and create a CANopen TCP Client rsp if required
		resultCode, sentClientRspPacket, err := CANopenTcpServer_ProcClientRquPkt(rcvdCANopenTcpClientRquPacket)
		if err != nil {
			logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER err[%s]", connId, resultCode, err)
			break
		} else if resultCode != utilsPkg.OK {
			logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER ", connId, resultCode)
			break
		}
		if sentClientRspPacket != nil {
			// Send canOpenclientRspPacket back to the CANopen TCP Client
			_, err = canOpenClientConn.Write(sentClientRspPacket)
			if err != nil {
				resultCode = utilsPkg.ERR_COMM_WRITE
				logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER canOpenClientConn.Write err[%s]", connId, resultCode, err)
				break
			}
			logPkg.CtsLog.Warn("handleCANopencanOpenClientConnection:connId[%d]Sent Client Rsp Packet\n >> sentClientRspPacket[%X] to remoteAddr[%s]\n\n", connId, sentClientRspPacket, remoteAddr)
		}
	} // End of Continuous loop
	logPkg.CtsLog.Debug("handleCANopencanOpenClientConnection:connId[%d]resultCode[0x%X] remoteAddr[%s] Exiting!", connId, int(resultCode), remoteAddr)
}

// / ==================================================================================
// / Server: << CANopenTcpServer_ProcClientRquPkt - On Entry receives Client Request Packet
// / Server: >> On Exit return CanOpenClientRspPacket if is required and result code
// / ==================================================================================
func CANopenTcpServer_ProcClientRquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Debug("CANopenTcpServer_ProcClientRquPkt: << canOpenClientRquPacket[%X]", canOpenClientRquPacket)
	ucResultCode, canOpenClientRspPacket, err := CANOpenTcpServer_GenCANopenClientRspPkt(canOpenClientRquPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpServer_ProcClientRquPkt: err[%s]", err)
	} else if ucResultCode != 0x00 {
		logPkg.CtsLog.Error("CANopenTcpServer_ProcClientRquPkt: ucResultCode[0x%02X]", ucResultCode)
	} else {
		logPkg.CtsLog.Debug("CANopenTcpServer_ProcClientRquPkt: >> canOpenClientRspPacket[%X]", canOpenClientRspPacket)
	}
	return ucResultCode, canOpenClientRspPacket, err
}

func CANOpenTcpServer_GenCANopenClientRspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00
	err := ParseRawCANopenFrame(canOpenClientRquPacket) // canopenControl.go
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	//TODO: Generate specific Client Response Packet accordign with it Requests
	return ucResultCode, canOpenClientRquPacket, nil
}

func CANopenTcpServer_Create_SDO_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen SDO response packet
	var canOpenClientRspPacket []byte
	/// ===========================
	/// controlByte 01xx xxxx = SDO
	/// ===========================
	/// COB-ID: 0x580 + node ID = transmitting from Node
	/// COB-ID: 0x600 + node ID = receiving    to   Node
	/// To initiate a download, the SDO client sends the following data
	/// in a CAN message with the 'receive' COB-ID of the SDO channel.
	///       Byte 0(Nr)       Byte 1            Byte 2            Byte 3          B4-8
	///  7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	/// |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |--------|
	/// | cs  |x| n |e|s| |               idx               | |  sub-Index    | |<d=data>|
	///   cs = Client Specifier Command for SDO
	///        010 -SDO Request (client command specifier) with expedited transfer and size indicated
	///    x = Reserved
	///    n = number of data bytes of msg w/o contain data. valid if e & s = 1
	///    e = indicates an expedited transfer
	///    s = indicates data size is specified in n (if e is set) or data msg
	///  idx = index
	///  sub = sub-index
	/// TODO: Implements CANopen SDO Rsp packet
	err := ParseRawCANopenFrame(canOpenClientRquPacket) // canopenControl.go
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          canFrame.DLC,
		Data:         canFrame.Data,
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt FunctionCode[%d] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt       NodeID[0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt          ext[%v] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt          rtr[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt          ext[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt          DLC[0x%X] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt         Data[%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_SDO_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	canOpenClientRspPacket[0] = 0x58 + +byte(nodeid)&0xFF>>4
	canOpenClientRspPacket[1] = +byte(nodeid) & 0xFF << 4
	return ucResultCode, canOpenClientRspPacket, nil
}

func CANopenTcpServer_Create_PDO_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen PDO response packet
	var canOpenClientRspPacket []byte
	/// ===========================
	/// controlByte 01xx xxxx = PDO
	/// ===========================
	/// COB-ID: 0x180 + node ID = PDO1.Rx
	/// COB-ID: 0x200 + node ID = PDO1.Tx
	/// COB-ID: 0x280 + node ID = PDO2.Rx
	/// COB-ID: 0x300 + node ID = PDO2.Tx
	/// COB-ID: 0x380 + node ID = PDO3.Rx
	/// COB-ID: 0x400 + node ID = PDO3.Tx
	/// COB-ID: 0x480 + node ID = PDO4.Rx
	/// COB-ID: 0x500 + node ID = PDO4.Tx
	/// To initiate a download, the SDO client sends the following data
	/// in a CAN message with the 'receive' COB-ID of the SDO channel.
	///       Byte 0(Nr)       Byte 1            Byte 2            Byte 3          B4-8
	///  7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0   7 6 5 4 3 2 1 0
	/// |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |-|-|-|-|-|-|-|-| |--------|
	/// | cs  |x| n |e|s| |               idx               | |  sub-Index    | |<d=data>|
	///   cs = Client Specifier Command for PDO ???
	///        010 -PDO Request ??? (client command specifier) with expedited transfer and size indicated
	///    x = Reserved
	///    n = number of data bytes of msg w/o contain data. valid if e & s = 1
	///    e = indicates an expedited transfer
	///    s = indicates data size is specified in n (if e is set) or data msg
	///  idx = index
	///  sub = sub-index
	err := ParseRawCANopenFrame(canOpenClientRquPacket)
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          dlc,
		Data:         canOpenClientRquPacket[5:],
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt FunctionCode: [0x%X] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt       NodeID: [0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt     Extended: [0x%X] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt          RTR: [%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt          DLC: [%v] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt          DLC: [0x%X] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt         Data: [%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_PDO_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	canOpenClientRspPacket[0] = 0x18 + +byte(nodeid)&0xFF>>4
	canOpenClientRspPacket[1] = +byte(nodeid) & 0xFF << 4
	return ucResultCode, canOpenClientRspPacket, nil
}

func CANopenTcpServer_Create_NMT_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen PDO response packet
	var canOpenClientRspPacket []byte
	err := ParseRawCANopenFrame(canOpenClientRquPacket)
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          canFrame.DLC,
		Data:         canFrame.Data,
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt FunctionCode[%d] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt       NodeID[0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          ext[%v] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          RTR[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          DLC[%d] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         Data[%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	return ucResultCode, canOpenClientRspPacket, nil
}

func CANopenTcpServer_Create_EMCY_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen PDO response packet
	var canOpenClientRspPacket []byte
	err := ParseRawCANopenFrame(canOpenClientRquPacket)
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          canFrame.DLC,
		Data:         canFrame.Data,
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt FunctionCode[%d] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt       NodeID[0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          ext[%v] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          RTR[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          DLC[%d] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         Data[%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	return ucResultCode, canOpenClientRspPacket, nil
}

func CANopenTcpServer_Create_SYNC_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen PDO response packet
	var canOpenClientRspPacket []byte
	err := ParseRawCANopenFrame(canOpenClientRquPacket)
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          canFrame.DLC,
		Data:         canFrame.Data,
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt FunctionCode[%d] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt       NodeID[0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          ext[%v] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          RTR[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          DLC[%d] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         Data[%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	return ucResultCode, canOpenClientRspPacket, nil
}

func CANopenTcpServer_Create_TIME_RspPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode byte = 0x00 // TODO: CANopenTcpServer_Create CANopen PDO response packet
	var canOpenClientRspPacket []byte
	err := ParseRawCANopenFrame(canOpenClientRquPacket)
	if err != nil {
		///FAIL: Is not CANFrame
		ucResultCode = 0x81
		logPkg.CtsLog.Error("Error: err[%s]", err)
		return ucResultCode, nil, err
	}
	ext := (canOpenClientRquPacket[4] & 0x80) != 0x80
	rtr := (canOpenClientRquPacket[4] & 0x40) != 0x40
	dlc := canOpenClientRquPacket[4] & 0x0F

	canFrame := &CANOPEN_Frame{
		ID:   uint32(canOpenClientRquPacket[0]&0x0F)<<4 | uint32(canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:  dlc,
		Data: canOpenClientRquPacket[5:],
	}

	canOpenFrame := &CANopen_Frame{
		FunctionCode: canOpenClientRquPacket[0] >> 4,
		NodeID:       (canOpenClientRquPacket[0]&0x0F)<<4 | (canOpenClientRquPacket[1]&0xE0)>>4,
		DLC:          canFrame.DLC,
		Data:         canFrame.Data,
	}

	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt FunctionCode[%d] 4 bits - Byte  0 b4-7\n", canOpenFrame.FunctionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt       NodeID[0x%X] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", canOpenFrame.NodeID)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          ext[%v] 1 bits - Byte  4 b7\n", ext)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          RTR[%v] 1 bits - Byte  4 b6(Remote Transmision Request)\n", rtr)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt          DLC[%d] 4 bits - Byte  4 b0-3\n", canOpenFrame.DLC)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         Data[%X] 0-8 bytes- Bytes 5-7\n", canOpenFrame.Data)

	cobid := (canFrame.ID >> 20) & 0xF7F       // b0-10 - 11 bits BigEndian
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3  -  4 bits BigEndian
	nodeid := (canFrame.ID >> 20) & 0x7F       // b4-10 -  7 bits BigEndian
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         cobid [0x%X] 4 bits - Byte  0 b4-7\n", cobid)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt   functionCode[0x%d] 7 bits - Byte  1 b0-3 & Byte 2 b5-7 \n", functionCode)
	logPkg.CtsLog.Info("CANopenTcpServer_Create_NMT_RspPkt         nodeid[0x%X] 1 bits - Byte  4 b7\n", nodeid)
	canOpenClientRspPacket = canOpenClientRquPacket
	return ucResultCode, canOpenClientRspPacket, nil
}
