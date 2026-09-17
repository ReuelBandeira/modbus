package main

import (
	"fmt"
	"icts/modbus/configs/logger"
	"time"

	configsPkg "icts/modbus/configs"
	packagesPkg "icts/modbus/pkg"
)

func main() {

	logger.Log.Info("Starting Modbus.")
	err := startingLocalBroker()
	if err != nil {
		logger.Log.Error("Is not possible start communication with broker.")
		return
	}

	packagesPkg.SendStatusProtocol("running")
	time.Sleep(500 * time.Millisecond)
	packagesPkg.StartReadFile()

	configsPkg.StartRead()

	for {
		fmt.Println("Checking...")
		time.Sleep(20000 * time.Millisecond)
	}
}

func startingLocalBroker() error {
	configsPkg.SetMqttBroker()
	err := configsPkg.MqttCommunication()
	return err
}
