package canopen

import (
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"fmt"
	"net"
	"time"
)

// / ==============================================================================
// / CANopenTcpClient_HandleServerConnection - CANopen TCP(CIA309): Server (Device)
// / ==============================================================================
func CANopenTcpClient_HandleServerConnection(serverConn net.Conn) {
	if serverConn == nil {
		logPkg.CtsLog.Error("CANopenTcpClient_HandleServerConnection: canOpenClientConn=nil")
		return
	}
	defer serverConn.Close()

	var resultCode byte = utilsPkg.OK
	// Get the server address and port
	remoteAddr := serverConn.RemoteAddr().String()
	// Get the server Network
	remoteNetwork := serverConn.RemoteAddr().Network()

	logPkg.CtsLog.Warn("CANopenTcpClient_HandleServerConnection: remoteNetwork[%s] remoteAddr[%s] Connected ", remoteNetwork, remoteAddr)

	buffer := make([]byte, 65536) // CANopen Max Buffer Size allowed is 64K
	startTime := time.Now()

	for {
		// Waiting Packet from the CANopen TCP Server(CIA309)
		n, err := serverConn.Read(buffer)
		if n == 0 {
			// elapsedTime := time.Since(startTime)
			if time.Since(startTime) >= 10*time.Second {
				// Timeout Occurred w/o read nothing for 10 seconds after connected
				logPkg.CtsLog.Debug("CANopenTcpClient_HandleServerConnection:No CANopen TCP Event Msg n[%d] 4 err[%s]", n, err)
			}
		} else {
			// TODO: roccess CANopen TCP Server Packet may be an event
			rcvdCANopenTcpServerEvtPacket := buffer[:n]
			logPkg.CtsLog.Warn("handleCANopencanOpenClientConnection: Server Evt\n << rcvdCANopenTcpServerEvtPacket[%d:%X] from remoteAddr[%s]", len(rcvdCANopenTcpServerEvtPacket), rcvdCANopenTcpServerEvtPacket, remoteAddr)
		}

		startTime = time.Now()

		// TODO: Proccess a Pending CANopen TCP Client Request MSG
		/// Store Rcvd CANoopen TCP Client Packet into a buffer with it size
		/// Proccess Rcvd CANoopen TCP Client Packet and create a CANopen TCP Client rsp if required
		/// resultCode, sentClientRspPacket, err := CANopenTcpServer_ProcClientRquPkt(rcvdCANopenTcpClientRquPacket)
		/// if err != nil {
		///	logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER err[%s]", connId, resultCode, err)
		///	break
		///} else if resultCode != utilsPkg.OK {
		///	logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER ", connId, resultCode)
		///	break
		///}
		///if sentClientRspPacket != nil {
		///	// Send canOpenclientRspPacket back to the CANopen TCP Client
		///	_, err = canOpenClientConn.Write(sentClientRspPacket)
		///	if err != nil {
		///		resultCode = utilsPkg.ERR_COMM_WRITE
		///		logPkg.CtsLog.Error("handleCANopencanOpenClientConnection:connId[%d] resultCode[0x%X]ER canOpenClientConn.Write err[%s]", connId, resultCode, err)
		///		break
		///	}
		///	logPkg.CtsLog.Warn("handleCANopencanOpenClientConnection:connId[%d]Sent Client Rsp Packet\n >> sentClientRspPacket[%X] to remoteAddr[%s]\n\n", connId, sentClientRspPacket, remoteAddr)
		if time.Since(startTime) >= 10*time.Second {
			// Timeout Occurred w/o read nothing for 10 seconds after connected
			logPkg.CtsLog.Debug("CANopenTcpClient_HandleServerConnection:No CANopen TCP Event Msg n[%d] 4 err[%s]", n, err)
			break
		}
	}
	logPkg.CtsLog.Debug("CANopenTcpClient_HandleServerConnection:resultCode[0x%X] remoteAddr[%s] Exiting!", int(resultCode), remoteAddr)
}

func CANopenTcpClient_ConnectToServer(conn net.Conn, ipaddr string, port string) (net.Conn, error) {
	var err error
	tcpAddress := ipaddr + ":" + port
	if conn != nil {
		// CANopen over TCP device is already connected
		if !utilsPkg.AnyDeviceIsConnected {
			logPkg.CtsLog.Warn("CANopenTcpClient_ConnectToServer:WR tcpAddress[%s] conn[%v] Already Connected 1", tcpAddress, conn)
			utilsPkg.AnyDeviceIsConnected = true
		} else {
			logPkg.CtsLog.Warn("CANopenTcpClient_ConnectToServer:WR tcpAddress[%s] conn[%v] Already Connected 2", tcpAddress, conn)
		}
		return conn, nil
	}
	// Trying to connect to CANopen TCP Servers
	conn, err = net.DialTimeout("tcp", tcpAddress, 5*time.Second)
	if err != nil {
		// Fail to connect to CANopen TCP Server
		if conn != nil {
			conn.Close()
			conn = nil
		}
		if utilsPkg.AnyDeviceIsConnected {
			logPkg.CtsLog.Error("CANopenTcpClient_ConnectToServer:ER net.DialTimeout tcpAddress[%s] Disconnected", tcpAddress)
			utilsPkg.AnyDeviceIsConnected = false
		}
		return nil, err
	}
	// Connected to CANopen TCP Server
	if !utilsPkg.AnyDeviceIsConnected {
		utilsPkg.AnyDeviceIsConnected = true
		logPkg.CtsLog.Warn("CANopenTcpClient_ConnectToServer:OK tcpAddress[%s] Connected conn[%v]OK", tcpAddress, conn)
	}
	return conn, nil
}

