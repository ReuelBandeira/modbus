package pkg

import (
	"fmt"
	"icts/modbus/configs/logger"
	utilsPkg "icts/modbus/utils"
	"time"

	modbus "github.com/goburrow/modbus"
)

func ConnModbus(device utilsPkg.Device) bool {
	mdbsDevice := modbus.NewTCPClientHandler(device.Address + ":" + device.Port)
	timeWait := 5 * time.Second
	connectionTry := 1
	for {		
		select{
			//Verifica se o canal não foi fechado
			case <-utilsPkg.StopChannels[device.ID]:
				logger.Log.Infof("modbusCheck:ConnModbus: Stop Modbus Device Name[%v] ID[%v].",
					device.Name, device.ID)
				utilsPkg.DeviceWorkStatus[device.ID] = "stopped"
				return false
			default:
				//Tenta conectar com o dispositivo
				err := mdbsDevice.Connect()
				if err != nil {
					utilsPkg.DeviceWorkStatus[device.ID] = "trying"
					utilsPkg.MsgLog = utilsPkg.MessageLog{
						MessageType: "error",
						Data: utilsPkg.DataLog{
							Address:  device.Address,
							Port:     device.Port,
							Name:     device.Name,
							Protocol: device.Protocol,
							Time:     time.Now().Format("2006-01-02 15:04:05"),
							Message: fmt.Sprintf("Connection Fail. Error [%v]. Attempt [%d]. Waiting [%vs]",
								err, connectionTry, timeWait.Seconds()),
							Code:   "",
							Source: "modbusCheck:ConnModbus",
						},
					}
					Sender(utilsPkg.MsgLog, utilsPkg.Log)
					
					//Send this mesage only one time
					if connectionTry  < 2{
						SendStatusDevice(device, "disconnected")
						logger.Log.Errorf("modbusCheck:ConnModbus: Connection Fail - Device [%s] is Offline. Try connection[%d]",
						mdbsDevice.Address, connectionTry)
					}
					
					logger.Log.Debugf("routineMdbRead: Waiting [%v] seconds to try reconnect. Attempt [%d]", timeWait.Seconds(), connectionTry)
					time.Sleep(timeWait)
					connectionTry++
		
				} else {
					logger.Log.Infof("modbusCheck:ConnModbus: Connection Sucess - Device [%s] is Online.",
						mdbsDevice.Address)
					SendStatusDevice(device, "connected")
					return true		
				}

		}
		
		
	}

}
