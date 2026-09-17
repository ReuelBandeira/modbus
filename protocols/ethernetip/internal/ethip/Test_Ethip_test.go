package ethip

import (
	"testing"
)

func Test_Ethip_ConnEth(t *testing.T) {
	adrress := "1.1.1.1"
	port := "10"
	resp := isTCPConnected(adrress, port)
	if resp {
		t.Errorf("FAIL: TCP Connected Expected false received")
	}
	t.Logf("TEST:PASS: TCP Connected OK! Adress[%v] Port[%v]", adrress, port)
}

func Test_Ethip_ConnEthernetIP(t *testing.T) {
	adrress := "1.1.1.1"
	port := "10"
	_, err := ConnEthernetIP(adrress, port)
	if err == nil {
		t.Errorf("FAIL: ResolveTCPAddress expected false received")
	} else {
		t.Logf("TEST:PASS: ResolveTCPAddress connected OK! Adress[%v] Port[%v]", adrress, port)
	}
}
func Test_Ethip_ReadBitMemory(t *testing.T) {

	var class uint32 = 100
	var instance uint32 = 10
	var attribute uint32 = 1

	p, err := ConnEthernetIP("172.16.16.61", "44818")
	if err != nil {
		t.Errorf("FAIL: TCP not connected. Expected connected")
		return
	}
	readBitMemory, status := p.ReadBitMemory(class, instance, attribute)
	if status == 0 {
		t.Errorf("FAIL: Read bit memory expected reading fail")
	} else {
		t.Logf("TEST01:PASS: Read bit memory OK! BitMemory [%v] Class[%v] Instance[%v] Attribute[%v]", readBitMemory, class, instance, attribute)
	}
}

func Test_Ethip_ReadWordMemory(t *testing.T) {

	var class uint32 = 100
	var instance uint32 = 10
	var attribute uint32 = 1

	p, err := ConnEthernetIP("172.16.16.61", "44818")
	if err != nil {
		t.Errorf("FAIL: TCP not connected. Expected connected")
		return
	}
	readWordMemory, status := p.ReadWordMemory(class, instance, attribute)
	if status == 0 {
		t.Errorf("FAIL: Read word memory expected reading fail")
	} else {
		t.Logf("TEST01:PASS: Read bit memory OK! WordMemory [%v] Class[%v] Instance[%v] Attribute[%v]", readWordMemory, class, instance, attribute)
	}
}