func CANopenTcpClient_DisconnectFromServer(ipaddr string, port uint16) error {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_DisconnectFromServer Not Implemented")
	return fmt.Errorf("TODO:CANopenTcpClient_DisconnectFromServer Not Implemented")
}

func IsCANopenTcpServerConnected() error {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_DisconnectFromServer Not Implemented")
	return fmt.Errorf("TODO:CANopenTcpClient_DisconnectFromServer Not Implemented")
}

func CANopenTcpClient_ProcServerEvtPkt(canOpenServerEvtPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_ProcServerEvtPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_ProcServerEvtPkt Not Implemented")
}

func CANOpenTcpClient_GenCANopenClientRquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANOpenTcpClient_GenCANopenClientRquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANOpenTcpClient_GenCANopenClientRquPkt Not Implemented")
}

func CANopenTcpClient_Create_SDO_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
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
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_SDO_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_SDO_RquPkt Not Implemented")
}

func CANopenTcpClient_Create_PDO_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
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
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_PDO_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_PDO_RquPkt Not Implemented")
}

func CANopenTcpClient_Create_NMT_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_NMT_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_NMT_RquPkt Not Implemented")
}

func CANopenTcpClient_Create_EMCY_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_EMCY_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_EMCY_RquPkt Not Implemented")
}

func CANopenTcpClient_Create_SYNC_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_SYNC_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_SYNC_RquPkt Not Implemented")
}

func CANopenTcpClient_Create_TIME_RquPkt(canOpenClientRquPacket []byte) (byte, []byte, error) {
	logPkg.CtsLog.Error("TODO:CANopenTcpClient_Create_TIME_RquPkt Not Implemented")
	return 0x80, nil, fmt.Errorf("TODO:CANopenTcpClient_Create_TIME_RquPkts Not Implemented")
}

// Proccessing CANopen TCP Client SDO Msg
// STEP01: Build CANopen TCP Client SDO Rqu Msg
// STEP02: Send  CANopen TCP Client SDO Rqu Msg >> CANopen TCP Server
// STEP03: Recv  CANopen TCP Client SDO Rsp Msg << CANopen TCP Server
// STEP04: Parse CANopen TCP Client SDO Rsp Msg << Get Response Data and store it
func CANopenTcpClient_Proc_SDO_Download_Msg(conn net.Conn, NodeId uint8, rquData []byte) error {

	if conn == nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Download_Msg:ParseCANopenRawFrame Conn=nil Device Disconnected\n")
		return fmt.Errorf("CANopenTcpClient_Proc_SDO_Download_Msg:conn=nil Device Disconnected")
	}

	// TODO: Make CANopen SDO Client Rqu Msg
	rawCANopenframe := []byte{0x63, 0xA0, 0x23, 0x45, 0x08, 0x40, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	rawCANopenframe[0] = 0x60 + NodeId>>4
	rawCANopenframe[1] = NodeId << 4
	err := ParseRawCANopenFrame(rawCANopenframe)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Download_Msg: ParseCANopenRawFrame err[%s]\n", err)
		return err
	}

	// Sending CANopen SDO TCP Client Rqu Msg to CANopen TCP Server
	_, err = conn.Write(rawCANopenframe)
	if err != nil {
		if conn != nil {
			conn.Close()
			conn = nil
		}
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Download_Msg: conn.Write err[%s]\n", err)
		return err
	}
	logPkg.CtsLog.Debug("CANopenTcpClient_Proc_SDO_Download_Msg:\n\t  >> Sent: rawCANopenframe[%X]\n", rawCANopenframe)

	// Receiving CANopen SDO TCP Client Rsp Msg From CANopen TCP Server
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		if conn != nil {
			conn.Close()
			conn = nil
			logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Download_Msg: Error receiving data:", err)
		}
		return err
	}

	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	logPkg.CtsLog.Debug("CANopenTcpClient_Proc_SDO_Download_Msg:\n\t  << Rcvd:sClientRspPacket[%X]\n", ClientRspPacket)

	// TODO Unpacking Received packet too get response data
	err = ParseRawCANopenFrame(ClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Download_Msg: ParseCANopenRawFrame err[%s]\n", err)
		return err
	}
	return nil
}

