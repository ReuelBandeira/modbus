package main

import (
	utilsPkg "emuethernetip/utils"
	logPkg "emuethernetip/utils/gologtofile"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

var ushUpNewValue uint16 = 1 // This is the return value that will be used

// ===========================================
// Eumulating - CIP EthernetIP Server (Device)
// ===========================================
func handleClientConnection(protocolName string, cipEipClientConn net.Conn, connId int) {
	defer cipEipClientConn.Close()
	var err error = nil
	var r int = 0
	var w int = 0
	var NotifiedConnected bool = false
	var sentClientRspPacket []byte = nil
	var resultCode byte = utilsPkg.OK

	// Get the client's address and port
	remoteAddr := cipEipClientConn.RemoteAddr().String()
	// Get the client's address and port
	remoteNetwork := cipEipClientConn.RemoteAddr().Network()
	if connId == 4 || connId == 12 {
		logPkg.CtsLog.Warn("c:protocolName[%s] connId[%d] remoteNetwork[%s] remoteAddr[%s] Connected\r\n", protocolName, connId, remoteNetwork, remoteAddr)
	}

	buffer := make([]byte, 65536) // Adjust the buffer size as needed
	startTime := time.Now()

	// Set a read timeout of 30 seconds.
	timeout := 60 * time.Second
	err = cipEipClientConn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		cipEipClientConn.Close()
		logPkg.CtsLog.Error("c:protocolName[%s] connId[%d] remoteNetwork[%s] remoteAddr[%s] SetReadDeadline()fail\r\n", protocolName, connId, remoteNetwork, remoteAddr)
		return
	}

	for {
		// Read data from the CIP EthernetIP Client
		r, err = cipEipClientConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				resultCode = utilsPkg.ERR_COMM_READ_TIMEPUT
				if w != 0 {
					utilsPkg.TCPClientRxErrors++
					logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] Timeout err[%s]", protocolName, connId, err)
				}
				cipEipClientConn.Close()
				break
			}
			continue
		} else if r == 0 {
			// elapsedTime := time.Since(startTime)
			if time.Since(startTime) >= 10*time.Second {
				// Timeout Occurred w/o read nothing for 10 seconds after connected
				if r == 0 && w == 0 {
					utilsPkg.TCPClientRxErrors++
					logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d]Timeout No Data for more than 10Seconds r[%d] err[%s]", protocolName, connId, r, err)
				}
				cipEipClientConn.Close()
				break
			}
			if r == 0 && w == 0 {
				logPkg.CtsLog.Warn("handleClientConnection:protocolName[%s] connId[%d]No Data for more than 10Seconds r[%d] err[%s]", protocolName, connId, r, err)
			}
			continue
		} else if err == nil && r < 4 {
			utilsPkg.TCPClientRxErrors++
			err = fmt.Errorf("protocolname[%s] connId[%d] remoteAddr[%s] r[%d]< 4 Fail", protocolName, connId, remoteAddr, r)
			logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] r[%d]< 4 err[%s]", protocolName, connId, r, err)
			break
		}
		/// Last Successfull read from this connection
		startTime = time.Now()
		utilsPkg.TCPClientRxPackets++
		utilsPkg.TCPClientRxBytes += r

		readTimeout := 60 * time.Second
		err = cipEipClientConn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			cipEipClientConn.Close()
			logPkg.CtsLog.Error("c:protocolName[%s] connId[%d] remoteNetwork[%s] remoteAddr[%s] SetReadDeadline()fail\r\n", protocolName, connId, remoteNetwork, remoteAddr)
			return
		}

		// Check if received concatenated packets
		cipEipCmd := binary.LittleEndian.Uint16(buffer[0:2])
		cipEipLen := binary.LittleEndian.Uint16(buffer[2:4])
		rcvdCipEipClientRquPacket := buffer[:r]
		if (cipEipCmd == 0x0065) && len(rcvdCipEipClientRquPacket) > int(cipEipLen)+24 {
			///
			/// Double Client Rqu Packets Received
			///
			logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Rcvd Concatenated Rqu\n << rcvdCipEipClientRquPacket[%d:%X]Double Client Packet Received from remoteAddr[%s]", protocolName, connId, len(rcvdCipEipClientRquPacket), rcvdCipEipClientRquPacket, remoteAddr)
			///
			/// Proccess First Client Rqu Packet
			///
			rcvdCipEipClientRquPacket = buffer[:int(cipEipLen)+24]
			logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Rcvd First Client Rqu\n << rcvdCipEipClientRquPacket[%d:%X]First Packet from remoteAddr[%s]", protocolName, connId, len(rcvdCipEipClientRquPacket), rcvdCipEipClientRquPacket, remoteAddr)
			// Perform some operation on rcvdCipEipClientRquPacket (you can replace this with your logic)
			resultCode, sentClientRspPacket, err = procCipEipClientRquPacket(rcvdCipEipClientRquPacket)
			if err != nil || resultCode != utilsPkg.OK {
				logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER err[%s]", protocolName, connId, resultCode, err)
				break
			}

			// Send cipEipclientRspPacket back to the client
			w, err = cipEipClientConn.Write(sentClientRspPacket)
			if err != nil {
				utilsPkg.TCPClientTxErrors++
				resultCode = utilsPkg.ERR_COMM_WRITE
				logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER cipEipClientConn.Write err[%s]", protocolName, connId, resultCode, err)
				break
			}
			utilsPkg.TCPClientTxPackets++
			utilsPkg.TCPClientTxBytes += w
			/// Sent
			logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Sent First Client Rsp Packet\n >> sentClientRspPacket[%X] to remoteAddr[%s]\n\n", protocolName, connId, sentClientRspPacket, remoteAddr)

			///
			/// Proccess Second Client Rqu Packet
			///
			rcvdCipEipClientRquPacket = buffer[int(cipEipLen)+24 : r]
			logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Rcvd Second Client Rqu\n << rcvdCipEipClientRquPacket[%d:%X]from remoteAddr[%s]", protocolName, connId, len(rcvdCipEipClientRquPacket), rcvdCipEipClientRquPacket, remoteAddr)
			// Perform some operation on rcvdCipEipClientRquPacket (you can replace this with your logic)
			resultCode, sentClientRspPacket, err = procCipEipClientRquPacket(rcvdCipEipClientRquPacket)
			if err != nil || resultCode != utilsPkg.OK {
				logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER err[%s]", protocolName, connId, resultCode, err)
				break
			}

			// Send cipEipclientRspPacket back to the client
			w, err = cipEipClientConn.Write(sentClientRspPacket)
			if err != nil {
				utilsPkg.TCPClientTxErrors++
				resultCode = utilsPkg.ERR_COMM_WRITE
				logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER cipEipClientConn.Write err[%s]", protocolName, connId, resultCode, err)
				break
			}
			utilsPkg.TCPClientTxPackets++
			utilsPkg.TCPClientTxBytes += w
			/// Sent
			logPkg.CtsLog.Info("handleClientConnection:protocolName[%s] connId[%d]Sent Second Client Rsp Packet\n << sentClientRspPacket[%X] to remoteAddr[%s]\n\n", protocolName, connId, sentClientRspPacket, remoteAddr)
			continue
		}
		///
		/// Single Client Packet Received
		///

		///
		/// Proccess Client Rqu Packet
		///
		rcvdCipEipClientRquPacket = buffer[:r]
		logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Rcvd Client Rqu\n << rcvdCipEipClientRquPacket[%d:%X] from remoteAddr[%s]", protocolName, connId, len(rcvdCipEipClientRquPacket), rcvdCipEipClientRquPacket, remoteAddr)

		// Perform some operation on rcvdCipEipClientRquPacket (you can replace this with your logic)
		resultCode, sentClientRspPacket, err = procCipEipClientRquPacket(rcvdCipEipClientRquPacket)
		if err != nil || resultCode != utilsPkg.OK {
			logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER err[%s]", protocolName, connId, resultCode, err)
			break
		}

		// Send cipEipclientRspPacket back to the client
		w, err = cipEipClientConn.Write(sentClientRspPacket)
		if err != nil {
			utilsPkg.TCPClientTxErrors++
			resultCode = utilsPkg.ERR_COMM_WRITE
			logPkg.CtsLog.Error("handleClientConnection:protocolName[%s] connId[%d] resultCode[0x%X]ER cipEipClientConn.Write err[%s]", protocolName, connId, resultCode, err)
			break
		}
		if w > 0 && !NotifiedConnected {
			NotifiedConnected = true
			utilsPkg.TCPClientsConnected++
		}
		utilsPkg.TCPClientTxPackets++
		utilsPkg.TCPClientTxBytes += w
		logPkg.CtsLog.Debug("handleClientConnection:protocolName[%s] connId[%d]Sent Client Rsp Packet\n >> sentClientRspPacket[%X] to remoteAddr[%s]\n\n", protocolName, connId, sentClientRspPacket, remoteAddr)
	}
	if r != 0 || w != 0 {
		if utilsPkg.TCPClientsConnected > 0 {
			utilsPkg.TCPClientsConnected--
		}
		logPkg.CtsLog.Warn("handleClientConnection:protocolName[%s] connId[%d]resultCode[0x%X] remoteAddr[%s] Exiting! r[%d] w[%d]\r\n", protocolName, connId, int(resultCode), remoteAddr, r, w)
	}
}

