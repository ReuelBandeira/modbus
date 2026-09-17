package main

import (
	utilsPkg "emuethernetip/utils"
	logPkg "emuethernetip/utils/gologtofile"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var ConsoleNetwork string = fmt.Sprintf("%s:%d", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpConsolePort)
var ConsoleListener net.Listener = nil
var ConsoleConn net.Conn = nil

// Console Terminal TCP Server
func ConsoleGoRotine() {
	var err error = nil
	var clientConnectionId = 0

	// listen on TCP Server on Console Port
	ConsoleNetwork = fmt.Sprintf("%s:%d", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpConsolePort)
	ConsoleListener, err = net.Listen("tcp", ConsoleNetwork)
	if err != nil {
		resultCode := utilsPkg.ERR_NET_LISTEN
		logPkg.CtsLog.Error("   Console:Server:ResultCode[0x%X]ER CslId[%d]   err[%v]\n", int(resultCode), clientConnectionId, err)
		if !GoUnitTestRunning {
			os.Exit(3)
		}
	}
	defer ConsoleListener.Close()
	for utilsPkg.ServerListenning {
		logPkg.CtsLog.Warn("     Console:Server[%s:%d] CslId[%d] Waiting for Console Connection\n", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpConsolePort, clientConnectionId)
		if ConsoleListener == nil {
			break
		}
		// Accept incoming connections
		conn, err := ConsoleListener.Accept()
		if err != nil {
			if !utilsPkg.ServerListenning {
				break
			}
			resultCode := utilsPkg.ERR_NET_LISTEN_ACCEPT
			logPkg.CtsLog.Error("    Console:Server:ResultCode[0x%X]ER CslId[%d] err[%v]\n", int(resultCode), clientConnectionId, err)
			continue
		}
		// Handle the console connection in a goroutine
		logPkg.CtsLog.Warn("     Console:Client[%s] CslId[%d] Conneced\n", conn.RemoteAddr(), clientConnectionId)
		go handleConsoleConnection(conn, clientConnectionId)
		clientConnectionId++
	}
	logPkg.CtsLog.Warn("     Console:Server[%s:%d] CslId[%d] Exiting\n", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpConsolePort, clientConnectionId)
}

func handleConsoleConnection(conn net.Conn, clientConnectionId int) {
	defer conn.Close()
	var err error = nil
	var bIgnoreNextChar = false
	var numBytesRead int = 0
	if conn == nil {
		return
	}
	// Send a welcome message to the client
	welcomeMsg := "Welcome to " + utilsPkg.LogFileName + " Console!\r\n Type Your Cmd Msg or 'quit,exit,CTL-C' to exit.\r\n$ "
	conn.Write([]byte(welcomeMsg))
	utilsPkg.TCPConsolesConnected++

	fullUserInCmdRcvd := ""
	buf := make([]byte, 1) // Read one byte at a time
	for utilsPkg.ServerListenning {
		if conn == nil {
			break
		}
		numBytesRead, err = conn.Read(buf)
		if err != nil {
			utilsPkg.TCPConsoleRxErrors++
			logPkg.CtsLog.Error("   Console:Client[%s] CslId[%d] is Disconnected by Read Error err[%s]", conn.RemoteAddr(), clientConnectionId, err)
			if utilsPkg.TCPConsolesConnected > 0 {
				utilsPkg.TCPConsolesConnected--
			}
			return
		}
		utilsPkg.TCPConsoleRxBytes += numBytesRead

		char := string(buf[:numBytesRead])

		if char == "\r" {
			// Ignore carriage return
			continue
		}

		if char == "\n" {
			utilsPkg.TCPConsoleRxMsgs++
			// Split the fullUserInCmdRcvd string by commas
			// <cmd1,args1>,.. <cmdn,argsn> CRLF
			fullCmdsWithArgs := strings.Split(fullUserInCmdRcvd, ",")

			// Iterate through the values
			for _, singleCmdWithArgs := range fullCmdsWithArgs {
				// Proccess a Single Cmd + Args
				ConsoleRspMsg, err := proccessConsoleMsg(singleCmdWithArgs, clientConnectionId)
				if !utilsPkg.ServerListenning {
					if ConsoleConn != nil {
						logPkg.CtsLog.Warn("    Console:Client[%s] CslId[%d] ConsoleConn:closed utilsPkg.ServerListenning[%v] err[%v]", conn.RemoteAddr(), clientConnectionId, utilsPkg.ServerListenning, err)
						ConsoleConn.Close()
					}
					if ConsoleListener != nil {
						logPkg.CtsLog.Warn("    Console:Client[%s] CslId[%d] ConsoleListener:closed sutilsPkg.ServerListenning[%v] err[%v]", conn.RemoteAddr(), clientConnectionId, utilsPkg.ServerListenning, err)
						ConsoleListener.Close()
					}
					break
				}
				if err == nil {
					// PASS: Send Console Response
					if conn == nil {
						break
					}
					_, err = conn.Write([]byte(ConsoleRspMsg))
					if err != nil {
						utilsPkg.TCPConsoleTxErrors++
						if utilsPkg.TCPConsolesConnected > 0 {
							utilsPkg.TCPConsolesConnected--
						}
						logPkg.CtsLog.Error("    Console:Client[%s] CslId[%d] is Disconnected by Write Error err[%s]", conn.RemoteAddr(), clientConnectionId, err)
						return
					}
					utilsPkg.TCPConsoleTxMsgs++
					utilsPkg.TCPConsoleTxBytes += len(ConsoleRspMsg)

					if ConsoleRspMsg == "Exiting!!!" {
						// This command was exit or quit
						logPkg.CtsLog.Warn("    Console:Client[%s] CslId[%d] is Disconnected", conn.RemoteAddr(), clientConnectionId)
						if utilsPkg.TCPConsolesConnected > 0 {
							utilsPkg.TCPConsolesConnected--
						}
						return
					}
					if conn == nil {
						break
					}
					_, err = conn.Write([]byte("\r\n$ "))
					if err != nil {
						utilsPkg.TCPConsoleTxErrors++
						if utilsPkg.TCPConsolesConnected > 0 {
							utilsPkg.TCPConsolesConnected--
						}
						logPkg.CtsLog.Error("   Console:Client[%s] CslId[%d] is Disconnected by Write Error err[%s]", conn.RemoteAddr(), clientConnectionId, err)
						return
					}
					utilsPkg.TCPConsoleTxBytes += 4

				} else {
					// FAIL: Convert the error to a string
					if utilsPkg.TCPConsolesConnected > 0 {
						utilsPkg.TCPConsolesConnected--
					}
					logPkg.CtsLog.Error("   Console:Client[%s] CslId[%d] is Disconnected by Error err[%s]", conn.RemoteAddr(), err, clientConnectionId)
					return
				}
				if !utilsPkg.ServerListenning {
					if ConsoleConn != nil {
						ConsoleConn.Close()
					}
					if ConsoleListener != nil {
						ConsoleListener.Close()
					}
					break
				}
			}

			// Reset the fullUserInCmdRcvd buffer to get next userCmdLine
			fullUserInCmdRcvd = ""
		} else {
			// Storing char by char into fullUserInCmdRcvd while not CR+LF
			if char[0] == 0x1B {
				logPkg.CtsLog.Warn("   Console:Client id[%d] Typed char:%X rcvd 1B fullUserInCmdRcvd[%s]", clientConnectionId, char[0], fullUserInCmdRcvd)
				bIgnoreNextChar = true
			} else if char[0] == 0x1A || char[0] == 0x03 || char[0] == 0x08 {
				if char[0] == 0x08 {
					logPkg.CtsLog.Warn("   Console:Client id[%d]Typed char:[0x%02X]BS", clientConnectionId, char[0])
					if len(fullUserInCmdRcvd) > 0 {
						fullUserInCmdRcvd = fullUserInCmdRcvd[:len(fullUserInCmdRcvd)-1]
						if conn == nil {
							break
						}
						_, err = conn.Write([]byte("\r\n$ ")) // Space
						if err != nil {
							utilsPkg.TCPConsoleTxErrors++
							logPkg.CtsLog.Error("   Console:Client[%s] CslId[%d] is Disconnected by Write Error err[%s]", conn.RemoteAddr(), clientConnectionId, err)
							if utilsPkg.TCPConsolesConnected > 0 {
								utilsPkg.TCPConsolesConnected--
							}
							return
						}
						utilsPkg.TCPConsoleTxBytes += 4

						if conn == nil {
							break
						}
						_, err = conn.Write([]byte(fullUserInCmdRcvd))
						if err != nil {
							utilsPkg.TCPConsoleTxErrors++
							if utilsPkg.TCPConsolesConnected > 0 {
								utilsPkg.TCPConsolesConnected--
							}
							logPkg.CtsLog.Error("   Console:Client[%s] CslId[%d] is Disconnected by Write Error err[%s]", conn.RemoteAddr(), clientConnectionId, err)
							return
						}
						utilsPkg.TCPConsoleTxBytes += len(fullUserInCmdRcvd)
						utilsPkg.TCPConsoleTxMsgs++

					}
				} else if char[0] == 0x03 { // CTL-C
					// This command was exit or quit
					logPkg.CtsLog.Warn("    Console:Client[%s] CslId[%d] typed[0x%02X]CTL-C quit!!!", conn.RemoteAddr(), clientConnectionId, char[0])
					if conn == nil {
						break
					}
					conn.Write([]byte("Exiting!!!"))
					if utilsPkg.TCPConsolesConnected > 0 {
						utilsPkg.TCPConsolesConnected--
					}
					break
				} else {
					logPkg.CtsLog.Warn("    Console:Client id[%d]Typed char:[0x%02X]", clientConnectionId, char[0])
				}
			} else if char[0] < 0x20 { // ASCII Control Chars 0x00 to 0x1F
				logPkg.CtsLog.Warn("    Console:Client id[%d] Rcvd ASCII Control Chars:[0x%02X](0x00-0x1F)", clientConnectionId, char[0])
			} else {
				if bIgnoreNextChar {
					logPkg.CtsLog.Warn("    Console:Client id[%d] Rcvd ASCII char:[0x%02X] Ignored", clientConnectionId, char[0])
					if len(fullUserInCmdRcvd) > 0 {
						fullUserInCmdRcvd = fullUserInCmdRcvd[:len(fullUserInCmdRcvd)-1]
					}
					bIgnoreNextChar = false
				} else {
					fullUserInCmdRcvd += char
				}
			}
		}
	}
	if err == nil {
		// Normal exit
		logPkg.CtsLog.Warn("    Console:Client[%s] CslId[%d] Sucessfull exited!!!", conn.RemoteAddr(), clientConnectionId)
	} else {
		logPkg.CtsLog.Error("    Console:Client[%s] CslId[%d] Exit With Error[%s]", conn.RemoteAddr(), clientConnectionId, err)
	}
}

// UserCmd + It Args separated by space
func proccessConsoleMsg(ConsoleRquMsg string, clientConnectionId int) (string, error) {
	var consoleRspMsg = ""
	// Process the command after receiving LF
	args := strings.Fields(ConsoleRquMsg)
	if len(args) < 1 {
		logPkg.CtsLog.Warn(">> ConsoleRquMsg: %s cslId[%d] is Empty ", ConsoleRquMsg, clientConnectionId)
		return "Please type any command + params", nil
	}
	cmd := args[0]
	params := args[1:]
	nroParams := len(params)
	logPkg.CtsLog.Info(">>   ConsoleRquMsg[%s]cslId[%d] ", ConsoleRquMsg, clientConnectionId)
	logPkg.CtsLog.Info("                        cmd[%s]", cmd)
	logPkg.CtsLog.Info("                     params[%d] %s", nroParams, params)
	switch cmd {
	case "quit":
		// Exiting force connection be closed
		return "Exiting!!!", nil
	case "exit":
		// Exiting force connection be closed
		return "Exiting!!!", nil

	case "set":
		// set
		if nroParams > 0 && params[0] == "crcerror" {
			// set crcerror
			if nroParams > 1 && params[1] == "on" {
				// set crcerror on
				if !utilsPkg.ForceCrcError {
					utilsPkg.ForceCrcError = true
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Forced CRC Error..\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s CRC Error already forced..", cmd, nroParams, params)
				}
			} else if nroParams > 1 && params[1] == "off" {
				// set crcerror off
				if utilsPkg.ForceCrcError {
					utilsPkg.ForceCrcError = false
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Forced CRC OK..\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s CRC already OK..", cmd, nroParams, params)
				}
			} else {
				// set crcerror ? expected on/off
				consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Expected on/off", cmd, nroParams, params)
			}
		} else if nroParams > 0 && params[0] == "autoupdate" {
			// set autoupdate
			if nroParams > 1 && params[1] == "on" {
				// set autoupdate on
				if !utilsPkg.AutoUpdate {
					utilsPkg.AutoUpdate = true
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s ON\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Already on", cmd, nroParams, params)
				}
			} else if nroParams > 1 && params[1] == "off" {
				// set autoupdate off
				if utilsPkg.AutoUpdate {
					utilsPkg.AutoUpdate = false
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s off\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Already off", cmd, nroParams, params)
				}
			} else {
				// set autoupdate ? expected on/off
				consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Expected error on/off", cmd, nroParams, params)
			}
		} else if nroParams > 0 && params[0] == "response" {
			// set response
			if nroParams > 1 && params[1] == "enable" {
				// set response enable
				if utilsPkg.ResponseDisabled {
					utilsPkg.ResponseDisabled = false
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Response Released\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Response not Blocked..", cmd, nroParams, params)
				}
			} else if nroParams > 1 && params[1] == "disable" {
				// set response disable
				if utilsPkg.ResponseDisabled {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Already Blocked..", cmd, nroParams, params)
				} else {
					utilsPkg.ResponseDisabled = true
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Response Blocked\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				}
			} else {
				// set response expected enable/disable
				consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Expected enable/disable", cmd, nroParams, params)
			}
		} else if nroParams > 0 && params[0] == "clients" {
			// set clients
			if nroParams > 1 && params[1] == "enable" {
				// set clients enable
				if !utilsPkg.TcpClientConnEnabled {
					utilsPkg.ServerStopped = false
					utilsPkg.TcpClientConnEnabled = true
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Enabled\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Already Enabled..", cmd, nroParams, params)
				}
			} else if nroParams > 1 && params[1] == "disable" {
				// set clients disable
				if !utilsPkg.TcpClientConnEnabled {
					if utilsPkg.ServerStopped {
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Already Stopped..", cmd, nroParams, params)
					} else {
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Stop Pending\r\n", cmd, nroParams, params)
					}
				} else {
					utilsPkg.TcpClientConnEnabled = false
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s disabled\r\n", cmd, nroParams, params)
					consoleRspMsg += GetStrStatus()
				}
			} else {
				// set clients ? expected enable/disable
				consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Expected enable/disable", cmd, nroParams, params)
			}
		} else if nroParams > 0 && params[0] == "listen" {
			// set listen
			if nroParams > 1 && params[1] == "port" {
				// set listen port
				if nroParams > 2 {
					// set listen port
					if iNewListenPort, err := strconv.Atoi(params[2]); err == nil {
						// set listen port n
						if (iNewListenPort <= 0) || iNewListenPort > 65535 {
							// set listen port n(<=0 or >65565) invalid
							consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s invalid Expected 1 to 65535", cmd, nroParams, params)
						} else {
							// set listen port n(>=1 & <=65565) Valid
							// TODO: Try to set iSetListenPort
							// close any existing client connection and listenner then
							// try to listen on it new port
							utilsPkg.TcpSvrPort = uint16(iNewListenPort)
							utilsPkg.TcpConsolePort = utilsPkg.TcpSvrPort + 1
							utilsPkg.ServerListenning = false
							utilsPkg.TcpServerListener.Close()
							utilsPkg.TcpServerListener = nil
							consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s new Listen Port[%d]\r\n", cmd, nroParams, params, iNewListenPort)
							consoleRspMsg += GetStrStatus()
						}
					} else {
						// set listen port ?n(NotNumeric)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s ?n(NotNumeric)", cmd, nroParams, params)
					}
				} else {
					// set listen port expected nroParams > 2
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s expected nroParams > 2", cmd, nroParams, params)
				}
			} else {
				// set listen ? expected port n
				consoleRspMsg = fmt.Sprintf("cmd:[%s] expected nroParams[%d]>2 and params%s port n", cmd, nroParams, params)
			}
		} else if nroParams > 1 && params[0] == "loglevel" {
			// set loglevel
			if iSetLogLevel, err := strconv.Atoi(params[1]); err == nil {
				if iSetLogLevel != utilsPkg.LogLevel {
					// set loglevel n(new)
					// set loglevel n
					switch iSetLogLevel {
					case 0: //disabled
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.DisableLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Disabled\r\n", cmd, nroParams, params)
					case 1: //Fatal
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.FatalLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Fatal\r\n", cmd, nroParams, params)
					case 2: //Error
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.ErrorLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Error\r\n", cmd, nroParams, params)
					case 3: //Warning
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.WarnLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Warn\r\n", cmd, nroParams, params)
					case 4: //Info
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.InfoLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Info\r\n", cmd, nroParams, params)
					case 5: //Debug
						utilsPkg.LogLevel = logPkg.CtsLog.SetLogLevel(logPkg.DebugLevel)
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Debug\r\n", cmd, nroParams, params)
					default:
						// set loglevel n? Invalid
						consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s invalid logLenvel[%d]No Changed\r\n", cmd, nroParams, params, utilsPkg.LogLevel)
					}
					consoleRspMsg += GetStrStatus()
				} else {
					consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s No Changed\r\n", cmd, nroParams, params)
				}
			} else {
				// set loglevel ?
				consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Invalid", cmd, nroParams, params)
			}
		} else {
			// set ? expected params
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s expected parameters", cmd, nroParams, params)
		}

	case "reset":
		if nroParams > 0 && params[0] == "status" {
			ResetStatus()
			consoleRspMsg = GetStrStatus()
		} else if nroParams > 0 && params[0] == "statistics" {
			ResetStatistics()
			consoleRspMsg = GetStrStatistics()
		} else {
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d]>0 params%s default[reset status/statistic]\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		}

	case "show":
		if nroParams > 0 && params[0] == "status" {
			consoleRspMsg = GetStrStatus()
		} else if nroParams > 0 && params[0] == "statistics" {
			consoleRspMsg = GetStrStatistics()
		} else {
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d]>0 params%s default[show status]\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		}

	case "help":
		// help
		consoleRspMsg = GetStrUsage()
	default:
		// Test if cmd is a menu driven (0-99)
		switch cmd {
		case "1":
			// 1,2  set   clients     enable/disable  Client Connection
			utilsPkg.ServerStopped = false
			utilsPkg.TcpClientConnEnabled = true
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Enabled\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "2":
			// 1,2  set   clients     enable/disable  Client Connection
			utilsPkg.ServerStopped = false
			utilsPkg.TcpClientConnEnabled = false
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Disabled\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "3":
			// 3    set   listen      port n          Server Listen Port
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s logLenvel[%d]No Changed\r\n", cmd, nroParams, params, utilsPkg.LogLevel)
		case "4":
			// 4,5  set   response    enable/disable  Server response to Client
			utilsPkg.ResponseDisabled = true
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Response Released Enable\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "5":
			// 4,5  set   response    enable/disable  Server response to Client
			utilsPkg.ResponseDisabled = false
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Response Released Disable\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "6":
			// 6,7  set   autoupdate  on/off          Server response to Client
			utilsPkg.AutoUpdate = true
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s ON\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "7":
			// 6,7  set   autoupdate  on/off          Server response to Client
			utilsPkg.AutoUpdate = false
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s OFF\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "8":
			// 8,9  set   crcerror    on/off          force crcerror
			utilsPkg.ForceCrcError = true
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Forced CRC Error ON\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()
		case "9":
			// 8,9  set   crcerror    on/off          force crcerror
			utilsPkg.ForceCrcError = false
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s Forced CRC Error OFF\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrStatus()

		case "10":
			//10    set   loglevel    n               set   loglevel n(0-5)
			consoleRspMsg = fmt.Sprintf("cmd:[%s] nroParams[%d] params%s logLenvel[%d]No Changed\r\n", cmd, nroParams, params, utilsPkg.LogLevel)
		case "11":
			//11    reset status                      reset status
		case "12":
			//12    reset statistics                  reset statistics
			ResetStatistics()
			consoleRspMsg = GetStrStatistics()
		case "13":
			//13    show  status                      show  current status
		case "14":
			//14    show  statistics                  show  statistics
		case "98":
			//98    help                              show usage
			consoleRspMsg = GetStrUsage()
		case "0":
			//99    quit/exit/CTL-C                   exiting console
			// Exiting force connection be closed
			return "Exiting!!!", nil
		case "99":
			//99    quit/exit/CTL-C                   exiting console
			// Exiting force connection be closed
			return "Exiting!!!", nil
		default:
			// ? invalid command
			consoleRspMsg = fmt.Sprintf("cmd:[%s]Invalid nroParams[%d] params%s\n\n Valid Commands:\r\n", cmd, nroParams, params)
			consoleRspMsg += GetStrUsage()
		}
	}
	logPkg.CtsLog.Debug("cslId[%d] << ConsoleRspMsg:\r\n%s\r\n", clientConnectionId, consoleRspMsg)
	return consoleRspMsg, nil
}

