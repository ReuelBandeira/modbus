package database

import (
	"gwisi40server/models"
	"gwisi40server/utils"
)

var (
	language  models.Language
	users     []models.User
	protocols []models.Protocol
	types     []models.Type
	servers   []models.Server
	plcs      []models.PLC
	// devices        []models.Device
	mqqt           []models.MQTTSettings
	configurations []models.Configurations
	version        []models.DataVersion
)

func InitializeDB() {

	// Inspect our database. If no information exists, creates our initial configuration.

	initializeDefaults()
	initializeUsers()
	initializeProtocols()
	initializeTypes()
	initializeServers()
	initializePlcs()
	// initializeDevices()
	// initializeDeviceConfiguration()
	initializeMqttSettings()
	initializeConfigurations()
	initializeDataVersion()

}

func initializeDataVersion() {
	result := DB.Find(&version)

	if result.RowsAffected == 0 {
		version := []models.DataVersion{
			{
				DataType: "Configurations",
				Version:  1,
			},
			{
				DataType: "Protocol",
				Version:  0,
			},
			{
				DataType: "MQTT",
				Version:  1,
			},
		}

		for _, c := range version {
			DB.Create(&c)
		}
	}
}

var devethernetip1 = `	
{
  "bitMemories": [
    {
      "attribute": 0,
      "class": 849,
      "instance": 1,
      "name": "Y00"
    },
    {
      "attribute": 1,
      "class": 849,
      "instance": 1,
      "name": "Y01"
    },
    {
      "attribute": 2,
      "class": 849,
      "instance": 1,
      "name": "Y02"
    },
    {
      "attribute": 3,
      "class": 849,
      "instance": 1,
      "name": "Y03"
    },
    {
      "attribute": 4,
      "class": 849,
      "instance": 1,
      "name": "Y04"
    },
    {
      "attribute": 5,
      "class": 849,
      "instance": 1,
      "name": "Y05"
    }
  ],
  "wordMemories": [
    {
      "attribute": 0,
      "class": 850,
      "instance": 2,
      "name": "D0"
    },
    {
      "attribute": 1,
      "class": 850,
      "instance": 2,
      "name": "D1"
    },
    {
      "attribute": 2,
      "class": 850,
      "instance": 2,
      "name": "D2"
    },
    {
      "attribute": 3,
      "class": 850,
      "instance": 2,
      "name": "D3"
    },
    {
      "attribute": 4,
      "class": 850,
      "instance": 2,
      "name": "D4"
    },
    {
      "attribute": 5,
      "class": 850,
      "instance": 2,
      "name": "D5"
    }
  ]
}
`