// ======================================================================
// << procCipEipClientRquPacket - On Entry receives Client Request Packet
// ======================================================================
func procCipEipClientRquPacket(cipEipClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode uint8 = 0 // OK
	var iBytesCopyed int = 0
	var strCmdName string = ""
	var strErrorMsg string = ""
	var cipEipClientRspPacket []byte = nil
	if cipEipClientRquPacket == nil {
		// ER: cipEipClientRquPacket = nil
		ucResultCode = utilsPkg.ERR_BAD_PACKET_FORMAT
		strErrorMsg = "cipEipClientRquPacket = nil"
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) < 4 {
		// ER: cipEipClientRquPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d < 4", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) < 24 { // STX(1) + Cmd(1) + PayLoadLen(1) + ETX(1) + CRC16(2)
		// ER: cipEipClientRquPacket < 24
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d < 24", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) > 65536 {
		// ER: cipEipClientRquPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d > 65536", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Parsing CIP EthernetIP Request Header
	cipEipCmd := binary.LittleEndian.Uint16(cipEipClientRquPacket[0:2])
	cipEipLen := binary.LittleEndian.Uint16(cipEipClientRquPacket[2:4])
	if int(cipEipLen) != len(cipEipClientRquPacket)-24 {
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN
		strErrorMsg = fmt.Sprintf("int(cipEipLen)[%d] != len(cipEipClientRquPacket)-24[%d]", int(cipEipLen), len(cipEipClientRquPacket)-24)
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	cipEipSenderHandler := binary.LittleEndian.Uint32(cipEipClientRquPacket[4:8])
	cipEipStatus := binary.LittleEndian.Uint32(cipEipClientRquPacket[8:12])
	cipEipSenderCtx := binary.LittleEndian.Uint64(cipEipClientRquPacket[12:20])
	cipEipOptions := binary.LittleEndian.Uint32(cipEipClientRquPacket[20:24])
	/// Parsing CIP EthernetIP Request Data
	cipEipEncapData := make([]byte, cipEipLen)
	if cipEipLen > 0 {
		iBytesCopyed = copy(cipEipEncapData, cipEipClientRquPacket[24:cipEipLen+24+1])
	}
	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		strCmdName = "RegisterSession"
	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		strCmdName = "SENDRRRdata"
	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strCmdName = "UNKNOWN"
	}

	if len(cipEipClientRquPacket)-24 != int(cipEipLen) {
		// ER Bad cipEipClientEncapDataLen
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)-24[%d] != int(cipEipClientEncapDataLen)[%d]", len(cipEipClientRquPacket)-24, int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if int(cipEipLen) > int(65511) {
		// ER Bad cipEipClientEncapDataLen > 65511 Data Len Too Big
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN_TOO_BIG
		strErrorMsg = fmt.Sprintf("len(cipEipLen)[%d] > 65511", int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Show CIP EthernetIP Request Packet & Parsed Request data
	logPkg.CtsLog.Debug(">> Rqu cipEipClientRquPacket[%X]", cipEipClientRquPacket)
	logPkg.CtsLog.Debug("   Rqu cipEipCmd            [%04X]%s", cipEipCmd, strCmdName)
	logPkg.CtsLog.Debug("   Rqu cipEipLen            [%04X]", cipEipLen)
	logPkg.CtsLog.Debug("   Rqu cipEipSenderHandler  [%08X]", cipEipSenderHandler)
	logPkg.CtsLog.Debug("   Rqu cipEipStatus         [%08X]", cipEipStatus)
	logPkg.CtsLog.Debug("   Rqu cipEipSenderCtx      [%016X]", cipEipSenderCtx)
	logPkg.CtsLog.Debug("   Rqu cipEipOptions        [%04X]", cipEipOptions)
	logPkg.CtsLog.Debug("   Rqu cipEipEncapData      [%d:%X]", iBytesCopyed, cipEipEncapData)
	// Parsing CIP EthernetIP Request Data
	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		/// ===========================================================
		/// > Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H >SH = 00000000         4 Bytes Sender Handler(00000000)
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--ED-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// >> 65000400000000000000000000000000000000000000000001000000
		/// ===========================================================
		/// < Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H <SH*= 06000000         2 Bytes Sender Handler(00000006) From Server
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///P  PL = Variable         n Bytes PayLoad size LN
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// << 65000400060000000000000000000000000000000000000001000000
		/// ===========================================================
		cipEipVersion := binary.LittleEndian.Uint16(cipEipEncapData[0:2])
		cipEipFlags := binary.LittleEndian.Uint16(cipEipEncapData[2:4])
		logPkg.CtsLog.Debug("   Rqu cipEipVersion        [%04X]", cipEipVersion)
		logPkg.CtsLog.Debug("   Rqu cipEipFlags          [%04X]", cipEipFlags)
		/// Creating Client Response Packet
		cipEipClientRspPacket = cipEipClientRquPacket
		cipEipClientRspPacket[4] = 0x06 // Force SenderHandler= 0x000006

	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		/// ===========================================================
		/// >> Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN = 1A00             2 Bytes Len(001A) Len(26 Bytes)
		///H  SH = 06000000         4 Bytes Sender Handler(00000006)Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(26 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P >TO = 0A00             2 Bytes Timeout(000A)=10
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P  TI = 0000             2 Bytes Type ID Null Item(0000)
		///P  L1 = 0000             2 Bytes Len(0000)
		///P  TI = B200             2 Bytes Tipe ID Unconnected Data Item(B200)
		///P  L2 = 0A00             2 Bytes Len(000A)=10
		///P >SV = 0E               1 Byte  Service(0E) = GetAttributeSingleRqu) b7=0
		///P  RS = 04               1 Byte  RquPathSz(word)(04)
		///P  PS = 2100             2 Bytes Path Segment 1 16-bit Class Segment(0021)
		///                          001      Path Segment Type: Logical segment(1)
		///                             000   Logical Segment Type: Class ID(0)
		///                                01 Logical Segment Format: 16-bit segment
		///P  CL = 5003             2 Bytes Class(0350)
		///P  PS = 24               1 Byte  Path Segment 8-bit Instance Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             001   Logical Segment Type: Instance ID(1)
		///                                00 Logical Segment Format: 8-bit segment
		///P  IN = 01               1 Byte     Instance 01
		///P  PS = 30               1 Byte  Path Segment 8-bit Attribute Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             100   Logical Segment Type: Attribute ID(4)
		///                                00 Logical Segment Format: 8-bit segment
		///P  AT = 01               1 Byte     Attribute 01
		///                                                    <----------------------Payload--------------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TI><L1><TI><L2>SVRS<PS><CL>PSINPSAT
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
		/// >> 6F001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100500324013001
		/// ===========================================================
		/// << Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN*= 1500/1600        2 Bytes Len(0015/0x0016) Len(21/22 Bytes)
		///H *SH = 06000000         4 Bytes Sender Handler(00000006) Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(21/22 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P <TO = 0000             2 Bytes Timeout(0000)=0
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P <TP = 0000             2 Bytes Type(0000)
		///P <LT = 0000             2 Bytes Len(0000)
		///P <TI = B200             2 Bytes Type ID(00B2) Unconnected Data Item
		///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		///P <SV = 8E00             2 Byte  Service(008E)
		///P <ST = 0000             2 Bytes Status
		///P <DT*= XX               1 Byte  Data LN=(0005)
		///P <DT*= XXXX             2 Bytes Data LN=(0006)
		///                                                    <---------------------Payload---------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TP><LT><TI><LD><SV><ST><DT>
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		/// << 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E000000XX
		/// << 6F0016000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20006008E000000XXXX
		/// =================================================================================================
		cipEipCmdData := binary.LittleEndian.Uint32(cipEipEncapData[0:4])    ///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		cipEipTimeout := binary.LittleEndian.Uint16(cipEipEncapData[4:6])    ///P >TO = 0A00             2 Bytes Timeout(000A)=10
		cipEipItemCnt := binary.LittleEndian.Uint16(cipEipEncapData[6:8])    ///P  IC = 0200             2 Bytes ItemCount(0002)=2
		cipEipTupe1 := binary.LittleEndian.Uint16(cipEipEncapData[8:10])     ///P  TI = 0000             2 Bytes Type ID Null Item(0000)
		cipEipLen1 := binary.LittleEndian.Uint16(cipEipEncapData[10:12])     ///P  L1 = 0000             2 Bytes Len(0000)
		cipEipTupe2 := binary.LittleEndian.Uint16(cipEipEncapData[12:14])    ///P  TI = B200             2 Bytes Tipe ID Unconnected Data Item(B200)
		cipEipLen2 := binary.LittleEndian.Uint16(cipEipEncapData[14:16])     ///P  L2 = 0A00             2 Bytes Len(000A)=10
		cipEipSvc2 := cipEipEncapData[16]                                    ///P >SV = 0E               1 Byte  Service(0E) = GetAttributeSingleRqu) b7=0
		RquPathSz2 := cipEipEncapData[17]                                    ///P  RS = 04               1 Byte  RquPathSz(word)(04)
		cipEipPathSeg2 := binary.LittleEndian.Uint16(cipEipEncapData[18:20]) ///P  PS = 2100             2 Bytes Path Segment 1 16-bit Class Segment(0021)
		///                                                                              00000000 00100001b(0021h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                         *000   Logical Segment Type: Class ID(0)
		///                                                                                            *01 Logical Segment Format: 16-bit segment
		cipEipClass := binary.LittleEndian.Uint16(cipEipEncapData[20:22]) /// --P *CL = 5003             2 Bytes Class(0350)
		cipEipPathSeg3 := cipEipEncapData[22]                             /// --P  PS = 24               1 Byte  Path Segment 8-bit Attribute Segment
		///                                                                                       00100100b(0024h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                         *001   Logical Segment Type: Instance ID(1)
		///                                                                                            *00 Logical Segment Format: 8-bit segment
		cipEipInstanse := cipEipEncapData[23] /// ------------------------------P *IN = 01               1 Byte     Instance(01)
		cipEipPathSeg4 := cipEipEncapData[24] /// ------------------------------P  PS = 30               1 Byte  Path Segment 8-bit Instance Segment
		///                                                                                       00110000b(0030h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                          100   Logical Segment Type: Attribute ID(4)
		///                                                                                             00 Logical Segment Format: 8-bit segment
		cipEipAttribute := cipEipEncapData[25] /// -----------------------------P *AT = 01               1 Byte     Attribute(01)
		logPkg.CtsLog.Debug("   Rqu cipEipCmdData        [%08X]", cipEipCmdData)
		logPkg.CtsLog.Debug("   Rqu cipEipTimeout        [%04X]", cipEipTimeout)
		logPkg.CtsLog.Debug("   Rqu cipEipItemCnt        [%04X]", cipEipItemCnt)
		logPkg.CtsLog.Debug("   Rqu cipEipTupe1          [%04X]", cipEipTupe1)
		logPkg.CtsLog.Debug("   Rqu cipEipLen1           [%04X]", cipEipLen1)
		logPkg.CtsLog.Debug("   Rqu cipEipTupe2          [%04X]", cipEipTupe2)
		logPkg.CtsLog.Debug("   Rqu cipEipLen2           [%04X]", cipEipLen2)
		logPkg.CtsLog.Debug("   Rqu cipEipSvc2           [%02X]", cipEipSvc2)
		logPkg.CtsLog.Debug("   Rqu RquPathSz2           [%02X]", RquPathSz2)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg2       [%04X]", cipEipPathSeg2)
		logPkg.CtsLog.Debug("   Rqu cipEipClass          [(%04d)%04X]*", cipEipClass, cipEipClass)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg3       [%02X]", cipEipPathSeg3)
		logPkg.CtsLog.Debug("   Rqu cipEipInstanse       [(%02d)%02X]*", cipEipInstanse, cipEipInstanse)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg4       [%02X]", cipEipPathSeg4)
		logPkg.CtsLog.Debug("   Rqu cipEipAttribute      [(%02d)%02X]*", cipEipAttribute, cipEipAttribute)
		/// Create SendRRData Client Response Packet
		var tmpRspLen uint16 = 21
		if cipEipInstanse == 2 {
			tmpRspLen = 22
		}
		cipEipClientRspData := make([]byte, int(tmpRspLen))
		cipEipClientRspPacket = make([]byte, 24+int(tmpRspLen))
		binary.LittleEndian.PutUint16(cipEipClientRspPacket[0:2], cipEipCmd)           //  0-1  cipEipCmd
		binary.LittleEndian.PutUint16(cipEipClientRspPacket[2:4], tmpRspLen)           //  2-3  21 or 22
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[4:8], cipEipSenderHandler) //  4-7  cipEipSenderHdl
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[8:12], cipEipStatus)       //  8-11 cipEipStatus
		binary.LittleEndian.PutUint64(cipEipClientRspPacket[12:20], cipEipSenderCtx)   // 12-19 cipEipSenderContext
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[20:24], cipEipOptions)     // 20-23 cipEipSenderContext
		// CIP: Creating EthernetIP Response Data
		binary.LittleEndian.PutUint32(cipEipClientRspData[0:4], 0x00000000) ///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		binary.LittleEndian.PutUint16(cipEipClientRspData[4:6], 0x0000)     ///P <TO = 0000             2 Bytes Timeout(0000)=0
		binary.LittleEndian.PutUint16(cipEipClientRspData[6:8], 0x0002)     ///P  IC = 0200             2 Bytes ItemCount(0002)=2
		binary.LittleEndian.PutUint16(cipEipClientRspData[8:10], 0x0000)    ///P <TP = 0000             2 Bytes Type(0000)
		binary.LittleEndian.PutUint16(cipEipClientRspData[10:12], 0x0000)   ///P <LT = 0000             2 Bytes Len(0000)
		binary.LittleEndian.PutUint16(cipEipClientRspData[12:14], 0x00B2)   ///P <TI = B200             2 Bytes Type ID(00B2) Unconnected Data Item
		if cipEipInstanse == 1 {
			binary.LittleEndian.PutUint16(cipEipClientRspData[14:16], 0x0005) ///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		} else {
			binary.LittleEndian.PutUint16(cipEipClientRspData[14:16], 0x0006) ///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		}
		cipEipClientRspData[16] = 0x8E // Service
		cipEipClientRspData[17] = 0x00 //
		cipEipClientRspData[18] = 0x00 // Status
		cipEipClientRspData[19] = 0x00
		if cipEipInstanse == 1 {
			// BitMemory
			cipEipClientRspData[20] = byte(ushUpNewValue & 0x03)
		} else {
			// WordMemory
			cipEipClientRspData[20] = byte(ushUpNewValue & 0xFF)
			cipEipClientRspData[21] = byte(ushUpNewValue >> 8)
		}
		if utilsPkg.AutoUpdate {
			ushUpNewValue++
		}
		copy(cipEipClientRspPacket[24:], cipEipClientRspData)
		logPkg.CtsLog.Debug("   --- cipEipClientRspPacket[%d:%X]*", len(cipEipClientRspPacket), cipEipClientRspPacket)
		logPkg.CtsLog.Debug("   ---   cipEipClientRspData[%d:%X]*", len(cipEipClientRspData), cipEipClientRspData)

	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strErrorMsg = fmt.Sprintf("cipEipCmd[%04X]BAD_CMD", cipEipCmd)
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	/// Parse Client Created Client Response Packet before return it
	ucResultCode, cipEipClientRspPacketOut, err := ParseEipCipClientRspPacket(cipEipClientRspPacket)
	if err != nil {
		logPkg.CtsLog.Error("   --- cipEipClientRspPacketOut[%d:%X]\n err[%s]", len(cipEipClientRspPacketOut), cipEipClientRspPacketOut, err)
	} else {
		logPkg.CtsLog.Debug("   --- cipEipClientRspPacketOut[%d:%X]*", len(cipEipClientRspPacketOut), cipEipClientRspPacketOut)
	}
	return ucResultCode, cipEipClientRspPacket, err
}