func GetStrUsage() string {
	strUsage := "\r\n Usage:\r\n"
	strUsage += "\r\n"
	strUsage += "    <cmd1>  [arg1.1 arg1.x],.....<cmdn> [argn.1 argn.x]\r\n"
	strUsage += "\r\n"
	strUsage += "Option:\r\n"
	strUsage += " 1,2  set   clients     enable/disable  Client Connection\r\n"
	strUsage += " 3    set   listen      port n          Server Listen Port\r\n"
	strUsage += " 4,5  set   response    enable/disable  Server response to Client\r\n"
	strUsage += " 6,7  set   autoupdate  on/off          Server response to Client\r\n"
	strUsage += " 8,9  set   crcerror    on/off          force crcerror\r\n"
	strUsage += "10    set   loglevel    n               set   loglevel n(0-5)\r\n"
	strUsage += "11    reset status                      reset status\r\n"
	strUsage += "12    reset statistics                  reset statistics\r\n"
	strUsage += "13    show  status                      show  current status\r\n"
	strUsage += "14    show  statistics                  show  statistics\r\n"
	strUsage += "98    help                              show usage\r\n"
	strUsage += "0,99  quit/exit/CTL-C                   exiting console\r\n"
	return strUsage
}

func GetStrStatus() string {
	nLogLevel, strLogLevel := logPkg.CtsLog.GetLogLevel()
	strStatus := fmt.Sprintf("ConsoleIpAddressP:ort     [%s:%d]\r\n", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpConsolePort)
	strStatus += fmt.Sprintf(" TcpSvrIpAddress:Port     [%s:%d]\r\n", utilsPkg.TcpSvrIpAddress, utilsPkg.TcpSvrPort)
	strStatus += fmt.Sprintf("flag:ServerListenning     [%v]\r\n", utilsPkg.ServerListenning)
	strStatus += fmt.Sprintf("flag:TcpClientConnId      [%d]\r\n", utilsPkg.TcpClientConnId)
	strStatus += fmt.Sprintf("flag:TcpClientConnEnabled [%v]\r\n", utilsPkg.TcpClientConnEnabled)
	strStatus += fmt.Sprintf("flag:ServerStopped        [%v]\r\n", utilsPkg.ServerStopped)
	strStatus += fmt.Sprintf("flag:ResponseDisabled     [%v]\r\n", utilsPkg.ResponseDisabled)
	strStatus += fmt.Sprintf("flag:ForceCrcError        [%v]\r\n", utilsPkg.ForceCrcError)
	strStatus += fmt.Sprintf("flag:AutoUpdate           [%v]\r\n", utilsPkg.AutoUpdate)
	strStatus += fmt.Sprintf("flag:LogLevel             [%d]%s\r\n", nLogLevel, strLogLevel)
	return strStatus
}