// Proccessing CANopen TCP Client SDO Msg
// STEP01: Build CANopen TCP Client SDO Rqu Msg
// STEP02: Send  CANopen TCP Client SDO Rqu Msg >> CANopen TCP Server
// STEP03: Recv  CANopen TCP Client SDO Rsp Msg << CANopen TCP Server
// STEP04: Parse CANopen TCP Client SDO Rsp Msg << Get Response Data and store it
func CANopenTcpClient_Proc_SDO_Upload_Msg(conn net.Conn, NodeId uint8, rawRquData []byte) ([]byte, error) {
	if conn == nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Upload_Msg:ParseCANopenRawFrame Conn=nil Device Disconnected\n")
		return nil, fmt.Errorf("CANopenTcpClient_Proc_SDO_Upload_Msg:conn=nil Device Disconnected")
	}

	//TODO: Make CANopen SDO Client Upload Rsp Data
	rawRspData := []byte{0x63, 0xA0, 0x23, 0x45, 0x08, 0x40, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	// TODO: Make CANopen SDO Client Rqu Msg
	rawCANopenframe := []byte{0x63, 0xA0, 0x23, 0x45, 0x08, 0x40, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	rawCANopenframe[0] = 0x60 + NodeId>>4
	rawCANopenframe[1] = NodeId << 4
	err := ParseRawCANopenFrame(rawCANopenframe)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Upload_Msg: ParseCANopenRawFrame err[%s]\n", err)
		return nil, err
	}

	// Sending CANopen SDO TCP Client Rqu Msg to CANopen TCP Server
	_, err = conn.Write(rawCANopenframe)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Upload_Msg: conn.Write err[%s]\n", err)
		return nil, err
	}
	logPkg.CtsLog.Debug("CANopenTcpClient_Proc_SDO_Upload_Msg:\n\t  >> Sent: rawCANopenframe[%X]\n", rawCANopenframe)

	// Receiving CANopen SDO TCP Client Rsp Msg From CANopen TCP Server
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Upload_Msg: Error receiving data:", err)
		return nil, err
	}

	// Trim the response to the actual data received
	ClientRspPacket = ClientRspPacket[:n]
	logPkg.CtsLog.Warn("CANopenTcpClient_Proc_SDO_Upload_Msg:\n\t  << Rcvd:sClientRspPacket[%X]\n", ClientRspPacket)

	// TODO Unpacking Received packet too get response data
	err = ParseRawCANopenFrame(ClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_SDO_Upload_Msg: ParseCANopenRawFrame err[%s]\n", err)
		return nil, err
	}
	return rawRspData, nil
}

// Proccessing CANopen TCP Client PDO Msg
// STEP01: Build CANopen TCP Client PDO Rqu Msg
// STEP02: Send  CANopen TCP Client PDO Rqu Msg >> CANopen TCP Server
// STEP03: Recv  CANopen TCP Client PDO Rsp Msg << CANopen TCP Server
// STEP04: Parse CANopen TCP Client PDO Rsp Msg << Get Response Data and store it
func CANopenTcpClient_Proc_PDO_Msg(conn net.Conn, pdotype byte, NodeId uint8, index uint16, subindex uint8, rquData []byte) ([]byte, error) {
	if conn == nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_PDO_Msg:ParseCANopenRawFrame Conn=nil Device Disconnected\n")
		return nil, fmt.Errorf("CANopenTcpClient_Proc_PDO_Msg:conn=nil Device Disconnected")
	}
	// TODO: Make CANopen PDO Client Rqu Msg
	rawRspData := []byte{0x4B, 0xA1, 0x23, 0x45, 0x08, 0x01, 0x18, 0x03, 0x01, 0x05, 0x06, 0x07, 0x08}

	// TODO: Make CANopen PDO Client Rqu Msg
	rawCANopenframe := []byte{0x4B, 0xA1, 0x23, 0x45, 0x08, 0x01, 0x18, 0x03, 0x01, 0x05, 0x06, 0x07, 0x08}

	//rawCANopenframe[0] = 0x18 + pdotype*8<<7 + NodeId>>4 //NodeID id.0-2  11 bits
	//rawCANopenframe[1] = NodeId << 4                     //NodeID id.3-10 11 bits
	err := ParseRawCANopenFrame(rawCANopenframe)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_PDO_Msg:ParseCANopenRawFrame err[%s]\n", err)
		return nil, err
	}

	// Sending CANopen PDO TCP Client Rqu Msg to CANopen TCP Server
	_, err = conn.Write(rawCANopenframe)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_PDO_Msg:conn.Write err[%s]\n", err)
		return nil, err
	}
	logPkg.CtsLog.Debug("CANopenTcpClient_Proc_PDO_Msg:\n\t >> Sent: rawCANopenframe[%X]\n", rawCANopenframe)

	// Receiving CANopen PDO TCP Client Rsp Msg from CANopen TCP Server
	ClientRspPacket := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := conn.Read(ClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_PDO_Msg: Error receiving data:", err)
		return nil, err
	}

	// Trim the response to the actual data receiveds
	ClientRspPacket = ClientRspPacket[:n]
	logPkg.CtsLog.Debug("CANopenTcpClient_Proc_PDO_Msg:\n\t << Rcvd:ClientRspPacket[%X]\n", ClientRspPacket)

	// TODO Unpacking Received packet too get response data
	err = ParseRawCANopenFrame(ClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("CANopenTcpClient_Proc_PDO_Msg: ParseCANopenRawFrame err[%s]\n", err)
		return nil, err
	}
	return rawRspData, err
}
