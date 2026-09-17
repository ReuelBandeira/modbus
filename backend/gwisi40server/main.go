package main

import (
	"gwisi40server/database"
	"gwisi40server/globals"
	"gwisi40server/mqtt"
	"gwisi40server/routes"
	"gwisi40server/utils"
	"io"
	"log"
	"os"
	"time"

	_ "gwisi40server/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

//	@title			ISI 4.0 server API
//	@version		1.0
//	@description	ISI 4.0 server API

//	@contact.name	ICTS
//	@contact.url	http://www.grupoicts.com.br

//	@SecurityDefinitions	BasicAuth
//	@In						header
//	@Name					Authorization

// @license.name	Proprietary
func main() {

	if utils.Running() {
		log.Println("=====> WARNING: Another instance of the server is already running")
		os.Exit(0)
	}

	start := time.Now()

	args := os.Args[1:]

	var initialize bool = true

	if len(args) != 0 {
		if args[0] == "clean" {
			initialize = false
		} else {
			initialize = true
		}
	}

	viper.SetConfigName("backend_conf")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file was not found
			globals.Conf.JWT.Expiration = 3600
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
	}

	// Configures logger
	logFile := &lumberjack.Logger{
		Filename:   "backend.log", // The name of the log file
		MaxSize:    1,             // Max size of the log file (in megabytes)
		MaxBackups: 3,             // Max number of old log files to retain
		MaxAge:     28,            // Max number of days to keep old log files
		Compress:   true,          // Whether to compress old log files (gzip)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Println("===================================================================================")
	log.Println("Starting API")
	log.Println("OS:", os.Getenv("OS"))

	// Creates a new Fiber instance
	app := fiber.New(fiber.Config{
		StrictRouting:         false,
		DisableStartupMessage: true,
	})

	// Configures CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://*, https://*, *",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin" +
			", X-Custom-Configurations-Version, X-Custom-Protocol-Version, X-Custom-MQTT-Version",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS",
	}))

	// Configures API logger
	app.Use(logger.New(logger.Config{
		Output:     logFile,                                           // Use the Lumberjack logger instance
		Format:     "${time} ${ip} - ${status} - ${method} ${path}\n", // Specify log format
		TimeFormat: "2006/01/02 15:04:05",                             // Specify the date and time format
	}))

	// Connects to database and migrate tables if needed
	if err := database.ConnectDB(initialize); err != nil {
		panic("Could not connect to database")
	}

	// Serves frontend static files
	app.Static("/", "./web/")

	// Configures route to serve Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Configures API endpoints routes
	routes.Setup(app)

	// Starts local MQTT server
	mqtt.StartLocalMQTTServer()

	// Connects to external MQTT server
	go mqtt.InitializeExternalBrokerConnection()

	// Connects to local MQTT server
	if !mqtt.InitializeLocalBrokerConnection() {
		log.Println("=====> WARNING: Could not connect to local MQTT Broker")
	}

	// Subscribes to MQTT topics that indicates the device's status
	mqtt.SubscribeStatus()

	// Let's run...
	duration := time.Since(start)
	log.Println("=====> Backend started in ", duration)

	log.Println("Server started in port 8585")
	app.Listen(":8585")
}
