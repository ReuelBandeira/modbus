package main

import (
	"gateway/globals"
	"gateway/keys"
	"gateway/models"
	"gateway/mqtt"
	"gateway/process"
	"gateway/protocols"
	"gateway/utils"
	"gateway/version"
	"log"
	"time"

	"os"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	// Check if gateway is already running
	if utils.Running() {
		log.Println("Gateway is already running")
		os.Exit(0)
	}

	start := time.Now()

	// Will contain all configurations from config file
	globals.Conf = &models.Configurations{}

	// Read config file
	viper.SetConfigName("gateway_conf")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; configure log files with default values
			log.SetOutput(&lumberjack.Logger{
				Filename:   "gateway.log",
				MaxSize:    1,    // megabytes
				MaxBackups: 3,    // number of files
				MaxAge:     28,   // days
				Compress:   true, // disabled by default
			})
		} else {
			// Config file was found but another error was produced
			log.Println("Error reading configuration file: ", err)
		}
	} else {

		// Get all configurations from config file
		err = viper.Unmarshal(globals.Conf)
		if err != nil {
			log.Printf("unable to decode into config struct, %v", err)
		}

		// Configure log files with values from configuration file
		log.SetOutput(&lumberjack.Logger{
			Filename:   globals.Conf.Log.Filename,
			MaxSize:    globals.Conf.Log.MaxSize,    // megabytes
			MaxBackups: globals.Conf.Log.MaxBackups, // number of files
			MaxAge:     globals.Conf.Log.MaxAge,     // days
			Compress:   globals.Conf.Log.Compress,   // disabled by default
		})
	}

	// Log gateway version
	log.Println("===================================================================================")
	log.Println("ISI 4.0 Gateway - Vrs. 1.0.0")

	// Gets path separator from the underlying OS
	globals.PathSeparator = string(os.PathSeparator)

	// Connects to local MQTT broker
	if !mqtt.ConnectLocalBroker() {
		log.Println("ERROR: local MQTT connection failed")
		return
	}

	// Connects to external MQTT broker
	mqtt.ConnectGlobalBroker()

	// Verify protocols' versions
	version.CheckVersions()

	// Starts all enabled protocols
	process.Initialize()

	// Get initial protocols' devices configurations
	protocols.GetDevices()

	// Monitoring keys and leds on the hardware
	go keys.ScanKeys()

	// Keep alive protocols
	if globals.Conf.KeepAlive.Interval > 0 {
		go process.SendProtocolsKeepAlive()
	}

	duration := time.Since(start)
	log.Println("=====> Gateway started in ", duration)

	select {}
}
