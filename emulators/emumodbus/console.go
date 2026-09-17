package main

import (
	utilsPkg "emumodbus/utils"
	logPkg "emumodbus/utils/gologtofile"
	"fmt"
	"net"
	"os"
	"strings"
)

// Console Terminal TCP Server
func ConsoleGoRotine() {
	// listen on TCP Server on Console Port
	ConsoleNetwork := fmt.Sprintf("%s:%d", utilsPkg.TcpSvrpAddress, utilsPkg.TcpConsolePort)
	consoleListener, err := net.Listen("tcp", ConsoleNetwork)
	if err != nil {
		resultCode := 0x80 //utilsPkg.NP_ERR_NET_LISTEN
		logPkg.CtsLog.Error("ConsoleGoRotine:resultCode[0x%X]ER listening [%v]", int(resultCode), err)
		if !GoUnitTestRunning {
			os.Exit(3)
		}
	}
	defer consoleListener.Close()
	for {
		logPkg.CtsLog.Warn("ConsoleGoRotine:CanOpenSvrpAddress[%s] is listening  CanOpenConsolePort:%d", utilsPkg.TcpSvrpAddress, utilsPkg.TcpConsolePort)
		// Accept incoming connections
		conn, err := consoleListener.Accept()
		if err != nil {
			resultCode := 0x81 // utilsPkg.NP_ERR_NET_LISTEN_ACCEPT
			logPkg.CtsLog.Warn("ConsoleGoRotine:resultCode[0x%X]ER listening [%v]", int(resultCode), err)
			continue
		}
		// Handle the console connection in a goroutine
		go handleConsoleConnection(conn)
	}
}

func handleConsoleConnection(conn net.Conn) {
	defer conn.Close()

	// Send a welcome message to the client
	welcomeMsg := "Welcome to the Emulator Console!\r\n Type Your Cmd Msg or 'quit' to exit.\r\n$ "
	conn.Write([]byte(welcomeMsg))

	fullCmd := ""
	buf := make([]byte, 1) // Read one byte at a time

	for {
		n, err := conn.Read(buf)
		if err != nil {
			logPkg.CtsLog.Warn("Client disconnected")
			return
		}

		char := string(buf[:n])

		if char == "\r" {
			// Ignore carriage return
			continue
		}

		if char == "\n" {
			// Process the command after receiving LF
			args := strings.Fields(fullCmd)
			if len(args) < 1 {
				logPkg.CtsLog.Warn("fullCmd: %s is Empty", fullCmd)
			}
			cmd := args[0]
			params := args[1:]
			nroParams := len(params)
			var rsp = ""
			switch cmd {
			case "quit":
				// Exiting force connection be closed
				conn.Write([]byte("Exiting!!!\n"))
				return
			case "start":
				// Proc Start Command
				if nroParams > 0 {
					if params[0] == "server" {
						if utilsPkg.ServerStopped || utilsPkg.StopServer {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Started..\r\n$ ", cmd, nroParams, params)
							utilsPkg.ServerStopped = false
							utilsPkg.StopServer = false
						} else {
							if !utilsPkg.StartServer {
								rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Started..\r\n$ ", cmd, nroParams, params)
							} else {
								rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Already Started..\r\n$ ", cmd, nroParams, params)
							}
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}

			case "stop":
				if nroParams > 0 {
					if params[0] == "server" {
						if utilsPkg.StopServer || utilsPkg.ServerStopped {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Already Stopped..\r\n$ ", cmd, nroParams, params)
						} else {
							if !utilsPkg.ServerStopped {
								utilsPkg.StopServer = true
							} else {
								rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Stopping..\r\n$ ", cmd, nroParams, params)
							}
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}

			// Proc Stop Command
			case "restart":
				if nroParams > 0 {
					if params[0] == "server" {
						// Proc Restart Command
						if utilsPkg.RestartServer {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Already Restarted..\r\n$ ", cmd, nroParams, params)
						} else {
							utilsPkg.RestartServer = true
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Restarted..\r\n$ ", cmd, nroParams, params)
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}

			case "block":
				if nroParams > 0 {
					if params[0] == "response" {
						if utilsPkg.BlockResponse {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Already Blocked..\r\n$ ", cmd, nroParams, params)
						} else {
							utilsPkg.BlockResponse = true
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Response Blocked..\r\n$ ", cmd, nroParams, params)
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}
			case "release":
				if nroParams > 0 {
					if params[0] == "response" {
						if utilsPkg.BlockResponse {
							utilsPkg.BlockResponse = true
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Response Released..\r\n$ ", cmd, nroParams, params)
						} else {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Response not Blocked..\r\n$ ", cmd, nroParams, params)
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid Arg\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}

			case "force":
				if nroParams > 0 {
					if params[0] == "crcerror" {
						if !utilsPkg.ForceCrcError {
							utilsPkg.ForceCrcError = true
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Forced CRC Error..\r\n$ ", cmd, nroParams, params)
						} else {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s CRC Error already forced..\r\n$ ", cmd, nroParams, params)
						}
					} else if params[0] == "crcok" {
						if utilsPkg.ForceCrcError {
							utilsPkg.ForceCrcError = false
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Forced CRC OK..\r\n$ ", cmd, nroParams, params)
						} else {
							rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s CRC already OK..\r\n$ ", cmd, nroParams, params)
						}
					} else {
						rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid Arg\r\n$ ", cmd, nroParams, params)
					}
				} else {
					rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Missing Args\r\n$ ", cmd, nroParams, params)
				}

			case "help":
				// Proc Restart Command
				rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Proccessed\r\n$ ", cmd, nroParams, params)
			default:
				rsp = fmt.Sprintf("\r\ncmd:[%s] nroParams[%d] params%s Invalid Command\r\n$ ", cmd, nroParams, params)
			}

			logPkg.CtsLog.Warn(">> fullCmd: %s", fullCmd)
			logPkg.CtsLog.Warn("<<     Rsp: %s", rsp)

			// Send generated response
			conn.Write([]byte(rsp))
			// Reset the command buffer to get next command
			fullCmd = ""
		} else {
			fullCmd += char
		}
	}
}
