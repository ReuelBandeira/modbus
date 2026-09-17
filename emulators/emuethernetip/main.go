package main

import (
	utilsPkg "emuethernetip/utils"
	logPkg "emuethernetip/utils/gologtofile"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

var EipSvrIpAddress string = "localhost"
var EipSvrPort uint16 = 44818
var EipConnectionId int = 0

var GoUnitTestRunning bool = false //Runnign go unit Tets

func main() {

	var resultCode byte = utilsPkg.OK

	utilsPkg.Hostname, _ = os.Hostname()
	// Initialize CtsLogs with default parameters
	rc, err := logPkg.InitCtsLogs(
		utilsPkg.LogFilePath,
		utilsPkg.LogFileName,
		utilsPkg.DeleteExistingLogFiles,
		utilsPkg.LogLevel,
		utilsPkg.LogFileMaxSize,
	)
	if (err != nil) || (rc < 0) {
		/// Fail to setup Logs
		/// TODO: Channel to Send Fail Information to Backend
		logPkg.CtsLog.Error("main:InitCtsLogs:Parameters rc[%d]\n            OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d] MaxSize[%d]\n         err[%s]\n",
			rc,
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			utilsPkg.LogFilePath,
			utilsPkg.LogFileName,
			utilsPkg.LogLevel,
			utilsPkg.LogFileMaxSize,
			err)
		if !GoUnitTestRunning {
			os.Exit(1)
		}
	}
	logPkg.CtsLog.Debug("main:InitCtsLogs:Parameters\n Running  OS[%s:%s] %d bits Host[%s]\n LogFilePath[%s]\n LogFileName[%s]\n    LogLevel[%d] MaxSize[%d]\n",
		runtime.GOOS,
		runtime.GOARCH,
		utilsPkg.OsBits,
		utilsPkg.Hostname,
		utilsPkg.LogFilePath,
		utilsPkg.LogFileName,
		utilsPkg.LogLevel,
		utilsPkg.LogFileMaxSize,
	)

	// Force Warning Log Levels Onpy
	utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(3)

	// Use net.Interfaces to get a list of all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		resultCode = utilsPkg.ERR_NET_INTERFACES
		logPkg.CtsLog.Error("resultCode[0x%X] err[%v]", int(resultCode), err)
		if !GoUnitTestRunning {
			os.Exit(2)
		}
	}

	logPkg.CtsLog.Debug("Connected IPv4 addresses:")

	for _, iface := range interfaces {
		// Check if the interface is up and not a loopback interface
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// Get the interface's addresses
			addrs, err := iface.Addrs()
			if err != nil {
				resultCode = utilsPkg.ERR_NET_IFC_ADDRS
				logPkg.CtsLog.Debug("resultCode[0x%X] err[%v]", int(resultCode), err)
				continue
			}

			// Iterate through the addresses associated with the interface
			for _, addr := range addrs {
				// Check if the address is an IPv4 address
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					resultCode = utilsPkg.ERR_NET_PARSE_CIDR
					logPkg.CtsLog.Debug("resultCode[0x%X] err[%v]", int(resultCode), err)
					continue
				}

				// Check if the address is an IPv4 address
				if ip.To4() != nil {
					logPkg.CtsLog.Info("%s: %s\n", iface.Name, ip)
					if iface.Name == "eth0" { // linux
						EipSvrIpAddress = ip.String()
						break
					} else if iface.Name == "wlan0" { // linux
						EipSvrIpAddress = ip.String()
						break
					} else if iface.Name == "Ethernet" { // windows
						EipSvrIpAddress = ip.String()
						break
					} else if iface.Name == "Wi-Fi" { // windows
						EipSvrIpAddress = ip.String()
						break
					}
				}
			}
		}
	}
	utilsPkg.TcpSvrIpAddress = EipSvrIpAddress
	utilsPkg.TcpSvrPort = EipSvrPort
	utilsPkg.TcpConsolePort = utilsPkg.TcpSvrPort + 1
	utilsPkg.ServerListenning = true
	utilsPkg.TcpClientConnEnabled = true

	// Start TCP Console Server
	go ConsoleGoRotine()

	for {
		err := RunTcpServer(utilsPkg.LogFileName, EipSvrIpAddress, EipSvrPort)
		if err != nil {
			if utilsPkg.ServerListenning {
				logPkg.CtsLog.Error("err[%v]", err)
				break
			}
		}
		logPkg.CtsLog.Warn(" Listen:Server[%s:%d] Closed utilsPkg.ServerListenning[%v]", EipSvrIpAddress, EipSvrPort, utilsPkg.ServerListenning)
		logPkg.CtsLog.Warn("Console:Server[%s:%d] Closed", EipSvrIpAddress, EipSvrPort+1)
		EipSvrIpAddress = utilsPkg.TcpSvrIpAddress
		EipSvrPort = utilsPkg.TcpSvrPort
		time.Sleep(10 * time.Second)
		if !utilsPkg.ServerListenning {
			utilsPkg.ServerListenning = true
			go ConsoleGoRotine()
		}
		logPkg.CtsLog.Warn(" Listen:Server[%s:%d] Restarted utilsPkg.ServerListenning[%v]", EipSvrIpAddress, EipSvrPort, utilsPkg.ServerListenning)
		logPkg.CtsLog.Warn("Console:Server[%s:%d] Restarted", EipSvrIpAddress, EipSvrPort+1)
	}
	if err != nil {
		logPkg.CtsLog.Error("Listen:Server[%s:%d] Exiting err[%v]", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpSvrPort, err)
	} else {
		logPkg.CtsLog.Warn("Listen:Server[%s:%d] Exiting", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpSvrPort)
	}
}

