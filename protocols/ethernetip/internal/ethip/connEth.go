package ethip

import (
	utilsPkg "ethernetip/utils"
	golog "ethernetip/utils/gologtofile"
)

// ConnEth make a connection with a EthernetIp Device
func ConnEthernetIp(device utilsPkg.DevSettings) bool {

	if !isTCPConnected(device.Address, device.Port) {
		golog.CtsLog.Error("FAIL to Connect with EthernetIp[%s:%s]", device.Address, device.Port)
		return false

	} else {
		golog.CtsLog.Info("Connected with EthernetIp[%s:%s]", device.Address, device.Port)
		return true
	}
}
