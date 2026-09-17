package logger

import (
	"io" // Adicione esta linha para importar o pacote io
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Configurar um formato personalizado para o horário (time)
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	//Change the level to see specific informations.
	Log.SetLevel(logrus.DebugLevel)

	// Create log file
	logFile, err := os.Create("modbus.log")
	if err != nil {
		Log.Fatal("Could not create log file:", err)
	}

	// Create a io.MultiWriter to send logs for file and prompt
	logOutput := io.MultiWriter(os.Stdout, logFile)
	Log.SetOutput(logOutput)
}
