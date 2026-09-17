package main

import (
	initMqtt "canopen/mqtt" // Added 24/08/2023
	protocol "canopen/protocol"
	utilsPkg "canopen/utils"
	cfgPkg "canopen/utils/config"
	logPkg "canopen/utils/gologtofile"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {

	// Check if it proccess is running
	if utilsPkg.Running() {
		fmt.Println("=====> WARNING: Another instance of the canopen is already running")
		os.Exit(0)
	}

	// Get Running Hostname
	utilsPkg.Hostname, _ = os.Hostname()

	// Initialize CtsLogs with default parameters
	rc, err := logPkg.InitCtsLogs(
		utilsPkg.DefaultLogFilePath,
		utilsPkg.DefaultLogFileName,
		utilsPkg.DeleteExistingLogFiles,
		utilsPkg.DefaultLogLevel,
		utilsPkg.LogFileMaxSize,
	)
	if (err != nil) || (rc < 0) {
		/// Fail to setup Logs
		/// TODO: Channel to Send Fail Information to Backend using Internal MQTT Brocker
		logPkg.CtsLog.Error("main:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n         err[%s]\n\n",
			rc,
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.DefaultLogFilePath,
			utilsPkg.DefaultLogFileName,
			utilsPkg.DefaultLogLevel,
			err)
	} else {
		logPkg.CtsLog.Info("main:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n\n",
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.DefaultLogFilePath,
			utilsPkg.DefaultLogFileName,
			utilsPkg.DefaultLogLevel)
	}

	// Force Warning Log Levels Onpy
	logPkg.CtsLog.SetLogLevel(logPkg.WarnLevel)

	// Connecting to Local MQTT Brocker
	err = initMqtt.ConnectToLocalMqttBrocker()
	if err != nil {
		err = initMqtt.StartLocalMQTTServer()
		if err != nil {
			/// FAIL to Connect to Local MQTT Brocker
			logPkg.CtsLog.Error("main:ConnectToLocalMqttBrocker()\n err[%v]\n", err)
			return
		}
		err = initMqtt.ConnectToLocalMqttBrocker()
		if err != nil {
			logPkg.CtsLog.Error("main:ConnectToLocalMqttBrocker()\n err[%v]\n", err)
			return
		}
	}

	// Sending Status that Protocol is Running
	if !utilsPkg.ProtocolIsRunning {
		utilsPkg.ProtocolIsRunning = true
		protocol.MQTTSendProtocolStatus("running")
	}

	//Start Read This Industrial Protocols
	cfgPkg.StartProtocol()

	// =============
	// Infinite Loop
	// =============
	logPkg.CtsLog.Debug("Loop Starting ... \n")
	for {
		time.Sleep(5000 * time.Millisecond)
	}
}
