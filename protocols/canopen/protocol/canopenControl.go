package canopen

import (
	logPkg "canopen/utils/gologtofile"
	"encoding/binary"
	"fmt"
)

// OSI Layer 1: Physical layer
// PLS Physical Layer Signaling   ISO 11898-1(CAN frame type)
// MDI Medium Dependent Interface ISO 11898-2(High Speed), ISO 11898-3
// PMA Physical medium attachment ISO 11898-2(High Speed), ISO 11898-3
// OSI Layer 2:  Data link layer ISO 11898-1(CAN frame type)
// OSI Layer 7:  Application layer (For numerical data types the encoding is little endian style)

// IMPORTANT: same struct of *can.frame from can.bus
type CANOPEN_Frame struct {
	//            Idx
	/// A standard CAN frame has an 11-bit ID
	/// A extended CAN frame has an 29-bit ID
	// CAN-ID(11 bits) = FunctionCode(4 bits) NODE-ID( 7 bits) + RTR(1 bit)(Remote Transmision Request)
	// CAN-ID(29 bits) = FunctionCode(4 bits) NODE-ID(25 bits) + RTR(1 bit)(Remote Transmision Request)
	// 7.1.3 Bit sequences
	// Examples: 10110100b, 1b, 101b, etc. are bit sequences
	ID  uint32 //  0-3
	DLC uint8  //  4  Byte  4 b0-3, RTR Byte 4 b6, EXT byte 4 b7(Remote Transmision Request)
	/// Up to 8 bytes
	Data []byte // 5  Bytes 5-12 0-8 bytes Data[0]=CANopen ControlByte
}

// CAN 2.0A *SOF 1 bit + ID=11 bits + *RTR 1 bit. *Control 6 bits + Data 0-64 bits + CRC 16 bits + EOF 7 bits(Remote Transmision Request)
// CAN 2.0B *SOF 1 bit + ID=29 bits + *RTR 1 bit. *Control 6 bits + Data 0-64 bits + CRC 16 bits + EOF 7 bits(Remote Transmision Request)
type CANopen_Frame struct {
	// CAN-ID(11 bits) = FunctionCode(4 bits) NODE-ID( 7 bits) + RTR(1 bit)(Remote Transmision Request)
	// CAN-ID(29 bits) = FunctionCode(4 bits) NODE-ID(25 bits) + RTR(1 bit)(Remote Transmision Request)
	FunctionCode uint8 // ID=11 bits 4 bits - Byte 0 b4-7
	NodeID       uint8 // ID=11 bits 7 bits - Byte 0 b0-3 &  Byte 1 B7-6
	DLC          uint8 // 4 bits - Byte 3 b0-3
	/// Up to 8 bytes
	Data []byte // 5
}

type COB_ID struct {
	FunctionCode uint8 // COB_ID 3 bits - b0-b3  BitEndian
	NODE_ID      uint8 // COB_ID 7 bits - b4-b10 BigEndian
}

type CANOpenData struct {
	ControlByte uint8  // CANopenFrame byte  5
	DataBytes   []byte // CANopenFrame Bytes 6-13
}

type StructCANOpenCtlByteInfo struct {
	cs uint8 // controlByte 3 bits b5-7
	x  uint8 // controlByte 1 bit  b4
	n  uint8 // controlByte 2 bits b2-3
	s  uint8 // controlByte 1 bit  b1
	e  uint8 // controlByte 1 bit  b0
}

var CANOpenCtlByteInfo StructCANOpenCtlByteInfo

// + =========================================================== +
// + CANopen Object Dictionary groups                            +
// + =========================================================== +
// + Index       + Object Group                                  +
//   =========================================================== +
// + 0000h       + Reserved                                      +
// + 0001h-009Fh + Static and Complex Data Types                 +
// + 00A0h-0FFFh + Reserved                                      +
// + 1000h-1FFFh + Communication profiles (e.g., DS 301, DS 302) +
// +             + Ex:PDO Range  1400h-1A7Fh                     +
// + 2000h-25FFh + Manufacturer Specific Device Profile          +
// + 6000h-9FFFh + Standardized device profiles                  +
// + A000h-FFFFh + Reserved                                      +

// =============================================================================================================== +
// CANopen Object Dictionary Communication profiles (e.g., DS 301, DS 302)                                         +
// =============================================================================================================== +
// + Index + Object Name    + Sub-Index + Description                      + Type + Access + Notes                 +
// =============================================================================================================== +
// + 1000h + Device Type    + 00h       + Type of Device                   + U32  + RO     + 0000 0000h No Profile +
// + 1001h + Error Register + 00h       + Error Register connected to EMCY + U32  + RO     + 0000 0000h No Profile +
// +                                    + bit 0 indicate generic error     + U8   + RO     + -                     +
// + 1003h + pre-defined    + 00h       + Number of Errors Writting a 0    + U32  + RO     +                       +
// +         error field    +           + to this sub-index clear the      + U8   + RW     + -                     +
// +                        +           + Error list                       +                                       +
// +                        + 01h-10h   + List of Error Recent at top      +                                       +
// + 1005h + COB-ID Sinc    + 00h       + ID of Sync message               + U32  + RO     + -                     +
// + 1006h + Comm Cycle     + 00h       + Comm Cycle Periodh               + U32  + RW     + -                     +
// + 1007h + Sync Win Len   + 00h       + Sync Windows Length              + U32  + RW     + Available Sync Supp   +
// + 1008h + ManDevName     + 00h       + Manufacturer Device Name         + str  + RO     + SI CANopen            +
// + 1009h + ManHwrVer      + 00h       + Manufactorer Hardware Version    + str  + RO     + Hwr Version           +
// + 100Ah + ManSwrVer      + 00h       + Manufacturer Software Version    + str  + RW     + Swr Version           +
// + 100Ch + GuardTime      + 00h       + Node Life Time(ms)               + U16  + RW     + 0000h(Default)        +
// + 100Dh + LifeTimeFactor + 00h       + Life Time Factor(ms)             + U16  + RW     + 0000h(Default)        +

