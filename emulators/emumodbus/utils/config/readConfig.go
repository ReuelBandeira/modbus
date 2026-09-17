package config

import (
	"encoding/json"
	/// TODO: change newprotocol with your protocolname
	utilsPkg "emumodbus/utils"
	logPkg "emumodbus/utils/gologtofile"
	"os"
	"sync"
)

var SettingsMqtt utilsPkg.SettingsMQTT

// ==============================================
// GetDeviceConfigInfo: >>> deviceConfig.json <<<
// ==============================================
// [
//
//		{
//		  "devices": [   // var devInf []utilsPkg.Devices
//		    {
//		      "address": "can", // var itemInf []utilsPkg.DevSettings
//		      "port": "0",
//		      "name": "CANopen",
//		      "protocol": "canopen",
//		      "readingTime": 20,
//		      "topics": [
//		        "alldevices",
//		        "canopen-test"
//		      ],
//		      "data": [ // []map[string]interface{}
//		        {
//	 =========================================
//	 TODO: according with protocol_config.json
//	 =========================================
//		        } // END: "data" {}
//		      ] // END: "data"
//		    } // END: "devices" {}
//		  ] // END: "devices" []
//		}
//
// ]
func GetDeviceConfigInfo() ([]utilsPkg.DevSettings, error) {
	var devInf []utilsPkg.Devices
	var itemInf []utilsPkg.DevSettings
	var devConfig string = "deviceConfig.json"

	file, err := os.ReadFile(devConfig)
	if err != nil {
		logPkg.CtsLog.Error("Fail to read JSON File[%s]\nerr[%v]\n", devConfig, err)
		return itemInf, err
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		logPkg.CtsLog.Error("Fail to decode JSON File[%s]\nerr[%v]\n devInf[%v]\n", devConfig, err, devInf)
		return itemInf, err
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			itemInf = append(itemInf,
				utilsPkg.DevSettings{
					Address:     item.Address,
					Port:        item.Port,
					Name:        item.Name,
					Protocol:    item.Protocol,
					Id:          item.Id, // Added 24/08/2023
					ReadingTime: item.ReadingTime,
					Topics:      item.Topics,
					Data:        item.Data, // []map[string]interface{}
				})
		}
	}
	return itemInf, nil
}

// GoRotine:Channels Control
var (
	DoneChan      chan bool
	SyncOnceChan  sync.Once
	SyncWaitGroup sync.WaitGroup
)

func CreateGoRotineChannel() {
	DoneChan = make(chan bool)
}