var devethernetip2 = `
{
  "bitMemories": [
    {
      "attribute": 100,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ01_REC"
    },
    {
      "attribute": 101,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ01_AVAN"
    },
    {
      "attribute": 102,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ02_REC"
    },
    {
      "attribute": 103,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ02_AVAN"
    },
    {
      "attribute": 104,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ03_REC"
    },
    {
      "attribute": 105,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ03_AVAN"
    },
    {
      "attribute": 106,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ04_REC"
    },
    {
      "attribute": 107,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_CIL_AJ04_AVAN"
    },
    {
      "attribute": 108,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_STOPPER_01"
    },
    {
      "attribute": 109,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_STOPPER_02"
    },
    {
      "attribute": 110,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_STOPPER_SAIDA"
    },
    {
      "attribute": 111,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_ESPERA01"
    },
    {
      "attribute": 112,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_ESPERA02"
    },
    {
      "attribute": 113,
      "class": 851,
      "instance": 1,
      "name": "M01_IN_ESPERA03"
    },
    {
      "attribute": 118,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_ESPERA_01"
    },
    {
      "attribute": 119,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_ESPERA_02"
    },
    {
      "attribute": 120,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_SAIDA"
    },
    {
      "attribute": 121,
      "class": 851,
      "instance": 1,
      "name": "ABRE_ESTEIRA"
    },
    {
      "attribute": 122,
      "class": 851,
      "instance": 1,
      "name": "FECHA_ESTEIRA"
    },
    {
      "attribute": 123,
      "class": 851,
      "instance": 1,
      "name": "LIGA_ESTEIRA"
    },
    {
      "attribute": 128,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_ESPERA_01"
    },
    {
      "attribute": 129,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_ESPERA_02"
    },
    {
      "attribute": 130,
      "class": 851,
      "instance": 1,
      "name": "BAIXA_STOPPER_SAIDA"
    },
    {
      "attribute": 131,
      "class": 851,
      "instance": 1,
      "name": "ABRE_ESTEIRA"
    },
    {
      "attribute": 132,
      "class": 851,
      "instance": 1,
      "name": "FECHA_ESTEIRA"
    },
    {
      "attribute": 133,
      "class": 851,
      "instance": 1,
      "name": "LIGA_ESTEIRA"
    },
    {
      "attribute": 200,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_STOPPER"
    },
    {
      "attribute": 201,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_PRESENÇA_BDJ"
    },
    {
      "attribute": 202,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_CIL_LIFT_MANOR_REC"
    },
    {
      "attribute": 203,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_CIL_LIFT_MANOR_AVAN"
    },
    {
      "attribute": 204,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_CIL_LIFT_MAIOR_REC"
    },
    {
      "attribute": 205,
      "class": 851,
      "instance": 1,
      "name": "M02_IN_CIL_LIFT_MAIOR_AVAN"
    },
    {
      "attribute": 210,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_STOPPER_LIFT_01"
    },
    {
      "attribute": 211,
      "class": 851,
      "instance": 1,
      "name": "M14_SOBE_LIFT_UDIM"
    },
    {
      "attribute": 212,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_LIFT_UDIM"
    },
    {
      "attribute": 213,
      "class": 851,
      "instance": 1,
      "name": "M14_SOBE_LIFT_SUDIM"
    },
    {
      "attribute": 214,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_LIFT_SUDIM"
    },
    {
      "attribute": 220,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_STOPPER_LIFT_01"
    },
    {
      "attribute": 221,
      "class": 851,
      "instance": 1,
      "name": "M14_SOBE_LIFT_UDIM"
    },
    {
      "attribute": 222,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_LIFT_UDIM"
    },
    {
      "attribute": 223,
      "class": 851,
      "instance": 1,
      "name": "M14_SOBE_LIFT_SUDIM"
    },
    {
      "attribute": 224,
      "class": 851,
      "instance": 1,
      "name": "M14_DESCE_LIFT_SUDIM"
    },
    {
      "attribute": 300,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_STOPPER"
    },
    {
      "attribute": 301,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_PRESENÇA_BDJ"
    },
    {
      "attribute": 302,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_CIL_LIFT_MAIOR_AVAN"
    },
    {
      "attribute": 303,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_CIL_LIFT_MAIOR_REC"
    },
    {
      "attribute": 304,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_CIL_LIFT_MENOR_AVAN"
    },
    {
      "attribute": 305,
      "class": 851,
      "instance": 1,
      "name": "M03_IN_CIL_LIFT_MENOR_REC"
    }
  ]
}
`

var devmodbus1 = `
{
  "bitMemories": [
    {
      "address": 40964,
      "name": "Y0.4"
    }
  ],
  "slaveId": 1,
  "wordMemories": [
    {
      "address": 0,
      "format": 16,
      "name": "D0"
    },
    {
      "address": 1,
      "format": 16,
      "name": "D1"
    },
    {
      "address": 2,
      "format": 16,
      "name": "D2"
    },
    {
      "address": 3,
      "format": 16,
      "name": "D3"
    },
    {
      "address": 4,
      "format": 16,
      "name": "D4"
    },
    {
      "address": 5,
      "format": 16,
      "name": "D5"
    },
    {
      "address": 6,
      "format": 16,
      "name": "D6"
    }
  ]
}
`