// =============================================================================================== +
// CANopen PDO Object Dictionary Communication profiles                                            +
// =============================================================================================== +
// + Index        + Object Name          + Sub-Index + Description                 + Type + Access +
// =============================================================================================== +
// + 1400h-147Fh + Receive PDO Param     + 00h       + Largest Subindex Supported  + U32  + RO     +
// +                                     + 01h       + COB-ID used in PDO          + U32  + RW     +
// +                                     + 02h       + Transmission type           + U8   + RW     +
// + 1600h-167Fh + Receive PDO Mapping   + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
// +                                     + 01h       + Mapped Object #1            + U32  + RW     +
// +                                     + 02h       + Mapped Object #2            + U32  + RW     +
// +                                     + 03h       + Mapped Object #3            + U32  + RW     +
// +                                     + 04h       + Mapped Object #4            + U32  + RW     +
// +                                     + 05h       + Mapped Object #5            + U32  + RW     +
// +                                     + 06h       + Mapped Object #6            + U32  + RW     +
// +                                     + 07h       + Mapped Object #7            + U32  + RW     +
// +                                     + 08h       + Mapped Object #8            + U32  + RW     +
// + 1800h-187Fh + Transmit PDO Param    + 00h       + Largest Subindex Supported  + U32  + RO     +
// +                                     + 01h       + COB-ID used in PDO          + U32  + RW     +
// +                                     + 02h       + Transmission type           + U8   + RW     +
// +                                     + 03h       + Inhib Time(ms)              + U16  + RW     +
// +                                     + 05h       + Event TIme(ms)              + U16  + RW     +
// + 1A00h-1A7Fh + Transmit PDO Mapping  + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
// +                                     + 01h       + Mapped Object #1            + U32  + RW     +
// +                                     + 02h       + Mapped Object #2            + U32  + RW     +
// +                                     + 03h       + Mapped Object #3            + U32  + RW     +
// +                                     + 04h       + Mapped Object #4            + U32  + RW     +
// +                                     + 05h       + Mapped Object #5            + U32  + RW     +
// +                                     + 06h       + Mapped Object #6            + U32  + RW     +
// +                                     + 07h       + Mapped Object #7            + U32  + RW     +
// +                                     + 08h       + Mapped Object #8            + U32  + RW     +
// =============================================================================================== +

// ============================================================================== +
// CANopen NMT Object Dictionary                                                                                   +
// ========================================================================================== +
// + Index + Object Name                + Sub-Index + Description                      + Type +
// ========================================================================================== +
// + 1F80h + NMT Start                  + -         + Who is the MASTER                + U32  +
// + 1F81h + Slave Assign               + ARRAY     + All Slaves to be Managed         + U32  +
// + 1F82h + Rqu NMT                    + ARRAY     + Remote Control Initializatio     + U32  +
// + 1F83h + Rqu Guard                  + ARRAY     + Remote Control Start/Stop        + U32  +
// + 1F84h + Device Type Identification + ARRAY     + Verify Dev Type for Slave        + U32  +
// + 1F85h + Vendor Identification      + ARRAY     + VID for Slaves                   + U32  +
// + 1F86h + Protuct Code               + ARRAY     + PID for Slaves                   + U32  +
// + 1F87h + Revision Number            + ARRAY     + Revision for Slaves              + U32  +
// + 1F88h + Serial Number              + ARRAY     + Verify SN For Slaves             + U32  +
// + 1F89h + Boot Time                  + ARRAY     + Minnimum slaves Boot Time        + U32  +