func GetStrStatistics() string {
	strStatistic := fmt.Sprintf("TCPClientsConnected  [%d]\r\n", utilsPkg.TCPClientsConnected)
	strStatistic += fmt.Sprintf("TCPClientRxPackets   [%d]\r\n", utilsPkg.TCPClientRxPackets)
	strStatistic += fmt.Sprintf("TCPClientTxPackets   [%d]\r\n", utilsPkg.TCPClientTxPackets)
	strStatistic += fmt.Sprintf("TCPClientRxBytes     [%d]\r\n", utilsPkg.TCPClientRxBytes)
	strStatistic += fmt.Sprintf("TCPClientTxBytes     [%d]\r\n", utilsPkg.TCPClientTxBytes)
	strStatistic += fmt.Sprintf("TCPConsolesConnected [%d]\r\n", utilsPkg.TCPConsolesConnected)
	strStatistic += fmt.Sprintf("TCPConsoleRxMsgs     [%d]\r\n", utilsPkg.TCPConsoleRxMsgs)
	strStatistic += fmt.Sprintf("TCPConsoleTxMsgs     [%d]\r\n", utilsPkg.TCPConsoleTxMsgs+1)
	strStatistic += fmt.Sprintf("TCPConsoleRxBytes    [%d]\r\n", utilsPkg.TCPConsoleRxBytes)
	strStatistic += fmt.Sprintf("TCPConsoleTxBytes    [%d]\r\n", utilsPkg.TCPConsoleTxBytes)
	strStatistic += fmt.Sprintf("TCPClientConnErrors  [%d]\r\n", utilsPkg.TCPClientConnErrors)
	strStatistic += fmt.Sprintf("TCPClientRxErrors    [%d]\r\n", utilsPkg.TCPClientRxErrors)
	strStatistic += fmt.Sprintf("TCPClientTxErrors    [%d]\r\n", utilsPkg.TCPClientTxErrors)
	strStatistic += fmt.Sprintf("TCPConsoleConnErrors [%d]\r\n", utilsPkg.TCPClientConnErrors)
	strStatistic += fmt.Sprintf("TCPConsoleRxErrors   [%d]\r\n", utilsPkg.TCPConsoleRxErrors)
	strStatistic += fmt.Sprintf("TCPConsoleTxErrors   [%d]\r\n", utilsPkg.TCPConsoleTxErrors)
	return strStatistic
}

