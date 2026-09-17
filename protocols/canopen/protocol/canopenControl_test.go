package canopen

import (
	utilsPkg "canopen/utils"
	logPkg "canopen/utils/gologtofile"
	"encoding/hex"
	"os"
	"runtime"
	"testing"
)

var LogStarted bool = false

func Test001_CANopenControl_ParseRawCANopenFrame_Null_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test001 - CANopen Frame = null
	//                 <F<--ID->fl<--Data---->
	//                  0 1 2 3 4 5 6 7 8 9 0
	err := ParseRawCANopenFrame(nil)
	if err == nil {
		t.Errorf("FAIL: NULL CAN Frame Expected error received err[%v]", err)
	}
	logPkg.CtsLog.Debug("TEST01:PASS: Null CAN Frame OK\n\n")
}

func Test002_CANopenControl_ParseRawCANopenFrame_Short_rawCanOpenFrame(t *testing.T) {
	StartLog()
	// Test002 - CAN Frame = Short
	//                <F<--ID->fl<--Data---->
	//                  0 1 2 3 4 5 6 7 8 9 0
	strRawCANOpeFrame := "600123"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANOpenFrame, err := hex.DecodeString(strRawCANOpeFrame)
	if err != nil {
		logPkg.CtsLog.Error("Test02:conversion of Hex str to hex []byte(%s)fail", strRawCANOpeFrame)
		t.Errorf("Test02:conversion of Hex str to hex []byte(%s)fail", strRawCANOpeFrame)
	}
	//Run Test now
	err = ParseRawCANopenFrame(rawCANOpenFrame)
	if err == nil {
		t.Errorf("FAIL:rawCANOpenFrame[%X] Short CAN Frame Expected error received err[%v]", rawCANOpenFrame, err)
		return
	}
	logPkg.CtsLog.Debug("TEST02:PASS: Short CAN Frame OK\n\n")
}