var devmodbus2 = `
{
  "bitMemories": [
    {
      "address": 0,
      "name": "BOTAO_EMERGENCIA"
    },
    {
      "address": 1,
      "name": "BOTÃO START"
    },
    {
      "address": 2,
      "name": "BOTÃO STOP"
    },
    {
      "address": 3,
      "name": "BOTÃO REARME"
    },
    {
      "address": 4,
      "name": "AUX CORTINA"
    },
    {
      "address": 5,
      "name": "S_IND_C"
    },
    {
      "address": 6,
      "name": "DRIVER X ERRO"
    },
    {
      "address": 7,
      "name": "DRIVER X MOV"
    },
    {
      "address": 8,
      "name": "DRIVER Y ERRO"
    },
    {
      "address": 9,
      "name": "DRIVER Y MOV"
    },
    {
      "address": 10,
      "name": "DRIVER Z ERRO"
    },
    {
      "address": 11,
      "name": "DRIVER Z MOV"
    },
    {
      "address": 12,
      "name": "S PORTA AVA"
    },
    {
      "address": 13,
      "name": "S PORTA REC"
    },
    {
      "address": 14,
      "name": "S CIL A AVA"
    },
    {
      "address": 15,
      "name": "S CIL A REC"
    },
    {
      "address": 16,
      "name": "S CIL B AVA"
    },
    {
      "address": 17,
      "name": "S CIL B REC"
    },
    {
      "address": 18,
      "name": "S CIL C AVA"
    },
    {
      "address": 19,
      "name": "S CIL C REC"
    },
    {
      "address": 20,
      "name": "S CIL D AVA"
    },
    {
      "address": 21,
      "name": "S CIL D RED"
    },
    {
      "address": 22,
      "name": "S CIL LAT AVA"
    },
    {
      "address": 23,
      "name": "S CIL LAT REC"
    },
    {
      "address": 24,
      "name": "S CIL GARRA AVA"
    },
    {
      "address": 25,
      "name": "S CIL GARRA REC"
    },
    {
      "address": 26,
      "name": "S IND A"
    },
    {
      "address": 27,
      "name": "S IND B"
    },
    {
      "address": 28,
      "name": "DOOR CLOSE"
    },
    {
      "address": 29,
      "name": "PORTA CLOSE"
    },
    {
      "address": 30,
      "name": "TRAVA MAGNET OK"
    }
  ],
  "slaveId": 1,
  "wordMemories": []
}
`

var devcanopen = `
{
  "SDO": [
    {
      "name": "SDO-Download",
      "canid": "0x600",
      "nodeid": 1,
      "cmd": 2,
      "objIndex": 3,
      "subIndex": 4
    },
    {
      "name": "SDO-Upload",
      "canid": "0x580",
      "nodeid": 5,
      "cmd": 6,
      "objIndex": 7,
      "subIndex": 8
    }
  ],
  "PDO": [
    {
      "name": "PDO0",
      "address": 2000
    },
    {
      "name": "PDO1",
      "address": 2001
    },
    {
      "name": "PDO2",
      "address": 2002
    }
  ],
  "NMT": [
    {
      "name": "NMT0",
      "address": 3000
    },
    {
      "name": "NMT1",
      "address": 3001
    },
    {
      "name": "NMT2",
      "address": 3002
    }
  ],
  "SYNC": [
    {
      "name": "SYNC0",
      "address": 4000
    },
    {
      "name": "SYNC1",
      "address": 4001
    },
    {
      "name": "SYNC2",
      "address": 4002
    }
  ],
  "TIME": [
    {
      "name": "TIME0",
      "address": 5000
    },
    {
      "name": "TIME1",
      "address": 5001
    },
    {
      "name": "TIME2",
      "address": 5002
    }
  ],
  "EMCY": [
    {
      "name": "EMCY0",
      "address": 6000
    },
    {
      "name": "EMCY1",
      "address": 6001
    },
    {
      "name": "EMCY2",
      "address": 6002
    }
  ],
  "devid": 333
}
`