func RunTcpServer(protocolname string, serverAddress string, serverPort uint16) error {
	var err error = nil
	// Start listening on port 50000
	ServerNetwork := fmt.Sprintf("%s:%d", serverAddress, serverPort)
	utilsPkg.TcpServerListener, err = net.Listen("tcp", ServerNetwork)
	if err != nil {
		logPkg.CtsLog.Error(" %s:Server[%s:%d] CliId[%d] err[%v]", protocolname, serverAddress, serverPort, EipConnectionId, err)
		return err
	}
	defer utilsPkg.TcpServerListener.Close()

	utilsPkg.ServerListenning = true
	utilsPkg.TcpClientConnEnabled = true
	//
	// Server is Up, Listenning and waiting for Client Request Message
	//
	for utilsPkg.TcpServerRunning {
		if EipConnectionId == 0 || EipConnectionId == 4 || EipConnectionId == 12 {
			logPkg.CtsLog.Warn(" %s:Server[%s:%d] CliId[%d] Waiting for Client Connection", protocolname, serverAddress, serverPort, EipConnectionId)
		}
		if utilsPkg.TcpServerListener == nil {
			if !utilsPkg.TcpServerRunning {
				logPkg.CtsLog.Warn(" %s:Server[%s:%d] CliId[%d] listenner=nil Server not running", protocolname, serverAddress, serverPort, EipConnectionId)
			}
			break
		}
		// Accepting incoming connections
		conn, err := utilsPkg.TcpServerListener.Accept()
		if err != nil {
			if !utilsPkg.TcpServerRunning {
				break
			}
			logPkg.CtsLog.Error(" %s:Server[%s:%d] CliId[%d] err[%v]", protocolname, serverAddress, serverPort, EipConnectionId, err)
			continue
		}

		// After Client is connected Handle the each one in a goroutine
		if EipConnectionId == 0 || EipConnectionId == 4 || EipConnectionId == 12 {
			logPkg.CtsLog.Debug("%s:Client:[%s] SvrId[%d] Connected", protocolname, conn.RemoteAddr().String(), EipConnectionId)
		}
		go handleClientConnection(protocolname, conn, EipConnectionId)
		EipConnectionId++
	}
	logPkg.CtsLog.Warn(" %s:Server[%s:%d] CliId[%d] Exiting", protocolname, serverAddress, serverPort, EipConnectionId)
	return nil
}