// ====================================
// CIP EthernetIP Communication Example
// ====================================
// STEP01Rqu  << 65000400000000000000000000000000000000000000000001000000
// STEP01Rsp  >> 65000400060000000000000000000000000000000000000001000000
// STEP02Rqu  << 6F001A000600000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324013001
// STEP02Rsp  >> 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000000
// STEP03Rqu  >> 6F001A000600000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324013002
// STEP03Rsp  >> 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000021
// STEP04Rqu  << 6F001A000600000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324013003
// STEP04Rsp  >> 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000022
// STEP05Rqu  << 6F001A000600000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324013004
// STEP05Rsp  >> 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000023
// STEP06Rqu  << 6F001A000600000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100510324013005
// STEP06Rsp  >> 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000024
// STEP07Rqu  << 6F001A0006000000000000006F930AAB21996F7400000000000000000A00020000000000B2000A000E042100520324023000
// STEP07Rsp  >> 6F00160006000000000000006F930AAB21996F7400000000000000000000020000000000B20006008E0000002501

// =======================================================================================
// >> ParseEipCipClientRquPacket - On Successfull Exit returns with Client Response Packet
// =======================================================================================
func ParseCipEipClientRquPkg(cipEipClientRquPacket []byte) (byte, []byte, error) {
	var ucResultCode uint8 = 0 // OK
	var iBytesCopyed int = 0
	var strCmdName string = ""
	var strErrorMsg string = ""
	var cipEipClientRspPacket []byte = nil
	if cipEipClientRquPacket == nil {
		// ER: cipEipClientRquPacket = nil
		ucResultCode = utilsPkg.ERR_BAD_PACKET_FORMAT
		strErrorMsg = "cipEipClientRquPacket = nil"
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) < 4 {
		// ER: cipEipClientRquPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d < 4", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) < 24 { // STX(1) + Cmd(1) + PayLoadLen(1) + ETX(1) + CRC16(2)
		// ER: cipEipClientRquPacket < 24
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d < 24", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRquPacket) > 65536 {
		// ER: cipEipClientRquPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d > 65536", len(cipEipClientRquPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Parsing CIP EthernetIP Request Header
	cipEipCmd := binary.LittleEndian.Uint16(cipEipClientRquPacket[0:2])
	cipEipLen := binary.LittleEndian.Uint16(cipEipClientRquPacket[2:4])
	cipEipSenderHandler := binary.LittleEndian.Uint32(cipEipClientRquPacket[4:8])
	cipEipStatus := binary.LittleEndian.Uint32(cipEipClientRquPacket[8:12])
	cipEipSenderCtx := binary.LittleEndian.Uint64(cipEipClientRquPacket[12:20])
	cipEipOptions := binary.LittleEndian.Uint32(cipEipClientRquPacket[20:24])
	/// Parsing CIP EthernetIP Request Data
	cipEipEncapData := make([]byte, cipEipLen)
	if cipEipLen > 0 {
		iBytesCopyed = copy(cipEipEncapData, cipEipClientRquPacket[24:cipEipLen+24+1])
	}

	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		strCmdName = "RegisterSession"
	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		strCmdName = "SENDRRRdata"
	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strCmdName = "UNKNOWN"
	}

	if len(cipEipClientRquPacket)-24 != int(cipEipLen) {
		// ER Bad cipEipClientEncapDataLen
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)-24[%d] != int(cipEipClientEncapDataLen)[%d]", len(cipEipClientRquPacket)-24, int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if int(cipEipLen) > int(65511) {
		// ER Bad cipEipClientEncapDataLen > 65511 Data Len Too Big
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN_TOO_BIG
		strErrorMsg = fmt.Sprintf("len(cipEipLen)[%d] > 65511", int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Show CIP EthernetIP Request Packet & Parsed Request data
	logPkg.CtsLog.Debug(">> Rqu cipEipClientRquPacket[%X]", cipEipClientRquPacket)
	logPkg.CtsLog.Debug("   Rqu cipEipCmd            [%04X]%s", cipEipCmd, strCmdName)
	logPkg.CtsLog.Debug("   Rqu cipEipLen            [%04X]", cipEipLen)
	logPkg.CtsLog.Debug("   Rqu cipEipSenderHandler  [%08X]", cipEipSenderHandler)
	logPkg.CtsLog.Debug("   Rqu cipEipStatus         [%08X]", cipEipStatus)
	logPkg.CtsLog.Debug("   Rqu cipEipSenderCtx      [%016X]", cipEipSenderCtx)
	logPkg.CtsLog.Debug("   Rqu cipEipOptions        [%04X]", cipEipOptions)
	logPkg.CtsLog.Debug("   Rqu cipEipEncapData      [%d:%X]", iBytesCopyed, cipEipEncapData)
	// Parsing CIP EthernetIP Request Data
	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		/// ===========================================================
		/// > Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H >SH = 00000000         4 Bytes Sender Handler(00000000)
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--ED-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// >> 65000400000000000000000000000000000000000000000001000000
		/// ===========================================================
		/// < Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H <SH*= 06000000         2 Bytes Sender Handler(00000006) From Server
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///P  PL = Variable         n Bytes PayLoad size LN
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// << 65000400060000000000000000000000000000000000000001000000
		/// ===========================================================
		cipEipVersion := binary.LittleEndian.Uint16(cipEipEncapData[0:2])
		cipEipFlags := binary.LittleEndian.Uint16(cipEipEncapData[2:4])
		logPkg.CtsLog.Debug("   Rqu cipEipVersion        [%04X]", cipEipVersion)
		logPkg.CtsLog.Debug("   Rqu cipEipFlags          [%04X]", cipEipFlags)
		/// Create RegisterSession Client Response Packet
		cipEipClientRspPacket = cipEipClientRquPacket
		cipEipClientRspPacket[4] = 0x06 // Force Session Handlerr(0x000006)

	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		/// ===========================================================
		/// >> Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN = 1A00             2 Bytes Len(001A) Len(26 Bytes)
		///H  SH = 06000000         4 Bytes Sender Handler(00000006)Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(26 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P >TO = 0A00             2 Bytes Timeout(000A)=10
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P  TI = 0000             2 Bytes Type ID Null Item(0000)
		///P  L1 = 0000             2 Bytes Len(0000)
		///P  TI = B200             2 Bytes Tipe ID Unconnected Data Item(B200)
		///P  L2 = 0A00             2 Bytes Len(000A)=10
		///P >SV = 0E               1 Byte  Service(0E) = GetAttributeSingleRqu) b7=0
		///P  RS = 04               1 Byte  RquPathSz(word)(04)
		///P  PS = 2100             2 Bytes Path Segment 1 16-bit Class Segment(0021)
		///                          001      Path Segment Type: Logical segment(1)
		///                             000   Logical Segment Type: Class ID(0)
		///                                01 Logical Segment Format: 16-bit segment
		///P  CL = 5003             2 Bytes Class(0350)
		///P  PS = 24               1 Byte  Path Segment 8-bit Instance Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             001   Logical Segment Type: Instance ID(1)
		///                                00 Logical Segment Format: 8-bit segment
		///P  IN = 01               1 Byte     Instance 01
		///P  PS = 30               1 Byte  Path Segment 8-bit Attribute Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             100   Logical Segment Type: Attribute ID(4)
		///                                00 Logical Segment Format: 8-bit segment
		///P  AT = 01               1 Byte     Attribute 01
		///                                                    <----------------------Payload--------------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TI><L1><TI><L2>SVRS<PS><CL>PSINPSAT
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
		/// >> 6F001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100500324013001
		/// ===========================================================
		/// << Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN*= 1500/1600        2 Bytes Len(0015/0x0016) Len(21/22 Bytes)
		///H *SH = 06000000         4 Bytes Sender Handler(00000006) Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(21/22 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P <TO = 0000             2 Bytes Timeout(0000)=0
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P <TP = 0000             2 Bytes Type(0000)
		///P <LT = 0000             2 Bytes Len(0000)
		///P <TI = B200             2 Bytes Type ID(00B2) Unconnected Data Item
		///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		///P <SV = 8E00             2 Byte  Service(008E)
		///P <ST = 0000             2 Bytes Status
		///P <DT*= XX               1 Byte  Data LN=(0005)
		///P <DT*= XXXX             2 Bytes Data LN=(0006)
		///                                                    <---------------------Payload---------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TP><LT><TI><LD><SV><ST><=DT.>
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2
		/// << 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000000XX
		/// << 6F0016000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20006008E00000000XXXX
		/// =================================================================================================
		cipEipCmdData := binary.LittleEndian.Uint32(cipEipEncapData[0:4])    ///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		cipEipTimeout := binary.LittleEndian.Uint16(cipEipEncapData[4:6])    ///P >TO = 0A00             2 Bytes Timeout(000A)=10
		cipEipItemCnt := binary.LittleEndian.Uint16(cipEipEncapData[6:8])    ///P  IC = 0200             2 Bytes ItemCount(0002)=2
		cipEipTupe1 := binary.LittleEndian.Uint16(cipEipEncapData[8:10])     ///P  TI = 0000             2 Bytes Type ID Null Item(0000)
		cipEipLen1 := binary.LittleEndian.Uint16(cipEipEncapData[10:12])     ///P  L1 = 0000             2 Bytes Len(0000)
		cipEipTupe2 := binary.LittleEndian.Uint16(cipEipEncapData[12:14])    ///P  TI = B200             2 Bytes Tipe ID Unconnected Data Item(B200)
		cipEipLen2 := binary.LittleEndian.Uint16(cipEipEncapData[14:16])     ///P  L2 = 0A00             2 Bytes Len(000A)=10
		cipEipSvc2 := cipEipEncapData[16]                                    ///P >SV = 0E               1 Byte  Service(0E) = GetAttributeSingleRqu) b7=0
		RquPathSz2 := cipEipEncapData[17]                                    ///P  RS = 04               1 Byte  RquPathSz(word)(04)
		cipEipPathSeg2 := binary.LittleEndian.Uint16(cipEipEncapData[18:20]) ///P  PS = 2100             2 Bytes Path Segment 1 16-bit Class Segment(0021)
		///                                                                              00000000 00100001b(0021h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                         *000   Logical Segment Type: Class ID(0)
		///                                                                                            *01 Logical Segment Format: 16-bit segment
		cipEipClass := binary.LittleEndian.Uint16(cipEipEncapData[20:22]) /// --P *CL = 5003             2 Bytes Class(0350)
		cipEipPathSeg3 := cipEipEncapData[22]                             /// --P  PS = 24               1 Byte  Path Segment 8-bit Attribute Segment
		///                                                                                       00100100b(0024h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                         *001   Logical Segment Type: Instance ID(1)
		///                                                                                            *00 Logical Segment Format: 8-bit segment
		cipEipInstanse := cipEipEncapData[23] /// ------------------------------P *IN = 01               1 Byte     Instance(01)
		cipEipPathSeg4 := cipEipEncapData[24] /// ------------------------------P  PS = 30               1 Byte  Path Segment 8-bit Instance Segment
		///                                                                                       00110000b(0030h)
		///                                                                                       001      Path Segment Type: Logical segment(1)
		///                                                                                          100   Logical Segment Type: Attribute ID(4)
		///                                                                                             00 Logical Segment Format: 8-bit segment
		cipEipAttribute := cipEipEncapData[25] /// -----------------------------P *AT = 01               1 Byte     Attribute(01)
		logPkg.CtsLog.Debug("   Rqu cipEipCmdData        [%08X]", cipEipCmdData)
		logPkg.CtsLog.Debug("   Rqu cipEipTimeout        [%04X]", cipEipTimeout)
		logPkg.CtsLog.Debug("   Rqu cipEipItemCnt        [%04X]", cipEipItemCnt)
		logPkg.CtsLog.Debug("   Rqu cipEipTupe1          [%04X]", cipEipTupe1)
		logPkg.CtsLog.Debug("   Rqu cipEipLen1           [%04X]", cipEipLen1)
		logPkg.CtsLog.Debug("   Rqu cipEipTupe2          [%04X]", cipEipTupe2)
		logPkg.CtsLog.Debug("   Rqu cipEipLen2           [%04X]", cipEipLen2)
		logPkg.CtsLog.Debug("   Rqu cipEipSvc2           [%02X]", cipEipSvc2)
		logPkg.CtsLog.Debug("   Rqu RquPathSz2           [%02X]", RquPathSz2)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg2       [%04X]", cipEipPathSeg2)
		logPkg.CtsLog.Debug("   Rqu cipEipClass          [%04X]*", cipEipClass)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg3       [%02X]", cipEipPathSeg3)
		logPkg.CtsLog.Debug("   Rqu cipEipInstanse       [%02X]*", cipEipInstanse)
		logPkg.CtsLog.Debug("   Rqu cipEipPathSeg4       [%02X]", cipEipPathSeg4)
		logPkg.CtsLog.Debug("   Rqu cipEipAttribute      [%02X]*", cipEipAttribute)
		/// Create SendRRData Client Response Packet
		var cipEipRquData []byte
		if cipEipInstanse == 1 {
			cipEipClientRspPacket = make([]byte, 24+21)
			cipEipRquData = make([]byte, 21)
		} else {
			cipEipClientRspPacket = make([]byte, 24+22)
			cipEipRquData = make([]byte, 21)
		}
		// CIP: EthernetIP Response Header
		binary.LittleEndian.PutUint16(cipEipClientRspPacket[0:2], cipEipCmd) //  0-1  cipEipCmd
		if cipEipInstanse == 1 {
			binary.LittleEndian.PutUint16(cipEipClientRspPacket[2:4], 21) //  2-3  21
		} else {
			binary.LittleEndian.PutUint16(cipEipClientRspPacket[2:4], 22) //  2-3  22
		}
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[4:8], cipEipSenderHandler) //  4-7  cipEipSenderHdl
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[8:12], cipEipStatus)       //  8-11 cipEipStatus
		binary.LittleEndian.PutUint64(cipEipClientRspPacket[12:20], cipEipSenderCtx)   // 12-19 cipEipSenderContext
		binary.LittleEndian.PutUint32(cipEipClientRspPacket[20:24], cipEipOptions)     // 20-23 cipEipSenderContext
		// CIP: EthernetIP Response Data
		copy(cipEipClientRspPacket[24:], cipEipRquData)

	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strErrorMsg = fmt.Sprintf("cipEipCmd[%04X]BAD_CMD", cipEipCmd)
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	/// Parse Client Created Client Response Packet before return it
	ucResultCode, _, err := ParseEipCipClientRspPacket(cipEipClientRspPacket)
	return ucResultCode, cipEipClientRspPacket, err
}