func initializeConfigurations() {
	result := DB.Find(&configurations)

	if result.RowsAffected == 0 {
		configurations := []models.Configurations{
			{
				Version:     1,
				Address:     "172.16.16.61",
				Port:        "44818",
				Name:        "CLP1",
				Protocol:    "ethernetip",
				ReadingTime: 1000,
				Topics:      "alldevices",
				Devices:     devethernetip1,
			},
			{
				Version:     1,
				Address:     "192.168.3.3",
				Port:        "44818",
				Name:        "CLP_PROJETO_ADATA",
				Protocol:    "ethernetip",
				ReadingTime: 1000,
				Topics:      "alldevices",
				Devices:     devethernetip2,
			},
			{
				Version:     1,
				Address:     "172.16.16.61",
				Port:        "502",
				Name:        "PLC asdasds",
				Protocol:    "modbus",
				ReadingTime: 300,
				Topics:      "plcmel,modbus-test",
				Devices:     devmodbus1,
			},
			{
				Version:     1,
				Address:     "172.16.16.61",
				Port:        "502",
				Name:        "PLC-AJAX-SIMULADOR",
				Protocol:    "modbus",
				ReadingTime: 2000,
				Topics:      "modbus,ajax",
				Devices:     devmodbus2,
			},
			{
				Version:     1,
				Address:     "can",
				Port:        "0",
				Name:        "CAN-Ifc",
				Protocol:    "canopen",
				ReadingTime: 5000,
				Topics:      "alldevices,canopen-test",
				Devices:     devcanopen,
			},
		}

		for _, c := range configurations {
			DB.Create(&c)
		}
	}
}

func initializeUsers() {
	result := DB.Find(&users)

	if result.RowsAffected == 0 {
		passwd, _ := utils.HashPassword("123456")
		user := models.User{
			Name:     "ISI40 Admin",
			Email:    "isi40@grupoicts.com.br",
			Password: passwd,
		}

		DB.Create(&user)

		passwd, _ = utils.HashPassword("a")
		user = models.User{
			Name:     "a",
			Email:    "a@a.com",
			Password: passwd,
		}

		DB.Create(&user)
	}
}

func initializeDefaults() {
	result := DB.First(&language)

	if result.RowsAffected == 0 {
		language = models.Language{
			Language: "pt",
		}

		DB.Create(&language)
	}
}

func initializeMqttSettings() {
	DB.Find(&mqqt)

	if len(mqqt) == 0 {
		mqtt := models.MQTTSettings{
			Server:       "191.252.214.157",
			Port:         "1883",
			Username:     "",
			Password:     "",
			InputInfo:    "inputinfo",
			ErrorInfo:    "errorinfo",
			StatusDevice: "statusdevice",
		}

		DB.Create(&mqtt)
	}
}

func initializeProtocols() {
	// Protocols...
	DB.Find(&protocols)

	if len(protocols) == 0 {
		protocols := []models.Protocol{
			{
				Protocol: "ethernetip",
				Version:  0,
				Config:   "",
			},
			{
				Protocol: "modbus",
				Version:  0,
				Config:   "",
			},
			{
				Protocol: "canopen",
				Version:  0,
				Config:   "",
			},
		}

		for _, c := range protocols {
			DB.Create(&c)
		}
	}
}

func initializeTypes() {
	// Types...
	DB.Find(&types)

	if len(types) == 0 {
		types := []models.Type{
			{
				Type: "DigitalOutput",
			},
			{
				Type: "AnalogOutput",
			},
			{
				Type: "AnalogInput",
			},
		}

		for _, c := range types {
			DB.Create(&c)
		}
	}
}

func initializePlcs() {
	// PLCs...
	DB.Find(&plcs)

	if len(plcs) == 0 {
		plcs := []models.PLC{
			{
				Name: "ETHERNETIP_PLC1",
			},
		}

		for _, c := range plcs {
			DB.Create(&c)
		}
	}
}

