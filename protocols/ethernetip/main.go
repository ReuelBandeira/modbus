package main

import (
	initMqtt "ethernetip/internal/mqtt"
	cfgPkg "ethernetip/utils"
	initPkg "ethernetip/utils/config"
	logPkg "ethernetip/utils/gologtofile"
)

func main() {
	// First: Initialize CtsLogs with default parameters
	rc, err := logPkg.InitCtsLogs(
		cfgPkg.DefaultGwIsi40LogFilePath,
		cfgPkg.DefaultGwIsi40LogFileName,
		cfgPkg.DeleteExistingLogFiles,
		cfgPkg.DefaultGwIsi40LogLevel,
	)

	if (err != nil) || (rc < 0) {
		logPkg.CtsLog.Error("InitCtsLogs:FAIL: First: Initialize CtsLogs with default parameters\n     LogFilePath[%s]\n     LogFileName[%s]\n        LogLevel[%d]\n             err:[%s]\n", cfgPkg.DefaultGwIsi40LogFilePath, cfgPkg.DefaultGwIsi40LogFileName, cfgPkg.DefaultGwIsi40LogLevel, err)
	} else {
		logPkg.CtsLog.Warn("InitCtsLogs:WARN First: Initialize CtsLogs with default parameters\n     LogFilePath[%s]\n     LogFileName[%s]\n        LogLevel[%d]\n", cfgPkg.DefaultGwIsi40LogFilePath, cfgPkg.DefaultGwIsi40LogFileName, cfgPkg.DefaultGwIsi40LogLevel)
	}

	err = initMqtt.StartingLocalBroker()
	if err != nil {
		logPkg.CtsLog.Error("FAIL - Is not possible start communication with broker")
		return
	}
	//Start Read Industrial Protocols EthernetIP
	initPkg.StartEthernetip()

	// =============
	// Infinite Loop
	// =============
	logPkg.CtsLog.Debug("Loop Starting ... \n")

	select {}
}