// =============================
// CANopen: Initialization Steps
// =============================
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// STEP01 - SDO Read Request: Device A initiates an SDO read request to read data from Device B's Object Dictionary
// >> COB-ID: 600h DLC: 8 Data: 40 01 20 03 00 00 00 00    [COB-ID: 600h + NodeID]
//
//	COB-ID: 600h (hexadecimal) - This is the CAN message identifier, including the node ID of Device B (0x20) and
//	                             the SDO bit set for client-server communication (0x40)
//
// DLC: 8 - This is the Data Length Code indicating the number of bytes in the data field.
// Data:
// 40 - SDO Request (client command specifier) with expedited transfer and size indicated.
// 01 - Object Index (high byte)
// 20 - Object Index (low byte)
// 03 - Sub-Index
// 00 00 00 00 - Data bytes (not applicable for read requests)
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// STEP02 - SDO Read Response: Device B (slave) receives the SDO read request and processes it.
// << COB-ID: 580h DLC: 8 Data: 43 01 20 03 2A 00 00 00  [COB-ID: 580h + NodeID]
//
//	COB-ID: 580h (hexadecimal) - This is the CAN message identifier, including the node ID of Device A (0x01) and
//	                             the SDO bit set for server-client communication (0x43).
//
// DLC: 8 - This is the Data Length Code indicating the number of bytes in the data field.
// Data:
// 43 - SDO Response (server response) with expedited transfer and size indicated.
// 01 - Object Index (high byte)
// 20 - Object Index (low byte)
// 03 - Sub-Index
// 2A 00 00 00 - Data bytes (requested data value)
//
// This Function decodes detailed of CANopen byte array message
// if all required data are present return nil else return error message explaining the reason of failure
func ParseRawCANopenFrame(rawCANOpenFrame []byte) error {
	// Converte Raw CAN Frame to CAN frame
	if rawCANOpenFrame == nil {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER can frame=nil CANRawFrame[%X]", len(rawCANOpenFrame))
		return fmt.Errorf("ParseRawCANopenFrame:ER can frame=nil CANRawFrame[%X]", len(rawCANOpenFrame))
	} else if len(rawCANOpenFrame) < 4 {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER can frame Too Short Exp len(CANRawFrame) >=4 rcv len(CANRawFrame)[%d]", len(rawCANOpenFrame))
		return fmt.Errorf("ParseRawCANopenFrame:ER can frame Too Short Exp len(CANRawFrame) >=4 rcv len(CANRawFrame)[%d]", len(rawCANOpenFrame))
	} else if len(rawCANOpenFrame) > 13 {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER can frame Too long Exp len(data) <=13 rcv len(data)[%d]", len(rawCANOpenFrame))
		return fmt.Errorf("ParseRawCANopenFrame:ER can frame too long Exp len(data) <=13 rcv len(data)[%d]", len(rawCANOpenFrame))
	}

	id := binary.BigEndian.Uint32(rawCANOpenFrame[:4]) // Bytes 0-3
	dlc := rawCANOpenFrame[4] & 0x0F                   // Byte 4 b0-3
	ext := (rawCANOpenFrame[4] & 0x80) != 0            // Byte 4 b7
	rtr := (rawCANOpenFrame[4] & 0x40) != 0            // Byte 4 b6
	if len(rawCANOpenFrame) < int(dlc)+5 {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER invalid can frame len(CANRawFrame)[%d] < int(dlc)+4[%d]", len(rawCANOpenFrame), int(dlc)+5)
		return fmt.Errorf("ParseRawCANopenFrame:ER invalid can frame len(CANRawFrame)[%d] < int(dlc)+4[%d]", len(rawCANOpenFrame), int(dlc)+5)
	} else if len(rawCANOpenFrame) > int(dlc)+5 {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER invalid can frame len(data)[%d] > int(dlc)+4[%d]", len(rawCANOpenFrame), int(dlc)+5)
		return fmt.Errorf("ParseRawCANopenFrame:ER invalid can frame len(data)[%d] > int(dlc)+4[%d]", len(rawCANOpenFrame), int(dlc)+5)
	}
	canFrame := &CANOPEN_Frame{
		// Standard  CanFrame(11-bits b0-10 BigEndian)+RTR(1 bit b11 BigEndian) 12 bits(Remote Transmision Request)
		// Extended CanFrrame(29 bits b0-28 BigEndian)+RTR 1 bit b20 BigEndian) 30 bits(Remote Transmision Request)
		ID:   id,                         // 32 bits bytes 0-3
		DLC:  dlc,                        //  4 Bits byte  4 bits (0-3) byte 4 b7 ext byte 4 b6 RNR byte 4 bits 4-5 N/A
		Data: rawCANOpenFrame[5 : 5+dlc], //  8 Bits Bytes 5-12
	}
	if canFrame.DLC != uint8(len(canFrame.Data)) {
		logPkg.CtsLog.Error("ParseRawCANopenFrame:ER invalid can frame: canFrame DLC[%d] != len(canFrame Data)[%d]", canFrame.DLC, len(canFrame.Data))
		return fmt.Errorf("ParseRawCANopenFrame:ER invalid can frame: canFrame DLC[%d] != len(canFrame Data)[%d]", canFrame.DLC, len(canFrame.Data))
	}
	if !ext {
		// CIA301a: ID = 11 bits
		err := ParseCANopen_CIA301aFrame(canFrame)
		if err != nil {
			///FAIL to Parse a CANopen Frame
			logPkg.CtsLog.Error("ParseRawCANopenFrame:ER err[%s]", err)
			return err
		}
	} else {
		// CIA301b: ID = 29 bits
		err := ParseCANopen_CIA301bFrame(canFrame)
		if err != nil {
			///FAIL to Parse a CANopen Frame
			logPkg.CtsLog.Error("ParseRawCANopenFrame:ER err[%s]", err)
			return err
		}
	}
	logPkg.CtsLog.Debug("ParseRawCANopenFrame:OK << ParseRawCANopenFrame[%X] id[%X] dlc[%d] ext[%v] rtr[%v] data[%X]", rawCANOpenFrame, id, dlc, ext, rtr, canFrame.Data)
	return nil
}