func Test003_CANopenControl_ParseRawCANopenFrame_rawCANopenFrame_DLC_Missnatch_Data_Size(t *testing.T) {
	StartLog()
	// Test003 - CAN Frame = Mismatch Data Size
	//                     <F<--ID->fl<--Data---->
	//                      0 1 2 3 4 5 6 7 8 9 0
	strRawCANopenFrame := "6001234585010203040506"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANOpenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("Test03:conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("Test03:conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	err = ParseRawCANopenFrame(rawCANOpenFrame)
	if err == nil {
		t.Errorf("FAIL: CAN Frame Expected DLC Missmatch Data Size err[%v] rawCANOpenFrame[%v]", err, rawCANOpenFrame)
		return
	}
	logPkg.CtsLog.Debug("TEST03:PASS: CAN Frrame DLC Missnath  Data Size OK\n\n")
}

func Test04_CANopenControl_ParseRawCANopenFrame_Long_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test004 - CANopen Frame = Long
	//                <F<--ID->fl<--Data[l]------->
	//                            0 1 2 3 4 5 6 7 x << Data
	//                  0 1 2 3 4 5 6 7 8 9 0 1 2 3 << Full CANFrame
	strRawCANopenFrame := "6001234589010203040506070809"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("Test04:conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("Test04:conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err == nil {
		t.Errorf("FAIL: Big CAN Frame Expected error received err[%v] CANRawFrame[%X]", err, rawCANopenFrame)
		return
	}
	logPkg.CtsLog.Debug("TEST04:PASS: Long  CANopen Frame OK\n\n")
}

func Test005_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RSDO_Node0x21_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                            0 1 2 3 4 5 6 7 << Data
	//                  0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "62112345080102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST05: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST05: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST05: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST05: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST05: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test006_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RSDO_Node0x21_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "62112345880102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST06: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST05: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST06: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST06: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST06: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test007_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TSDO_Node0x21_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test007 - CANopen Frame = Good
	//                    <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "5A112345080102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST07: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST07: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST07: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST07: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST07: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test008_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TSDO_Node0x21_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test008 - CANopen Frame = Good
	//                    <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "5A112345880102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST08: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST08: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST08: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST08: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST08: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test009_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RSDO_Node0x3A_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "63A12345080102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST09: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST09: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST09: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST09: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST09: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test010_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RSDO_Node0x3A_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "63A12345880102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST10: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST10: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST10: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST10: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST10: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test011_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TSDO_Node0x3A_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "5BA12345080102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST11: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST11: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST11: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST11: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST11: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test012_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TSDO_Node0x3A_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test005 - CANopen Frame = Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	strRawCANopenFrame := "5BA12345880102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST12: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST10: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST12: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST12: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST12: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test013_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TPDO_1_Node0x3A_Idx1800_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test013 - CANopen Frame = TPDO_1 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_1  +   001 1b  + 181h-1FFh +  1800h       + 01h +
	strRawCANopenFrame := "1BA12345080118000105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST13: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST13: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST13: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST13: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST13: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test014_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TPDO_1_Node0x3A_Idx1800_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test013 - CANopen Frame = TPDO_1 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_1  +   001 1b  + 18100000h-1FFFFFFFh +  1800h       + 01h +
	strRawCANopenFrame := "1BA12345880118000105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST13: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST13: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST13: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST13: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST13: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test015_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RPDO_1_Node0x3A_Idx1400_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test013 - CANopen Frame = RPDO_1 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_1  +   001 1b  + 200h-27Fh +  1400h       + 01h +
	strRawCANopenFrame := "23A12345080114000105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST15: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST15: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST15: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST15: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST15: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test016_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RPDO_1_Node0x3A_Idx1400_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test013 - CANopen Frame = RPDO_1 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_1  +   001 1b  + 20000000h-27FFFFFFh +  1400h       + 01h +
	strRawCANopenFrame := "23A12345880114000105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST16: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST16: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST16: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST16: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST16: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test017_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TPDO_2_Node0x3A_Idx1801_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test017 - CANopen Frame = TPDO_2 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_2  +   010 1b  + 281h-2FFh +  1801h       + 01h +
	strRawCANopenFrame := "2BA12345080118010105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST17: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST17: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST17: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST17: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST17: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test018_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TPDO_2_Node0x3A_Idx1801_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test017 - CANopen Frame = TPDO_2 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_2  +   010 1b  + 2810000h-2FFFFFFh +  1801h       + 01h +
	strRawCANopenFrame := "2BA12345880118010105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST17: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST17: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST17: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST17: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST17: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test019_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RPDO_2_Node0x3A_Idx1401_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test017 - CANopen Frame = RPDO_2 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_2  +   011 0b  + 301h-37Fh +  1401h       + 01h +
	strRawCANopenFrame := "33A12345080114010105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST19: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST19: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST19: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST19: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST19: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test020_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RPDO_2_Node0x3A_Idx1401_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test017 - CANopen Frame = RPDO_2 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_2  +   011 0b  + 30100000h-37FFFFFFh +  1401h       + 01h +
	strRawCANopenFrame := "33A12345880114010105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST20: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST20: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST20: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST20: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST20: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test021_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TPDO_3_Node0x3A_Idx1802_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test021 - CANopen Frame = TPDO_3 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_3  +   011 1b  + 381h-3FFh +  1802h       + 01h +
	strRawCANopenFrame := "3BA12345080118020105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST21: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST21: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST21: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST21: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST21: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test022_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TPDO_3_Node0x3A_Idx1802_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test022 - CANopen Frame = TPDO_3 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_3  +   011 1b  + 3810000h-3FFFFFFh +  1802h       + 01h +
	strRawCANopenFrame := "3BA12345880118020105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST22: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST22: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST22: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST22: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST22: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test023_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RPDO_3_Node0x3A_Idx1402_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test023 - CANopen Frame = RPDO_3 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_3  +   100 0b  + 401h-47Fh +  1402h       + 01h +
	strRawCANopenFrame := "43A12345080114020105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST23: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST23: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST23: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST23: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST23: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test024_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RPDO_3_Node0x3A_Idx1402_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test023 - CANopen Frame = RPDO_3 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	//   - RPDO_3  +   100 0b  + 40100000h-47FFFFFFh +  1402h       + 01h + OK
	strRawCANopenFrame := "43A12345880114020105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST23: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST23: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST23: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST23: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST23: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test025_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_TPDO_4_Node0x3A_Idx1802_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test025 - CANopen Frame = TPDO_4 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_4  +   100 1b  + 481h-4FFh +  1803h       + 01h +
	strRawCANopenFrame := "4BA12345080118030105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST25: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST25: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST25: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST25: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST25: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test026_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_TPDO_4_Node0x3A_Idx1802_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test025 - CANopen Frame = TPDO_4 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + TPDO_4  +   100 1b  + 481h-4FFh +  1803h       + 01h +
	strRawCANopenFrame := "4BA12345880118030105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST25: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST25: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST25: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST25: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST25: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test027_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_RPDO_4_Node0x3A_Idx1403_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test027 - CANopen Frame = RPDO_4 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_4  +   100 1b  + 501h-57Fh +  1403h       + 01h +
	strRawCANopenFrame := "53A12345080114030105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST27: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST27: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST27: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST27: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST27: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test028_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_RPDO_4_Node0x3A_Idx1403_sub_01_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test028 - CANopen Frame = RPDO_4 Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + RPDO_4  +   100 1b  + 501h-57Fh +  1403h       + 01h +
	strRawCANopenFrame := "53A12345880114030105060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST28: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST28: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST28: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST28: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST28: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test029_CANopenControl_ParseCANRawFrame_Good_MPDO_Node_CIA301a_0x3A_Idx0060_sub_01_CANRawFrame(t *testing.T) {
	StartLog()
	// ==========================================================
	//  MPDO example frame from Producer to consumer address 0x3A
	// ==========================================================
	// CAN Identifier: 0x400 + Node ID (0x3A) = 0x43A
	// Data Length Code (DLC): 8 bytes
	// Data Byte 1: 0x40 (MPDO COB-ID for a request)
	// Data Byte 2: Object Index (e.g., 0x60 for a specific object)
	// Data Byte 3: Sub-Index (e.g., 0x01 for a specific sub-index)
	// Data Byte 4 to 8: Data specific to the request
	// Test029 - CANopen Frame = Good
	//                    <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + MPDO_1   +   1xx 0b  + 181h-47Fh +    0060       + 01 +
	strRawCANopenFrame := "43A12345084000600105060708"

	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST29: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST29: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST29: CANRawFrame[0x%X]Hex", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST29: FAIL Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST29: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test030_CANopenControl_ParseCANRawFrame_Good_MPDO_Node_CIA301b_0x3A_Idx0060_sub_01_CANRawFrame(t *testing.T) {
	StartLog()
	// ===============================================================
	//  MPDO example frame from Producer to consumer address 0x3A00000
	// ===============================================================
	// CAN Identifier: 0x40000000 + Node ID (0x3A00000) = 0x43A00000
	// Data Length Code (DLC): 8 bytes
	// Data Byte 1: 0x40 (MPDO COB-ID for a request)
	// Data Byte 2: Object Index (e.g., 0x60 for a specific object)
	// Data Byte 3: Sub-Index (e.g., 0x01 for a specific sub-index)
	// Data Byte 4 to 8: Data specific to the request
	// Test030 - CANopen Frame = Good
	//                    <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + MPDO_1   +   1xx 0b  + 18100000h-47F00000h +    0060       + 01 +
	strRawCANopenFrame := "43A12345884000600105060708"

	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST30: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST30: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST30: CANRawFrame[0x%X]Hex", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST30: FAIL Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST30: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test031_CANopenControl_ParseRawCANopenFrame_Good_CIA301a_NMT_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test029 - CANopen Frame = NMT Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + NMT_EC  +   111 0b  + 701h-77Fh +    N/A       + N/A +
	strRawCANopenFrame := "70112345080102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST31: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST31: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST31: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST31: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST31: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func Test032_CANopenControl_ParseRawCANopenFrame_Good_CIA301b_NMT_rawCANopenFrame(t *testing.T) {
	StartLog()
	// Test029 - CANopen Frame = NMT Good
	//                <F<--ID->fl<--Data[l]----->
	//                               0 1 2 3 4 5 6 7 << Data
	//                     0 1 2 3 4 5 6 7 8 9 0 1 2 << Full CANFrame
	// + NMT_EC  +   111 0b  + 701h-77Fh +    N/A       + N/A +
	strRawCANopenFrame := "70112345880102030405060708"
	// Convert strCANRawFrame to []byte CANRawFrame
	rawCANopenFrame, err := hex.DecodeString(strRawCANopenFrame)
	if err != nil {
		logPkg.CtsLog.Error("TEST31: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
		t.Errorf("TEST31: conversion of Hex str to hex []byte(%s)fail", strRawCANopenFrame)
	}
	//Run Test now

	logPkg.CtsLog.Debug("TEST31: CANRawFrame[0x%X]HexOn Entry", rawCANopenFrame)
	err = ParseRawCANopenFrame(rawCANopenFrame)
	if err != nil {
		t.Errorf("TEST31: FAIL: Big CAN Frame Expected good received err[%v]", err)
		return
	}
	logPkg.CtsLog.Debug("TEST31: PASS:OK rawCANopenFrame[%X]\n\n", rawCANopenFrame)
}

func StartLog() {
	if !LogStarted {
		utilsPkg.Hostname, _ = os.Hostname()
		// Initialize CtsLogs with default parameters
		rc, err := logPkg.InitCtsLogs(
			utilsPkg.LogFilePath,
			"emucanopentcp_test",
			utilsPkg.DeleteExistingLogFiles,
			utilsPkg.LogLevel,
			utilsPkg.LogFileMaxSize,
		)
		if (err != nil) || (rc < 0) {
			/// Fail to setup Logs
			logPkg.CtsLog.Error("StartLog:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n         err[%s]\n\n",
				rc,
				runtime.GOOS,
				runtime.GOARCH,
				utilsPkg.OsBits,
				utilsPkg.Hostname,
				utilsPkg.LogFilePath,
				utilsPkg.LogFileName,
				utilsPkg.LogLevel,
				err)
			return
		}
		logPkg.CtsLog.Info("StartLog:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n\n",
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.LogFilePath,
			utilsPkg.LogFileName,
			utilsPkg.LogLevel)
	}
}