// =======================================================================================
// << ParseEipCipClientRspPacket - On Successfull Exit returns with Client Response Packet
// =======================================================================================
func ParseEipCipClientRspPacket(cipEipClientRspPacket []byte) (byte, []byte, error) {
	var ucResultCode uint8 = 0 // OK
	var iBytesCopyed int = 0
	var strCmdName string = ""
	var strErrorMsg string = ""

	if cipEipClientRspPacket == nil {
		// ER: cipEipClientRspPacket = nil
		ucResultCode = utilsPkg.ERR_BAD_PACKET_FORMAT
		strErrorMsg = "cipEipClientRspPacket = nil"
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRspPacket) < 4 {
		// ER: cipEipClientRspPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRquPacket)=%d < 4", len(cipEipClientRspPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRspPacket) < 24 { // STX(1) + Cmd(1) + PayLoadLen(1) + ETX(1) + CRC16(2)
		// ER: cipEipClientRspPacket < 24
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRspPacket)=%d < 24", len(cipEipClientRspPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if len(cipEipClientRspPacket) > 65536 {
		// ER: cipEipClientRspPacket < 4
		ucResultCode = utilsPkg.ERR_BAD_PACKET_TOO_SHORT
		strErrorMsg = fmt.Sprintf("len(cipEipClientRspPacket)=%d > 65536", len(cipEipClientRspPacket))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Parsing CIP EthernetIP Response Header
	cipEipCmd := binary.LittleEndian.Uint16(cipEipClientRspPacket[0:2])
	cipEipLen := binary.LittleEndian.Uint16(cipEipClientRspPacket[2:4])
	cipEipSenderHandler := binary.LittleEndian.Uint32(cipEipClientRspPacket[4:8])
	cipEipStatus := binary.LittleEndian.Uint32(cipEipClientRspPacket[8:12])
	cipEipSenderCtx := binary.LittleEndian.Uint64(cipEipClientRspPacket[12:20])
	cipEipOptions := binary.LittleEndian.Uint32(cipEipClientRspPacket[20:24])
	/// Parsing CIP EthernetIP Response Data
	cipEipEncapData := make([]byte, cipEipLen)
	if cipEipLen > 0 {
		iBytesCopyed = copy(cipEipEncapData, cipEipClientRspPacket[24:])
	}
	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		strCmdName = "RegisterSession"
	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		strCmdName = "SENDRRRdata"
	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strCmdName = "UNKNOWN"
	}

	if len(cipEipClientRspPacket)-24 != int(cipEipLen) {
		// ER Bad cipEipClientEncapDataLen
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN
		strErrorMsg = fmt.Sprintf("len(cipEipClientRspPacket)-24[%d] != int(cipEipClientEncapDataLen)[%d]", len(cipEipClientRspPacket)-24, int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	} else if int(cipEipLen) > int(65511) {
		// ER Bad cipEipClientEncapDataLen > 65511 Data Len Too Big
		ucResultCode = utilsPkg.ERR_BAD_PACKET_DATA_LEN_TOO_BIG
		strErrorMsg = fmt.Sprintf("len(cipEipLen)[%d] > 65511", int(cipEipLen))
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	// Show CIP EthernetIP Response Packet & Parsed Request data
	logPkg.CtsLog.Debug("<< Rsp cipEipClientRspPacket[%X]", cipEipClientRspPacket)
	logPkg.CtsLog.Debug("   Rsp cipEipCmd            [%04X]%s", cipEipCmd, strCmdName)
	logPkg.CtsLog.Debug("   Rsp cipEipLen            [%04X]", cipEipLen)
	logPkg.CtsLog.Debug("   Rsp cipEipSenderHandler  [%08X]***", cipEipSenderHandler)
	logPkg.CtsLog.Debug("   Rsp cipEipStatus         [%08X]", cipEipStatus)
	logPkg.CtsLog.Debug("   Rsp cipEipSenderCtx      [%016X]", cipEipSenderCtx)
	logPkg.CtsLog.Debug("   Rsp cipEipOptions        [%04X]", cipEipOptions)
	logPkg.CtsLog.Debug("   Rsp cipEipEncapData      [%d:%X]", iBytesCopyed, cipEipEncapData)

	// Parsing CIP EthernetIP Response Data
	switch cipEipCmd {
	case utilsPkg.CIP_EIP_SVC_CODE_REGISTER_SESSION: // 0x0065
		/// ===========================================================
		/// > Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H >SH = 00000000         4 Bytes Sender Handler(00000000)
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--ED-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// >> 65000400000000000000000000000000000000000000000001000000
		/// ===========================================================
		/// < Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6500             2 Bytes ushRquCmd(0065)Register Session
		///H  LN = 0400             2 Bytes Len(0004)
		///H <SH*= 06000000         2 Bytes Sender Handler(00000006) From Server
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = 0000000000000000 8 Bytes Sender Context(0000000000000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///P  PL = Variable         n Bytes PayLoad size LN
		///   ED = Encapsulation Data(4 Bytes)
		///P *ED = 01000000         4 Bytes Command Specific Encapsulation Data(00000001)
		///        VVVV            *2 Bytes Protocol Version(0001)
		///            FFFF        *2 Bytes Option Flags(0000)
		///                                                    <Payld->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD-->
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3
		/// << 65000400060000000000000000000000000000000000000001000000
		/// ===========================================================
		cipEipVersion := binary.LittleEndian.Uint16(cipEipEncapData[0:2])
		cipEipFlags := binary.LittleEndian.Uint16(cipEipEncapData[2:4])
		logPkg.CtsLog.Debug("   Rsp cipEipVersion        [%04X]", cipEipVersion)
		logPkg.CtsLog.Debug("   Rsp cipEipFlags          [%04X]", cipEipFlags)
		/// Creating Client Response Packet

	case utilsPkg.CIP_EIP_SVC_CODE_SEND_RR_DATA: // 0x006F
		/// ===========================================================
		/// >> Request
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN = 1A00             2 Bytes Len(001A) Len(26 Bytes)
		///H  SH = 06000000         4 Bytes Sender Handler(00000006)Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(26 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P >TO = 0A00             2 Bytes Timeout(000A)=10
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P  TI = 0000             2 Bytes Type ID Null Item(0000)
		///P  L1 = 0000             2 Bytes Len(0000)
		///P  TI = B200             2 Bytes Tipe ID Unconnected Data Item(B200)
		///P  L2 = 0A00             2 Bytes Len(000A)=10
		///P >SV = 0E               1 Byte  Service(0E) = GetAttributeSingleRqu) b7=0
		///P  RS = 04               1 Byte  RquPathSz(word)(04)
		///P  PS = 2100             2 Bytes Path Segment 1 16-bit Class Segment(0021)
		///                          001      Path Segment Type: Logical segment(1)
		///                             000   Logical Segment Type: Class ID(0)
		///                                01 Logical Segment Format: 16-bit segment
		///P  CL = 5003             2 Bytes Class(0350)
		///P  PS = 24               1 Byte  Path Segment 8-bit Instance Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             001   Logical Segment Type: Instance ID(1)
		///                                00 Logical Segment Format: 8-bit segment
		///P  IN = 01               1 Byte     Instance 01
		///P  PS = 30               1 Byte  Path Segment 8-bit Attribute Segment
		///                          001      Path Segment Type: Logical segment(1)
		///                             100   Logical Segment Type: Attribute ID(4)
		///                                00 Logical Segment Format: 8-bit segment
		///P  AT = 01               1 Byte     Attribute 01
		///                                                    <----------------------Payload--------------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TI><L1><TI><L2>SVRS<PS><CL>PSINPSAT
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
		/// >> 6F001A000000000000000000FFFFFFFF0000000000000000000000000A00020000000000B2000A000E042100500324013001
		/// ===========================================================
		/// << Response
		/// ===========================================================
		///   EH = Encapsulation Header(24 Bytes)
		///H  CM = 6F00             2 Bytes ushRquCmd(006F)Send RR Data
		///H  LN*= 1500/1600        2 Bytes Len(0015/0x0016) Len(21/22 Bytes)
		///H *SH = 06000000         4 Bytes Sender Handler(00000006) Same Register Session Rsp
		///H  ST = 00000000         4 Bytes Status(00000000)success)
		///H  SC = FFFFFFFF00000000 8 Bytes Sender Context(FFFFFFFF00000000)
		///H  OP = 00000000         4 Bytes Options(00000000)
		///   ED = Encapsulation Data(21/22 Bytes)
		///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		///P <TO = 0000             2 Bytes Timeout(0000)=0
		///P  IC = 0200             2 Bytes ItemCount(0002)=2
		///P <TP = 0000             2 Bytes Type(0000)
		///P <LT = 0000             2 Bytes Len(0000)
		///P <TI = B200             2 Bytes Type ID(00B2) Unconnected Data Item
		///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		///P <SV = 8E00             2 Byte  Service(008E)
		///P <ST = 0000             2 Bytes Status
		///P <DT*= XX               1 Byte  Data LN=(0005)
		///P <DT*= XXXX             2 Bytes Data LN=(0006)
		///                                                    <---------------------Payload---------------->
		///    <CM><LN><--SH--><--ST--><-------SC-----><--OP--><--CD--><TO><IC><TP><LT><TI><LD><SV><ST><=DT.>
		///     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2
		/// << 6F0015000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20005008E00000000XX
		/// << 6F0016000600000000000000FFFFFFFF0000000000000000000000000000020000000000B20006008E00000000XXXX
		/// =================================================================================================

		///P <DT*= XX               1 Byte  Data LN=(0005)

		cipEipCmdData := binary.LittleEndian.Uint32(cipEipEncapData[0:4])   ///P  CD = 00000000         4 Bytes Command Specific Data(00000000) = CIP
		cipEipTimeout := binary.LittleEndian.Uint16(cipEipEncapData[4:6])   ///P <TO = 0000             2 Bytes Timeout(0000)=0
		cipEipItemCnt := binary.LittleEndian.Uint16(cipEipEncapData[6:8])   ///P  IC = 0200             2 Bytes ItemCount(0002)=2
		cipEipType1 := binary.LittleEndian.Uint16(cipEipEncapData[8:10])    ///P <TP = 0000             2 Bytes Type(0000)
		cipEipLen1 := binary.LittleEndian.Uint16(cipEipEncapData[10:12])    ///P <LT = 0000             2 Bytes Len(0000)
		cipEipType2 := binary.LittleEndian.Uint16(cipEipEncapData[12:14])   ///P <TI = B200             2 Bytes Type ID(00B2) Unconnected Data Item
		cipEipLen2 := binary.LittleEndian.Uint16(cipEipEncapData[14:16])    ///P <LD*= 0500/0600        2 Bytes Len(0005)DT=XX or(0006) DT=XXXX
		cipEipService := binary.LittleEndian.Uint16(cipEipEncapData[16:18]) ///P <SV = 8E00             2 Byte  Service(008E)
		cipEipStatus := binary.LittleEndian.Uint16(cipEipEncapData[18:20])  ///P <ST = 0000             2 Bytes Status
		var cipEipData8 uint8 = 0xFF
		var cipEipData16 uint16 = 0xFFFF
		if cipEipLen2 == 5 {
			// VarData = 8bits
			cipEipData8 = cipEipEncapData[20] ///-------------------------------P <DT*= XXXX             2 Bytes Data LN=(0006)
		} else if cipEipLen2 == 6 {
			// VarData = 16bits
			cipEipData16 = binary.LittleEndian.Uint16(cipEipEncapData[20:22]) ///P <DT*= XXXX             2 Bytes Data LN=(0006)
		}
		logPkg.CtsLog.Debug("   Rsp cipEipCmdData        [%08X]", cipEipCmdData)
		logPkg.CtsLog.Debug("   Rsp cipEipTimeout        [%04X]", cipEipTimeout)
		logPkg.CtsLog.Debug("   Rsp cipEipItemCnt        [%04X]", cipEipItemCnt)
		logPkg.CtsLog.Debug("   Rsp cipEipType1          [%04X]", cipEipType1)
		logPkg.CtsLog.Debug("   Rsp cipEipLen1           [%04X]", cipEipLen1)
		logPkg.CtsLog.Debug("   Rsp cipEipType2          [%04X]", cipEipType2)
		logPkg.CtsLog.Debug("   Rsp cipEipLen2           [%04X]", cipEipLen2)
		logPkg.CtsLog.Debug("   Rsp cipEipService        [%04X]", cipEipService)
		logPkg.CtsLog.Debug("   Rsp cipEipStatus         [%04X]", cipEipStatus)
		if cipEipLen2 == 5 {
			logPkg.CtsLog.Debug("   Rsp cipEipData8          [%02X]", cipEipData8)
		} else if cipEipLen2 == 6 {
			logPkg.CtsLog.Debug("   Rsp cipEipData16         [%04X]", cipEipData16)
		} else {
			logPkg.CtsLog.Error("   Rsp cipEipLen2          [%d]ER: Expected 5 or 6", cipEipLen2)
		}

	default:
		ucResultCode = utilsPkg.ERR_BAD_CMD
		strErrorMsg = fmt.Sprintf("cipEipCmd[%04X]BAD_CMD", cipEipCmd)
		logPkg.CtsLog.Error(strErrorMsg)
		return ucResultCode, cipEipClientRspPacket, fmt.Errorf(strErrorMsg)
	}
	return ucResultCode, cipEipClientRspPacket, nil
}