func initializeServers() {
	// Servers...
	DB.Find(&servers)

	if len(servers) == 0 {
		var ethernet models.Protocol
		DB.First(&ethernet, "protocol = ?", "ethernetip")

		var modbus models.Protocol
		DB.First(&modbus, "protocol = ?", "modbus")

		servers := []models.Server{
			{
				Name:       "Ethernet/IP Server",
				IP:         "127.0.0.1",
				Port:       44818,
				ProtocolID: ethernet.ID,
				Protocol:   ethernet,
			},
			{
				Name:       "Modbus Server",
				IP:         "127.1.1.1",
				Port:       44819,
				ProtocolID: modbus.ID,
				Protocol:   modbus,
			},
		}

		for _, c := range servers {
			DB.Create(&c)
		}
	}
}

// func initializeDevices() {
// 	// Servers...
// 	DB.Find(&devices)

// 	if len(devices) == 0 {
// 		var ethernet models.Protocol
// 		DB.First(&ethernet, "protocol = ?", "Ethernet/IP")

// 		var modbus models.Protocol
// 		DB.First(&modbus, "protocol = ?", "Modbus TCP")

// 		var digitalOutput models.Type
// 		DB.First(&digitalOutput, "type = ?", "DigitalOutput")

// 		var analogOutput models.Type
// 		DB.First(&analogOutput, "type = ?", "AnalogOutput")

// 		var analogInput models.Type
// 		DB.First(&analogInput, "type = ?", "AnalogInput")

// 		var plc1 models.PLC
// 		DB.First(&plc1, "name = ?", "ETHERNETIP_PLC1")

// 		devices := []models.Device{
// 			{
// 				Name:        "Caldeira",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      digitalOutput.ID,
// 				MemoryBlock: 1,
// 				Class:       4,
// 				Instance:    101,
// 				Attribute:   3,
// 				Address:     0,
// 			},
// 			{
// 				Name:        "Queimador",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      digitalOutput.ID,
// 				MemoryBlock: 1,
// 				Class:       4,
// 				Instance:    101,
// 				Attribute:   3,
// 				Address:     1,
// 			},
// 			{
// 				Name:        "Purgador",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      digitalOutput.ID,
// 				MemoryBlock: 7,
// 				Class:       4,
// 				Instance:    101,
// 				Attribute:   3,
// 				Address:     2,
// 			},
// 			{
// 				Name:        "ValvulaEntrada",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      analogOutput.ID,
// 				MemoryBlock: 2,
// 				Class:       4,
// 				Instance:    102,
// 				Attribute:   3,
// 				Address:     0,
// 			},
// 			{
// 				Name:        "IntensidadeQueimador",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      analogOutput.ID,
// 				MemoryBlock: 2,
// 				Class:       4,
// 				Instance:    102,
// 				Attribute:   3,
// 				Address:     1,
// 			},
// 			{
// 				Name:        "ControladorTemperatura",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      analogOutput.ID,
// 				MemoryBlock: 2,
// 				Class:       4,
// 				Instance:    102,
// 				Attribute:   3,
// 				Address:     2,
// 			},
// 			{
// 				Name:        "SensorNivelDoReservatorio",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      analogInput.ID,
// 				MemoryBlock: 3,
// 				Class:       4,
// 				Instance:    103,
// 				Attribute:   3,
// 				Address:     0,
// 			},
// 			{
// 				Name:        "SensorTemperatura",
// 				ProtocolID:  ethernet.ID,
// 				PLCID:       plc1.ID,
// 				TypeID:      analogInput.ID,
// 				MemoryBlock: 3,
// 				Class:       4,
// 				Instance:    103,
// 				Attribute:   3,
// 				Address:     1,
// 			},
// 		}

// 		for _, c := range devices {
// 			DB.Create(&c)
// 		}
// 	}

// }