//   - ------- + -------- + ------------+ ------------ + --- +
//   - Pre-defined CAN-IDs for CANopen(CIA301a) protocols    +
//   - ------- + -------- + ------------+ ------------ + --- +
//   - Message + FuncCode + COB-ID(11b) +     Index    + sub +
//   - ------- + -------- + ------------+ ------------ + --- +
//   - NMT     +   000 0b  + 000h       +    N/A       + N/A + OK
//   - SYNC    +   000 1b  + 080h       +  1005h-1007h + 00h + OK
//   - EMCY    +   000 1b  + 081h-0FFh  +  1014h-1015h + 00h + OK
//   - TIME    +   001 0b  + 100h       +  1012h       + 00h + OK
//   - TPDO_1  +   001 1b  + 181h-1FFh  +  1800h       + 01h + OK
//   - RPDO_1  +   010 0b  + 201h-27Fh  +  1400h       + 01h + OK
//   - TPDO_2  +   010 1b  + 281h-2FFh  +  1801h       + 01h + OK
//   - RPDO_2  +   011 0b  + 301h-37Fh  +  1401h       + 01h + OK
//   - TPDO_3  +   011 1b  + 381h-3FFh  +  1802h       + 01h + OK
//   - RPDO_3  +   100 0b  + 401h-47Fh  +  1402h       + 01h + OK
//   - TPDO_4  +   100 1b  + 481h-4FFh  +  1803h       + 01h + OK
//   - RPDO_4  +   100 1b  + 501h-57Fh  +  1403h       + 01h + OK
//   - MPDO_n  +   1xx 1b  + 181h-57Fh  +  0060h       + 01h + OK
//   - TSDO_1  +   101 1b  + 581h-5FFh  +    N/A       + N/A + OK
//   - RSDO_1  +   110 0b  + 601h-57Fh  +    N/A       + N/A + OK
//   - NMT_EC  +   111 0b  + 701h-77Fh  +    N/A       + N/A + OK
//   - ------- + -------- + ------------+ ------------ + --- +
//     TSDO(server-to-client)
//     RSDO(client-to-server)
//     NMT_EC = (NMT Error Control)Boot-up/Heartbeat
func ParseCANopen_CIA301aFrame(canFrame *CANOPEN_Frame) error {
	var cobid uint32 = 0
	var nodeid uint32 = 0
	var controlbyte uint8 = 0
	var index uint16 = 0
	var subindex uint8 = 0
	if len(canFrame.Data) != int(canFrame.DLC) {
		// Infalid CAN Frame
		return fmt.Errorf("ParseCANopen_CIA301aFrame:Invalid CANFrame Len Frame Data Len[%d] <> frame.DLC[%d]", len(canFrame.Data), int(canFrame.DLC))
	}
	if len(canFrame.Data) > 0 {
		controlbyte = canFrame.Data[0]
	}
	if len(canFrame.Data) >= 4 {
		index = uint16(canFrame.Data[1])<<8 + uint16(canFrame.Data[2])
		subindex = canFrame.Data[3]
	}
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3     4 bits BigEndian CANopen Only
	// CIA301a 11 bits
	cobid = (canFrame.ID >> 20) & 0xFFF // b0-10 11 bits BigEndian
	nodeid = (canFrame.ID >> 20) & 0x7F // b4-10  7 bits BigEndian (CIA304a)11 bita
	if functionCode == 1 || functionCode == 2 || functionCode == 4 {
		nodeid -= 0x080
		nodeid &= 0x0FF
	}
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CAN        canFrame.ID[0x%X]ID bo-10", canFrame.ID)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*        COB-ID[0x%X]ID bo-10", cobid)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*        nodeid[0x%X]ID b4-10", nodeid)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*  functionCode[0x%d]ID b0-3s", functionCode)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*   controlbyte[0x%02X]Data[0]    8 bits", controlbyte)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*         index[0x%04X]Data[1-2] 16 bits", index)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*      subindex[0x%02X]Data[3]    8 bits", subindex)
	CANOpenCtlByteInfo.cs = (controlbyte >> 5) & 0x07 // controlByte 3 bits b5-7
	CANOpenCtlByteInfo.x = (controlbyte >> 4) & 0x01  // controlByte 1 bit  b4
	CANOpenCtlByteInfo.n = (controlbyte >> 2) & 0x03  // controlByte 2 bits b2-3
	CANOpenCtlByteInfo.s = (controlbyte >> 1) & 0x01  // controlByte 1 bit  b1
	CANOpenCtlByteInfo.e = (controlbyte & 0x01)       // controlByte 1 bits b0
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*            cs[%d] controlByte 3 bits b5-7\n", CANOpenCtlByteInfo.cs)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*             x[%d] controlByte 1 bit  b4\n", CANOpenCtlByteInfo.x)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*             n[%d] controlByte 2 bits b2-3\n", CANOpenCtlByteInfo.n)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*             s[%d] controlByte 1 bits b1\n", CANOpenCtlByteInfo.s)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen*             e[%d] controlByte 1 bits b0\n", CANOpenCtlByteInfo.e)
	if cobid == 0x000 { // + NMT     +   000 0b  + 000h      +    N/A       + N/A +
		logPkg.CtsLog.Debug("CANopen:NMT Frame OK")
	} else if cobid == 0x080 { // + SYNC    +   000 1b  + 080h      +  1005h-1007h + 00h +
		if index >= 0x1005 && index <= 0x1007 {
			// SYNC Frame OK
			logPkg.CtsLog.Debug("CANopen:SYNC Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen:SYNC Frame ER: expected cobid[0x%04X]exp=[0x80] FunctionCode[%d] & index[0x%04X]exp=[0x1005-0x1007]", cobid, functionCode, index)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen:SYNC Frame ER: expected cobid[0x%04X]exp=[0x80] FunctionCode[%d] & index[0x%04X]exp=[0x1005-0x1007]", cobid, functionCode, index)
		}
	} else if cobid >= 0x081 && cobid <= 0x0FF { // + EMCY    +   000 1b  + 081h-0FFh +  1014h-1015h + 00h +
		if index >= 0x1014 && index <= 0x1014 && subindex == 0x00 {
			// EMCY Frame OK
			logPkg.CtsLog.Debug("CANopen:EMCY Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen:EMCY Frame ER: expected cobid[0x%04X]exp=[0x080-0xFF] FunctionCode[%d] & index[0x%04X]exp=[0x1014-0x1014] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen:EMCY Frame ER: expected cobid[0x%04X]exp=[0x080-0xFF] FunctionCode[%d] & index[0x%04X]exp=[0x1014-0x1014] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
		}
	} else if cobid == 0x100 { // + TIME    +   001 0b  + 100h      +  1012h       + 00h +
		if index == 0x1012 && subindex == 0x00 {
			// TIME Frame OK
			logPkg.CtsLog.Debug("CANopen:TIME Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen:TIME Frame ER: expected cobid[0x%04X]exp=[0x100] FunctionCode[%d] & index[0x%04X]exp=[0x1012] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen:TIME Frame ER: expected cobid[0x%04X]exp=[0x100] FunctionCode[%d] & index[0x%04X]exp=[0x1012] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
		}
	} else if cobid >= 0x181 && cobid <= 0x57F { // PDO    081h-57Fh
		objType, objData, err := GetPDO_OjbDicData(index, subindex)
		if err != nil {
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen:PDO Frame ER: cobid[0x%04X]exp=[0x18100000-0x57F] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen:PDO Frame ER: cobid[0x%04X]exp=[0x181-0x57F] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
		}
		logPkg.CtsLog.Debug("ParseCANopen_CIA301aFrame:CANopen:PDO Frame OK: cobid[0x%04X]exp=[0x181-0x57F] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
		if cobid >= 0x181 && cobid <= 0x1FF && index == 0x1800 && subindex == 0x01 {
			// + TPDO_1  +   001 1b  + 181h-1FFh +  1800h       + 01h +
			logPkg.CtsLog.Debug("CANopen:TPDO_1 Frame OK")
		} else if cobid >= 0x201 && cobid <= 0x27F && index == 0x1400 && subindex == 0x01 {
			// + RPDO_1  +   010 0b  + 201h-27Fh +  1400h       + 01h +
			logPkg.CtsLog.Debug("CANopen:RPDO_1 Frame OK")
		} else if cobid >= 0x281 && cobid <= 0x2FF && index == 0x1801 && subindex == 0x01 {
			// + TPDO_2  +   010 1b  + 281h-2FFh +  1801h       + 01h +
			logPkg.CtsLog.Debug("CANopen:TPDO_2 Frame OK")
		} else if cobid >= 0x301 && cobid <= 0x37F && index == 0x1401 && subindex == 0x01 {
			// + RPDO_2  +   011 0b  + 301h-37Fh +  1401h       + 01h +
			logPkg.CtsLog.Debug("CANopen:RPDO_2 Frame OK")
		} else if cobid >= 0x381 && cobid <= 0x3FF && index == 0x1802 && subindex == 0x01 {
			// + TPDO_3  +   011 1b  + 381h-3FFh +  1802h       + 01h +
			logPkg.CtsLog.Debug("CANopen:TPDO_3 Frame OK")
		} else if cobid >= 0x401 && cobid <= 0x47F && index == 0x1402 && subindex == 0x01 {
			// + RPDO_3  +   100 0b  + 401h-47Fh +  1402h       + 01h +
			logPkg.CtsLog.Debug("CANopen:RPDO_3 Frame OK")
		} else if cobid >= 0x481 && cobid <= 0x4FF && index == 0x1803 && subindex == 0x01 {
			// + TPDO_4  +   100 1b  + 481h-4FFh +  1803h       + 01h +
			logPkg.CtsLog.Debug("CANopen:TPDO_4 Frame OK")
		} else if cobid >= 0x501 && cobid <= 0x57F && index == 0x1403 && subindex == 0x01 {
			// + RPDO_4  +   100 1b  + 501h-57Fh +  1403h       + 01h +
			logPkg.CtsLog.Debug("CANopen:RPDO_4 Frame OK")
		} else if cobid >= 0x181 && cobid <= 0x57F && index == 0x0060 && subindex == 0x01 {
			// - MPDO_n  +   1xx 1b  + 181h-57Fh  +  0060h      + 01h + OK
			logPkg.CtsLog.Debug("CANopen:MPDO_n Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen:PDO Frame ER: expected cobid[0x%04X]exp=[0x181-0x57F] FunctionCode[%d] & index[0x%04X]exp=[0x140x-0x180x] & subindex[0x%02X]exp=[0x01]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen:PDO Frame ER: expected cobid[0x%04X]exp=[0x181-0x57F] FunctionCode[%d] & index[0x%04X]exp=[0x140x-0x180x] & subindex[0x%02X]exp=[0x01]", cobid, functionCode, index, subindex)
		}
	} else if cobid >= 0x581 && cobid <= 0x5FF { // + TSDO_1  +   101 1b  + 581h-5FFh +    N/A       + N/A +
		// TSDO_1 Frame OK
		logPkg.CtsLog.Debug("CANopen:TSDO_1 Frame OK")
	} else if cobid >= 0x601 && cobid <= 0x6FF { // + RSDO_1  +   110 0b  + 601h-6FFh +    N/A       + N/A +
		// RSDO_1 Frame OK
		logPkg.CtsLog.Debug("CANopen:RSDO_1 Frame OK")

	} else if cobid >= 0x701 && cobid <= 0x77F { // + NMT_EC  +   111 0b  + 701h-77Fh +    N/A       + N/A +
		// NMT_EC Frame OK
		logPkg.CtsLog.Debug("CANopen:NMT_EC Frame OK")
	} else {
		if index >= 0x0001 && index <= 0x009F {
			// + 0001h-009Fh + Static and Complex Data Types                 +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301aFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Static and Complex Data Types", cobid, functionCode, index, subindex)
		} else if index >= 0x1000 && index <= 0x1FFF {
			// + 1000h-1FFFh + Communication profiles (e.g., DS 301, DS 302) +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301aFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Communication profiles", cobid, functionCode, index, subindex)
		} else if index >= 0x2000 && index <= 0x25FF {
			// + 2000h-25FFh + Manufacturer Specific Device Profile          +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301aFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Manufacturer Specific Device Profile", cobid, functionCode, index, subindex)
		} else if index >= 0x6000 && index <= 0x9FFF {
			// + 6000h-9FFFh + Standardized device profiles                  +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301aFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Standardized device profiles", cobid, functionCode, index, subindex)
		} else {
			// Is Not a CANopen frame
			logPkg.CtsLog.Error("ParseCANopen_CIA301aFrame:CANopen Frame ER: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X] no valid CANopen Frame", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301aFrame:CANopen Frame ER: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X] no valid CANopen Frame", cobid, functionCode, index, subindex)
		}
	}
	return nil
}

//   - ------- + -------- + ---------------------+ ------------ + --- +
//   - Pre-defined CAN-IDs for CANopen(CIA301b) protocols        +
//   - ------- + -------- + -------------------- + ------------ + --- +
//   - Message + FuncCode + COB-ID(29b)          +     Index    + sub +
//   - ------- + -------- + -------------------- + ------------ + --- +
//   - NMT     +   000 0b  + 00000000h           +    N/A       + N/A + OK
//   - SYNC    +   000 1b  + 08000000h           +  1005h-1007h + 00h + OK
//   - EMCY    +   000 1b  + 08100000h-0FFFFFFFh +  1014h-1015h + 00h + OK
//   - TIME    +   001 0b  + 10000000h           +  1012h       + 00h + OK
//   - TPDO_1  +   001 1b  + 18100000h-1FFFFFFFh +  1800h       + 01h + OK
//   - RPDO_1  +   010 0b  + 20100000h-27FFFFFFh +  1400h       + 01h + OK
//   - TPDO_2  +   010 1b  + 28100000h-2FFFFFFFh +  1801h       + 01h + OK
//   - RPDO_2  +   011 0b  + 30100000h-37FFFFFFh +  1401h       + 01h + OK
//   - TPDO_3  +   011 1b  + 38100000h-3FFFFFFFh +  1802h       + 01h + OK
//   - RPDO_3  +   100 0b  + 40100000h-47FFFFFFh +  1402h       + 01h + OK
//   - MPDO_n  +   1xx 1b  + 18100000h-57FFFFFFh +  0060h       + 01h + OK
//   - TPDO_4  +   100 1b  + 48100000h-4FFFFFFFh +  1803h       + 01h + OK
//   - RPDO_4  +   100 1b  + 50100000h-57FFFFFFh +  1403h       + 01h + OK
//   - TSDO_1  +   101 1b  + 58100000h-5FFFFFFFh +    N/A       + N/A + OK
//   - RSDO_1  +   110 0b  + 60100000h-57FFFFFFh +    N/A       + N/A + OK
//   - NMT_EC  +   111 0b  + 70100000h-77FFFFFFh +    N/A       + N/A + OK
//   - ------- + -------- + ---------------------+ ------------ + --- +
//     TSDO(server-to-client)
//     RSDO(client-to-server)
//     NMT_EC = (NMT Error Control)Boot-up/Heartbeat
func ParseCANopen_CIA301bFrame(canFrame *CANOPEN_Frame) error {
	var cobid uint32 = 0
	var nodeid uint32 = 0
	var controlbyte uint8 = 0
	var index uint16 = 0
	var subindex uint8 = 0
	if len(canFrame.Data) != int(canFrame.DLC) {
		// Infalid CAN Frame
		logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:Invalid CANFrame Len Frame Data Len[%d] <> frame.DLC[%d]", len(canFrame.Data), int(canFrame.DLC))
		return fmt.Errorf("ParseCANopen_CIA301bFrame:Invalid CANFrame Len Frame Data Len[%d] <> frame.DLC[%d]", len(canFrame.Data), int(canFrame.DLC))
	}
	if len(canFrame.Data) > 0 {
		controlbyte = canFrame.Data[0]
	}
	if len(canFrame.Data) >= 4 {
		index = uint16(canFrame.Data[1])<<8 + uint16(canFrame.Data[2])
		subindex = canFrame.Data[3]
	}
	functionCode := (canFrame.ID >> 28) & 0x0F // b0-3     4 bits BigEndian CANopen Only
	// CIA301b 29 bits
	cobid = (canFrame.ID) & 0xFFFFFFFF  // b0-28 29 bits BigEndian
	nodeid = (canFrame.ID) & 0x0FFFFFFF // b4-28 25 bits BigEndian (CIA304b)29 bita
	if functionCode == 1 || functionCode == 2 || functionCode == 4 {
		nodeid -= 0x08000000
		nodeid &= 0x0FFFFFFF
	}
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CAN        canFrame.ID[0x%X]ID bo-28", canFrame.ID)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*        COB-ID[0x%X]ID bo-28", cobid)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*        nodeid[0x%X]ID b4-28", nodeid)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*  functionCode[0x%d]ID b0-3      4s bits", functionCode)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*   controlbyte[0x%02X]Data[0]    8 bits", controlbyte)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*         index[0x%04X]Data[1-2] 16 bits", index)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*      subindex[0x%02X]Data[3]    8 bits", subindex)
	CANOpenCtlByteInfo.cs = (controlbyte >> 5) & 0x07 // controlByte 3 bits b5-7
	CANOpenCtlByteInfo.x = (controlbyte >> 4) & 0x01  // controlByte 1 bit  b4
	CANOpenCtlByteInfo.n = (controlbyte >> 2) & 0x03  // controlByte 2 bits b2-3
	CANOpenCtlByteInfo.s = (controlbyte >> 1) & 0x01  // controlByte 1 bit  b1
	CANOpenCtlByteInfo.e = (controlbyte & 0x01)       // controlByte 1 bits b0
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*            cs[%d] controlByte 3 bits b5-7\n", CANOpenCtlByteInfo.cs)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*             x[%d] controlByte 1 bit  b4\n", CANOpenCtlByteInfo.x)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*             n[%d] controlByte 2 bits b2-3\n", CANOpenCtlByteInfo.n)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*             s[%d] controlByte 1 bits b1\n", CANOpenCtlByteInfo.s)
	logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen*             e[%d] controlByte 1 bits b0\n", CANOpenCtlByteInfo.e)
	if cobid == 0x00000000 { // + NMT     +   000 0b  + 000h      +    N/A       + N/A +
		logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:NMT Frame OK")
	} else if cobid == 0x08000000 { // + SYNC    +   000 1b  + 080h      +  1005h-1007h + 00h +
		if index >= 0x1005 && index <= 0x1007 {
			// SYNC Frame OK
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:SYNC Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:CANopen:SYNC Frame ER: expected cobid[0x%04X]exp=[0x80] FunctionCode[%d] & index[0x%04X]exp=[0x1005-0x1007]", cobid, functionCode, index)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen:SYNC Frame ER: expected cobid[0x%04X]exp=[0x80] FunctionCode[%d] & index[0x%04X]exp=[0x1005-0x1007]", cobid, functionCode, index)
		}
	} else if cobid >= 0x08100000 && cobid <= 0x0FFFFFFF { // + EMCY    +   000 1b  + 081h-0FFh +  1014h-1015h + 00h +
		if index >= 0x1014 && index <= 0x1014 && subindex == 0x00 {
			// EMCY Frame OK
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:EMCY Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:CANopen:EMCY Frame ER: expected cobid[0x%04X]exp=[0x080-0xFF] FunctionCode[%d] & index[0x%04X]exp=[0x1014-0x1014] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen:EMCY Frame ER: expected cobid[0x%04X]exp=[0x080-0xFF] FunctionCode[%d] & index[0x%04X]exp=[0x1014-0x1014] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
		}
	} else if cobid == 0x10000000 { // + TIME    +   001 0b  + 100h      +  1012h       + 00h +
		if index == 0x1012 && subindex == 0x00 {
			// TIME Frame OK
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:TIME Frame OK")
		} else {
			logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:CANopen:TIME Frame ER: expected cobid[0x%04X]exp=[0x100] FunctionCode[%d] & index[0x%04X]exp=[0x1012] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen:TIME Frame ER: expected cobid[0x%04X]exp=[0x100] FunctionCode[%d] & index[0x%04X]exp=[0x1012] & subindex[0x%02X]exp=[0x00]", cobid, functionCode, index, subindex)
		}
	} else if cobid >= 0x18100000 && cobid <= 0x57FFFFFF { // PDO    18100000h-57FFFFFFh
		objType, objData, err := GetPDO_OjbDicData(index, subindex)
		if err != nil {
			logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:CANopen:PDO Frame ER: cobid[0x%04X]exp=[0x18100000-0x57FFFFFF] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen:PDO Frame ER: cobid[0x%04X]exp=[0x18100000-0x57FFFFFF] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
		}
		logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:PDO Frame OK: cobid[0x%04X]exp=[0x18100000-0x57FFFFFF] FunctionCode[%d] & index[0x%04X] subindex[0x%02X] objType[%d] objData[%0X]", cobid, functionCode, index, subindex, objType, objData)
		if cobid >= 0x18100000 && cobid <= 0x1FFFFFFF && index >= 0x1800 && subindex == 0x01 {
			// + TPDO_1  +   001 1b  + 18100000h-1FFFFFFFh +  1800h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:TPDO_1 Frame OK")
		} else if cobid >= 0x20100000 && cobid <= 0x27FFFFFF && index >= 0x1400 && subindex == 0x01 {
			// + RPDO_1  +   010 0b  + 20100000h-27FFFFFFh +  1400h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:RPDO_1 Frame OK")
		} else if cobid >= 0x22810000 && cobid <= 0x2FFFFFFF && index >= 0x1801 && subindex == 0x01 {
			// + TPDO_2  +   010 1b  + 28100000h-2FFFFFFFh +  1801h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:TPDO_2 Frame OK")
		} else if cobid >= 0x30100000 && cobid <= 0x37FFFFFF && index >= 0x1401 && subindex == 0x01 {
			// + RPDO_2  +   011 0b  + 30100000h-37FFFFFFh +  1401h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:RPDO_2 Frame OK")
		} else if cobid >= 0x38100000 && cobid <= 0x3FFFFFFF && index >= 0x1802 && subindex == 0x01 {
			// + TPDO_3  +   011 1b  + 38100000h-3FFFFFFFh +  1802h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:TPDO_3 Frame OK")
		} else if cobid >= 0x40100000 && cobid <= 0x47FFFFFF && index >= 0x1402 && subindex == 0x01 {
			// + RPDO_3  +   100 0b  + 40100000h-47FFFFFFh +  1402h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:RPDO_3 Frame OK")
		} else if cobid >= 0x48100000 && cobid <= 0x4FFFFFFF && index >= 0x1803 && subindex == 0x01 {
			// + TPDO_4  +   100 1b  + 48100000h-4FFFFFFFh +  1803h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:TPDO_4 Frame OK")
		} else if cobid >= 0x50100000 && cobid <= 0x57FFFFFF && index >= 0x1403 && subindex == 0x01 {
			// + RPDO_4  +   100 1b  + 50100000h-57FFFFFFh +  1403h       + 01h +
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:RPDO_4 Frame OK")
		} else if cobid >= 0x18100000 && cobid <= 0x57FFFFFF && index == 0x0060 && subindex == 0x01 {
			// - MPDO_n  +   1xx 1b  + 18100000h-57FFFFFFh  +  0060h      + 01h + OK
			logPkg.CtsLog.Debug("CANopen:MPDO_n Frame OK")
		} else {
			logPkg.CtsLog.Debug("ParseCANopen_CIA301bFrame:CANopen:PDO Frame ER: expected cobid[0x%04X]exp=[0x18100000-0x57FFFFFF] FunctionCode[%d] & index[0x%04X]exp=[0x140x-0x180x] & subindex[0x%02X]exp=[0x01]", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen:PDO Frame ER: expected cobid[0x%04X]exp=[0x18100000-0x57FFFFFF] FunctionCode[%d] & index[0x%04X]exp=[0x140x-0x180x] & subindex[0x%02X]exp=[0x01]", cobid, functionCode, index, subindex)
		}
	} else if cobid >= 0x58100000 && cobid <= 0x5FFFFFFF { // + TSDO_1  +   101 1b  + 581h-5FFh +    N/A       + N/A +
		// TSDO_1 Frame OK
		logPkg.CtsLog.Debug("CANopen:TSDO_1 Frame OK")
	} else if cobid >= 0x60100000 && cobid <= 0x6FFFFFFF { // + RSDO_1  +   110 0b  + 601h-6FFh +    N/A       + N/A +
		// RSDO_1 Frame OK
		logPkg.CtsLog.Debug("CANopen:RSDO_1 Frame OK")
	} else if cobid >= 0x70100000 && cobid <= 0x77FFFFFF { // + NMT_EC  +   111 0b  + 701h-77Fh +    N/A       + N/A +
		// NMT_EC Frame OK
		logPkg.CtsLog.Debug("CANopen:NMT_EC Frame OK")
	} else {
		if index >= 0x0001 && index <= 0x009F {
			// + 0001h-009Fh + Static and Complex Data Types                 +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301bFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Static and Complex Data Types", cobid, functionCode, index, subindex)
		} else if index >= 0x1000 && index <= 0x1FFF {
			// + 1000h-1FFFh + Communication profiles (e.g., DS 301, DS 302) +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301bFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Communication profiles", cobid, functionCode, index, subindex)
		} else if index >= 0x2000 && index <= 0x25FF {
			// + 2000h-25FFh + Manufacturer Specific Device Profile          +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301bFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Manufacturer Specific Device Profile", cobid, functionCode, index, subindex)
		} else if index >= 0x6000 && index <= 0x9FFF {
			// + 6000h-9FFFh + Standardized device profiles                  +
			logPkg.CtsLog.Warn("ParseCANopen_CIA301bFrame:CANopen Frame OK: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X]Standardized device profiles", cobid, functionCode, index, subindex)
		} else {
			// Is Not a CANopen frame
			logPkg.CtsLog.Error("ParseCANopen_CIA301bFrame:CANopen Frame ER: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X] no valid CANopen Frame", cobid, functionCode, index, subindex)
			return fmt.Errorf("ParseCANopen_CIA301bFrame:CANopen Frame ER: cobid[0x%04X] FunctionCode[%d] index[0x%04X] subindex[0x%02X] no valid CANopen Frame", cobid, functionCode, index, subindex)
		}
	}
	return nil
}