func ResetStatus() {
	utilsPkg.TcpClientConnEnabled = true
	utilsPkg.ServerStopped = false
	utilsPkg.ResponseDisabled = false
	utilsPkg.ForceCrcError = false
	utilsPkg.AutoUpdate = true
}

func ResetStatistics() {
	// reset TCP Client Statistic
	// utilsPkg.TCPClientsConnected = 0
	utilsPkg.TCPClientConnErrors = 0
	utilsPkg.TCPClientRxPackets = 0
	utilsPkg.TCPClientTxPackets = 0
	utilsPkg.TCPClientRxBytes = 0
	utilsPkg.TCPClientTxBytes = 0
	utilsPkg.TCPClientRxErrors = 0
	utilsPkg.TCPClientTxErrors = 0
	// reset TCP Console Statistic
	// utilsPkg.TCPConsolesConnected = 0
	utilsPkg.TCPConsoleConnErrors = 0
	utilsPkg.TCPConsoleRxMsgs = 0
	utilsPkg.TCPConsoleTxMsgs = 0
	utilsPkg.TCPConsoleRxBytes = 0
	utilsPkg.TCPConsoleTxBytes = 0
	utilsPkg.TCPConsoleRxErrors = 0
	utilsPkg.TCPConsoleTxErrors = 0
}
