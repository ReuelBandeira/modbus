package newprotocol

import (
	utilsPkg "newprotocol/utils"
	logPkg "newprotocol/utils/gologtofile"
	"os"
	"runtime"
	"testing"
)

var bLogStarted = false

func Test01_IoMsg(t *testing.T) {
	startLog()
	var filePath string = "deviceConfig.json"
	devStruct, err := readJSON(filePath)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		logPkg.CtsLog.Debug("devStruct[%v]", devStruct)
	}
}

func startLog() {
	if !bLogStarted {
		utilsPkg.Hostname, _ = os.Hostname()
		// Initialize CtsLogs with default parameters
		rc, err := logPkg.InitCtsLogs(
			utilsPkg.DefaultLogFilePath,
			"newprotocol_test",
			utilsPkg.DeleteExistingLogFiles,
			utilsPkg.DefaultLogLevel,
			utilsPkg.LogFileMaxSize,
		)
		if (err != nil) || (rc < 0) {
			/// Fail to setup Logs
			logPkg.CtsLog.Error("Test01:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n         err[%s]\n\n",
				rc,
				runtime.GOOS,
				runtime.GOARCH,
				utilsPkg.OsBits,
				utilsPkg.Hostname,
				utilsPkg.DefaultLogFilePath,
				"newprotocol_test",
				utilsPkg.DefaultLogLevel,
				err)
			return
		}
		logPkg.CtsLog.Info("Test01:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d]\n\n",
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.DefaultLogFilePath,
			"newprotocol_test",
			utilsPkg.DefaultLogLevel)
	}
}
