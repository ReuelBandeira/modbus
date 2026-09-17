package gologtofile

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// ===========================================================
// 2022/06/22 - OLIVEIRA: BEGIN NEW LOGGING MODE WITH LOG MODE
// ===========================================================
var CtsLogFileHandle *os.File // LogFileHandleLog File Name: AAAAMMDD_LogFileHandle.log

type NewLogLevel int

const (
	DisableLevel NewLogLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type Logger struct {
	*log.Logger
	curLogLevel NewLogLevel
}

var CtsLog *Logger = nil

var ExitApp bool = false
var CurLoglevel int = int(InfoLevel)

// SetLogLevel function set current loglevel
// 0=disabled,1=fatal,2=error,3=warnning,4=info,5=debug,>5 unknown
func (l *Logger) SetLogLevel(newLevel NewLogLevel) int {
	CurLoglevel = int(newLevel)
	// Create ErrorLog, WarnLog and InfoLog  Type
	CtsLog = NewLogger(newLevel, CtsLogFileHandle, "", log.LstdFlags|log.Lshortfile|log.Ltime)
	CtsLog.SetOutput(io.MultiWriter(os.Stdout, CtsLogFileHandle))
	return CurLoglevel
}

// GetLogLevel function return the current loglevel
// 0=disabled,1=fatal,2=error,3=warnning,4=info,5=debug,>5 unknown
func (l *Logger) GetLogLevel() (int, string) {
	switch CurLoglevel {
	case int(DisableLevel):
		return CurLoglevel, "DisableLevel"
	case int(FatalLevel):
		return CurLoglevel, "FatalLevel"
	case int(ErrorLevel):
		return CurLoglevel, "ErrorLevel"
	case int(WarnLevel):
		return CurLoglevel, "WarnLevel"
	case int(InfoLevel):
		return CurLoglevel, "InfoLevel"
	case int(DebugLevel):
		return CurLoglevel, "DebugLevel"
	default:
		return CurLoglevel, "UnknownLevel"
	}

}

// GetLogLevel function return the current loglevel
// 0=disabled,1=fatal,2=error,3=warnning,4=info,5=debug,>5 unknown
func (l *Logger) GetLogLevelByName(LogLevelName string) int {
	switch LogLevelName {
	case "DisableLevel":
		return int(DisableLevel)
	case "FatalLevel":
		return int(FatalLevel)
	case "ErrorLevel":
		return int(ErrorLevel)
	case "WarnLevel":
		return int(WarnLevel)
	case "InfoLevel":
		return int(InfoLevel)
	case "DebugLevel":
		return int(DebugLevel)
	}
	if l.curLogLevel >= WarnLevel {
		l.Output(2, fmt.Sprintf("[WRN] "+"LogLevelName[%s]Unknown return [6]", LogLevelName))
	}
	return 6
}

// NewLogger Initialise a logger
// @param1 >> newLogLevel - 0=disable,1=fatal,2=error,3=warning,4=info,5=debug
// @param2 >> out - where the log will be written
// @param3 >> prefix - log prefix ex: [ERR]..
// @param4 >> flag   - Log Flags
// @return1 << *Logger - Pointer of desired logger
func NewLogger(newLogLevel NewLogLevel, out io.Writer, prefix string, flag int) *Logger {
	CurLoglevel = int(newLogLevel)
	return &Logger{
		Logger:      log.New(out, prefix, flag),
		curLogLevel: newLogLevel,
	}
}

// Func Fatal will log fatal errors if enabled
func (l *Logger) Fatal(format string, v ...interface{}) {
	if CtsLog == nil {
		fmt.Printf("[FTL] %v\n", v)
	} else if l.curLogLevel >= FatalLevel {
		l.Output(2, fmt.Sprintf("[FTL] "+format, v...))
	}
}

// Func Error will log error logs or above
func (l *Logger) Error(format string, v ...interface{}) {
	if CtsLog == nil {
		fmt.Printf("[ERR] %v\n", v)
	} else if l.curLogLevel >= ErrorLevel {
		l.Output(2, fmt.Sprintf("[ERR] "+format, v...))
	}
}

// Func Warn will log Warning logs or above
func (l *Logger) Warn(format string, v ...interface{}) {
	if CtsLog == nil {
		fmt.Printf("[WRN] %v\n", v)
	} else if l.curLogLevel >= WarnLevel {
		l.Output(2, fmt.Sprintf("[WRN] "+format, v...))
	}
}

// Func Info will log Info logs or above
func (l *Logger) Info(format string, v ...interface{}) {
	if CtsLog == nil {
		fmt.Printf("[INF] %v\n", v)
	} else if l.curLogLevel >= InfoLevel {
		l.Output(2, fmt.Sprintf("[INF] "+format, v...))
	}
}

// Func Debug will log Debug logs or above
func (l *Logger) Debug(format string, v ...interface{}) {
	if CtsLog == nil {
		fmt.Printf("[DBG] %v\n", v)
	} else if l.curLogLevel >= DebugLevel {
		l.Output(2, fmt.Sprintf("[DBG] "+format, v...))
	}
}

var DefaultFilePath string = "."
var DefaultFileName string = "log"

// func InitCtsLogs initialize LogFileHandleLogs to file
// @param1  >> iniFilePath = path where the log will be stored if do not exists will be created
// @param2  >> iniFileName = AAAAMMDD_<IniFileName>.log
// @param3  >> delLogFiles = Delete Existing logs file on folder
// @param4  >> startLogMonitor = Start LogMonitor GoRotine
// @param4  >> maxLogFileSize = the MaxFile size per LogFile when reach make backup
// @return1 << int = 0 Success, -1 Fail to create dir, -2 fail to open Log file
// @return2 << err = nil success, other error
func InitCtsLogs(inFilePath string, inFileName string, delLogFiles bool, inLogLevel int, maxLogFileSize int) (int, error) {
	var err error = nil
	var erctr int = 0
	var usedFilePath string = inFilePath
	var usedFileName string = inFileName
	_, erroFolder := os.Stat(usedFilePath)
	if os.IsNotExist(erroFolder) {
		errDir := os.MkdirAll(usedFilePath, 0755)
		if errDir != nil {
			errDir := os.MkdirAll(usedFilePath, 0755)
			fmt.Printf("InitCtsLogs:FAIL Log Directory[%s] Not Created\n err[%s]\n", usedFilePath, errDir)
			erctr++
			usedFilePath = DefaultFilePath
			//return -1, errDir
		} else {
			fmt.Printf("InitCtsLogs:PASS Log Directory[%s] Successfull Created\n", usedFilePath)
		}
	}

	readDirectory, _ := os.Open(usedFilePath)
	allFiles, _ := readDirectory.Readdir(0)
	for f := range allFiles {
		var filePath string
		file := allFiles[f]
		fileName := file.Name()
		filePath = usedFilePath + "/" + fileName

		if strings.Contains(fileName, ".log") {
			if delLogFiles {
				os.Remove(filePath)
				fmt.Printf("InitCtsLogs:PASS filePath[%s] has been deleted\n", filePath)
			} else {
				fmt.Printf("InitCtsLogs:PASS filePath[%s] Not deleted\n", filePath)
			}
		}
	}
	if len(inFileName) == 0 {
		usedFileName = DefaultFileName
	}

	// Create LogFileHandle File Name: AAAAMMDD_LogFileHandle.log
	year, month, day := time.Now().Date()
	var strLogFileName string = ""
	if len(inFilePath) > 0 {
		strLogFileName = fmt.Sprintf("%s/%04d%02d%02d_%s.log", usedFilePath, year, int(month), day, usedFileName)
	} else {
		strLogFileName = fmt.Sprintf("%04d%02d%02d_%s.log", year, int(month), day, usedFileName)
	}

	/// Open New AAAAMMDD_LogFileHandle.log File
	CtsLogFileHandle, err = openLogFile(strLogFileName)
	if err != nil {
		fmt.Printf("InitCtsLogs:FAIL openLogFile(strLogFileName[%s])\n err[%s]\n", strLogFileName, err)
		erctr++
		//return -2, err
	}

	CtsLog = NewLogger(NewLogLevel(inLogLevel), CtsLogFileHandle, "", log.LstdFlags|log.Lshortfile|log.Ltime)
	CtsLog.SetOutput(io.MultiWriter(os.Stdout, CtsLogFileHandle))
	go LogFileSizeMonGoRoutine(CtsLogFileHandle, inFilePath, inFileName, maxLogFileSize, inLogLevel)
	return erctr, err
}

func LogFileSizeMonGoRoutine(logFileHandle *os.File, inFilePath string, inFileName string, maxLogFileSize int, inLogLevel int) {
	var i int = 0
	// Continuous lool Scan LogFile every minute
	for {
		// Wait one Minut before monitor file
		time.Sleep(1 * time.Minute)

		// Create LogFileHandle File Name: AAAAMMDD_LogFileHandle.log
		year, month, day := time.Now().Date()
		var strFullInLogFileName string = ""
		if len(inFilePath) > 0 {
			strFullInLogFileName = fmt.Sprintf("%s/%04d%02d%02d_%s.log", inFilePath, year, int(month), day, inFileName)
		} else {
			strFullInLogFileName = fmt.Sprintf("%04d%02d%02d_%s.log", year, int(month), day, inFileName)
		}

		// Get CtsLogFileHandle information
		fileInfo, err := CtsLogFileHandle.Stat()
		if err != nil {
			CtsLog.Error("Error getting file information:", err)
			return
		}

		// Get the file size in bytes
		fileSize := fileInfo.Size()

		// CloseCtsLogFileHandle
		// CtsLogFileHandle.Close()

		if fileSize > int64(maxLogFileSize) {
			// The current fileName MUST be backuped
			// Create LogFileHandle File Name: AAAAMMDD_LogFileHandle.log
			year, month, day := time.Now().Date()
			hour := time.Now().Hour()
			minute := time.Now().Minute()
			i = 0
			for {
				i++
				var strFullbkupLogFileName string = ""
				if len(inFilePath) > 0 {
					strFullbkupLogFileName = fmt.Sprintf("%s/%04d%02d%02d_%02dh%02dm_%s.log%03d", inFilePath, year, int(month), day, int(hour), int(minute), inFileName, i)
				} else {
					strFullbkupLogFileName = fmt.Sprintf("%04d%02d%02d_%02dh%02dm_%s.log%03d", year, int(month), day, int(hour), int(minute), inFileName, i)
				}
				_, err := os.Stat(strFullbkupLogFileName)
				if err == nil {
					CtsLog.Warn("strbkupLogFileName[%s]Already Exists Skipped!!!\n", strFullbkupLogFileName)
					continue // File exists
				}
				if os.IsNotExist(err) {

					// File does not exist move current file to it file and delete it
					// Attempt to move the file
					CtsLog.Warn("strFullInLogFileName[%s] Moving to strbkupLogFileName[%s] Started!!!\n\n", strFullInLogFileName, strFullbkupLogFileName)

					CtsLog.SetOutput(io.MultiWriter(os.Stdout))

					//Close LogFile and Rename LogFile
					err = CtsLogFileHandle.Close()
					if err != nil {
						CtsLog.Error("Close Feil \r\n strFullbkupLogFileName[%s] CtsLogFileHandle[%v]\r\n err[%s]\n\n", strFullbkupLogFileName, CtsLogFileHandle, err)
						break
					}
					CtsLog.Warn("Close Success \r\n strFullbkupLogFileName[%s] logFileHandle[%v]\r\n", strFullbkupLogFileName, logFileHandle)

					// Successfull closed renaming now
					err = os.Rename(strFullInLogFileName, strFullbkupLogFileName)
					if err != nil {
						// Fail to rename[move]
						CtsLog.Error("Rename\r\n strFullbkupLogFileName[%s] to \r\n strbkupLogFileName[%s]fail  err[%s]\n", strFullInLogFileName, strFullbkupLogFileName, err)
					} else {
						// Successfull renamed[moved] initialize again
						CtsLogFileHandle, err = openLogFile(strFullInLogFileName)
						if err != nil {
							CtsLog.Error("openLogFile Feil \r\n strFullbkupLogFileName[%s] CtsLogFileHandle[%v]\r\n err[%s]\n\n", strFullbkupLogFileName, CtsLogFileHandle, err)
						} else {
							// Setup LogFile again
							// CtsLog = NewLogger(NewLogLevel(inLogLevel), CtsLogFileHandle, "", log.LstdFlags|log.Lshortfile|log.Ltime)
							CtsLog.SetOutput(io.MultiWriter(os.Stdout, CtsLogFileHandle))
							CtsLog.Warn("openLogFile Success \r\n strFullbkupLogFileName[%s] logFileHandle[%v]\r\n", strFullbkupLogFileName, logFileHandle)
						}
					}
					break
				} // End backup file not exists so move current logFile to backup
			} // Rename File Loop
		} // End FileSize > MaxFileSize
	} // GoLang Loop
}

// openLogFile will try open the log file or create a new
func openLogFile(InFilePath string) (*os.File, error) {
	LogFileHandle, err := os.OpenFile(InFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("openLogFile:FAIL os.OpenFile(InFilePath[%s]) err[%s]\n", InFilePath, err)
		return nil, err
	}
	//fmt.Printf("openLogFile:PASS os.OpenFile(InFilePath[%s]) Successfull Created!!\n", InFilePath)
	return LogFileHandle, nil
}

// Example:
/*
func main() {
	rc, err := InitCtsLogs(DefaultFilePath, DefaultFileName)
	if (err != nil) || (rc < 0) {
		fmt.Printf("FAIL: InitCtsLogs rc[%d] err[%s]\n", rc, err)
		return
	}
	CtsLog.Debug("PASS: InitCtsLogs rc[%d]\n", rc)

	rc, levelName := CtsLog.GetLogLevel()
	CtsLog.Fatal("PASS: GetLogLevel()[%d=%s]\n", rc, levelName)
	CtsLog.Fatal("CtsLog.Fatal[%d]", 111)
	CtsLog.Error("CtsLog.Error[%d]", 222)
	CtsLog.Warn("CtsLog.Warn [%d]", 333)
	CtsLog.Info("CtsLog.Info [%d]", 444)
	CtsLog.Debug("CtsLog.Debug[%d]", 555)
	CtsLog.SetLogLevel(ErrorLevel)
	rc, levelName = CtsLog.GetLogLevel()
	CtsLog.Error("PASS: CtsLog.GetLogLevel()[%d=%s]\n", rc, levelName)
	CtsLog.Fatal("CtsLog.Fatal[%d]", 666)
	CtsLog.Error("CtsLog.Error[%d]", 777)
	CtsLog.Warn("CtsLog.Warn [%d]", 888)
	CtsLog.Info("CtsLog.Info [%d]", 999)
	CtsLog.Debug("CtsLog.Debug[%d]", 121)
}
*/
// =========================================================
// 2022/06/22 - OLIVEIRA: END NEW LOGGING MODE WITH LOG MODE
// =========================================================