// Parse CANopen PDO ObjectDictionary Index+subindex params
// If found return the data and it type U8, U16 or U32
// R1 = Data_TypeCode 8=U8, 16=U16, 32=U32 all other are invalid
// R2 = Data stored on PDO Object Dictionary
// R3 = nil or error code
func GetPDO_OjbDicData(index uint16, subindex uint8) (uint8, uint32, error) {
	var pdoDataType uint8 = 0
	var pdoData uint32 = 0
	var err error = nil
	if index >= 0x1400 && index <= 0x147F {
		if subindex == 0x00 {
			// + 1400h-147Fh + Receive PDO Param    + 00h       + Largest Subindex Supported  + U32  + RO     +
			pdoDataType = 32
			pdoData = 8
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Param Largest Subindex Supported", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x01 {
			// +                                    + 01h       + COB-ID used in PDO          + U32  + RW     +
			pdoDataType = 32
			pdoData = 0
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Param COB-ID used in PDO", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x02 {
			// +                                    + 02h       + Transmission type           + U8   + RW     +
			pdoDataType = 8
			pdoData = 0 // TODO: Rreturn a Transmit Type
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Param Transmission type", index, subindex, pdoDataType, pdoData)
		} else {
			logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Param Expected sub-Index(00-02h) ", index, subindex, pdoDataType, pdoData)
			err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Param Expected sub-Index(00-02h) ", index, subindex, pdoDataType, pdoData)
		}
	} else if index >= 0x1600 && index <= 0x167F {
		// + 1600h-167Fh + Receive PDO Mapping  + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
		if subindex == 0x00 {
			// + 1600h-167Fh + Receive PDO Mapping  + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
			pdoDataType = 8
			pdoData = 0
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping  NumberOf Mapped Appl PDO", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x01 {
			// +                                    + 01h       + Mapped Object #1            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #1
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #1", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x02 {
			// +                                    + 02h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #2
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #2", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x03 {
			// +                                    + 03h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #3
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #3", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x04 {
			// +                                    + 04h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #4
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #4", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x05 {
			// +                                    + 05h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #5
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #5", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x06 {
			// +                                    + 06h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #6
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #6", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x07 {
			// +                                    + 07h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #7
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #7", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x08 {
			// +                                    + 08h       + Mapped Object #8            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #8
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Mapped Object #8", index, subindex, pdoDataType, pdoData)
		} else {
			logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Expected sub-Index(00-08h) ", index, subindex, pdoDataType, pdoData)
			err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Receive PDO Mapping Expected sub-Index(00-08h) ", index, subindex, pdoDataType, pdoData)
		}
	} else if index >= 0x1800 && index <= 0x187F {
		if subindex == 0x00 {
			// + 1800h-187Fh + Transmit PDO Param    + 00h       + Largest Subindex Supported  + U32  + RO     +
			pdoDataType = 32
			pdoData = 8
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Largest Subindex Supported", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x01 {
			// +                                    + 01h       + COB-ID used in PDO          + U32  + RW     +
			pdoDataType = 32
			pdoData = 0
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param COB-ID used in PDO", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x02 {
			// +                                    + 02h       + Transmission type           + U8   + RW     +
			pdoDataType = 8
			pdoData = 0 // TODO: Rreturn a Transmit Type
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Transmission type", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x03 {
			// +                                    + 02h       + Transmission type           + U8   + RW     +
			pdoDataType = 8
			pdoData = 0 // TODO: Rreturn a Transmit Type
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Transmission type", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x05 {
			// +                                    + 02h       + Transmission type           + U8   + RW     +
			pdoDataType = 8
			pdoData = 0 // TODO: Rreturn a Transmit Type
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Transmission type", index, subindex, pdoDataType, pdoData)
		} else {
			logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Expected sub-Index(00-02h) ", index, subindex, pdoDataType, pdoData)
			err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Param Expected sub-Index(00-02h) ", index, subindex, pdoDataType, pdoData)
		}
	} else if index >= 0x1A00 && index <= 0x1A7F {
		// + 1A00h-1A7Fh + Transmit PDO Mapping  + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
		if subindex == 0x00 {
			// + 1A00h-1A7Fh + Transmit PDO Mapping  + 00h       + NumberOf Mapped Appl PDO    + U8   + RW     +
			pdoDataType = 8
			pdoData = 0
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping  NumberOf Mapped Appl PDO", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x01 {
			// +                                    + 01h       + Mapped Object #1            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #1
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #1", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x02 {
			// +                                    + 02h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #2
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #2", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x03 {
			// +                                    + 03h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #3
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #3", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x04 {
			// +                                    + 04h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #4
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #4", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x05 {
			// +                                    + 05h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #5
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #5", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x06 {
			// +                                    + 06h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #6
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #6", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x07 {
			// +                                    + 07h       + Mapped Object #2            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #7
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #7", index, subindex, pdoDataType, pdoData)
		} else if subindex == 0x08 {
			// +                                    + 08h       + Mapped Object #8            + U32  + RW     +
			pdoDataType = 32
			pdoData = 0 // TODO: Rreturn Mapped Object #8
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Mapped Object #8", index, subindex, pdoDataType, pdoData)
		} else {
			logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Expected sub-Index(00-08h) ", index, subindex, pdoDataType, pdoData)
			err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Transmit PDO Mapping Expected sub-Index(00-08h) ", index, subindex, pdoDataType, pdoData)
		}
	} else if index == 0x0060 { // MPDO
		if subindex == 0x01 {
			// + 0060h      +       MPDO Mapping  + 01h       + MPDO_n    + U8   + RW     +
			pdoDataType = 8
			pdoData = 0
			logPkg.CtsLog.Debug("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] MPDO_n", index, subindex, pdoDataType, pdoData)
		} else {
			logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] MPDO_n Error", index, subindex, pdoDataType, pdoData)
			err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] MPDO_n Error", index, subindex, pdoDataType, pdoData)
		}
	} else {
		logPkg.CtsLog.Error("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Unknown", index, subindex, pdoDataType, pdoData)
		err = fmt.Errorf("GetPDO_OjbDicData:            Index[0x%04X] Sub-Index[0x%02X] pdoDataType[%02d] pdoDa0xta[0x%08X] Unknown", index, subindex, pdoDataType, pdoData)
	}
	return pdoDataType, pdoData, err
}

// Convert a CAN Frame to Byte array
func ConvertCANOPEN_FameToBytes(inCanFrame *CANOPEN_Frame) ([]byte, error) {
	// Determine the length of the CAN frame
	frameLength := len(inCanFrame.Data) + 5 // CanFame id(4) DLC(1)

	// CANopenTcpServer_Create a byte slice with the appropriate length
	outRawCanFrame := make([]byte, frameLength)
	outRawCanFrame[0] = byte(inCanFrame.ID >> 24)
	outRawCanFrame[1] = byte(inCanFrame.ID >> 16)
	outRawCanFrame[2] = byte(inCanFrame.ID >> 8)
	outRawCanFrame[3] = byte(inCanFrame.ID)
	outRawCanFrame[4] = byte(inCanFrame.DLC)
	for i := 0; i < len(inCanFrame.Data); i++ {
		outRawCanFrame[i+5] = inCanFrame.Data[i]
	}
	return outRawCanFrame, nil
}
