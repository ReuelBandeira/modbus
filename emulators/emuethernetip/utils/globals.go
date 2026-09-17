package utils

import "net"

// HostName
var Hostname string

// Operational Syatem 32 or 64 bits
const OsBits = 32 << (^uint(0) >> 63)

// LogFile parameters
var LogFilePath string = "."             // "."
var LogFileName string = "EmuEthernetIP" //  "EmuEthernetIP"
var LogLevel int = 4                     // "4=InfoLevel"
var LogFileMaxSize int = 500000          // 500Kbytes
var DelLogFiles bool = true              // "true"

// TCP Server Namd/IpAddress + Port
var TcpServerListener net.Listener = nil
var TcpServerRunning bool = true
var TcpSvrIpAddress string = "localhost"
var TcpSvrPort uint16 = 44818 // Emulate Ethernet/IP TCP Server Port
// TCP Console Server Port
var TcpConsolePort uint16 = TcpSvrPort + 1
var TcpClientConnId int = 0

// TCP Client Statistic
var TCPClientsConnected int = 0
var TCPClientRxPackets int = 0
var TCPClientTxPackets int = 0
var TCPClientRxBytes int = 0
var TCPClientTxBytes int = 0
var TCPClientRxErrors int = 0
var TCPClientTxErrors int = 0
var TCPClientConnErrors int = 0

// TCP Console Statistic
var TCPConsolesConnected int = 0
var TCPConsoleRxMsgs int = 0
var TCPConsoleTxMsgs int = 0
var TCPConsoleRxBytes int = 0
var TCPConsoleTxBytes int = 0
var TCPConsoleRxErrors int = 0
var TCPConsoleTxErrors int = 0
var TCPConsoleConnErrors int = 0

// TCP Console Cmd Flags
var ServerListenning = false
var TcpClientConnEnabled = false
var ServerStopped = false
var ResponseDisabled = false

var ForceCrcError = false
var CrcErrorForced = false
var AutoUpdate = true

// Constants
const (
	// ====================================
	// CIP EthernetIP Service Code Commands
	// ====================================
	// RegisterSession/UnRegisterSession Commands
	// These two commands are used to open and close an Encapsulation Session between two devices.
	// Once a Session is established, it is used to exchange more messages.
	//  Only one Session may exist between two devices.
	// The device receiving the RegisterSession request creates a Session Handle that it returns in
	// the RegisterSession reply.
	// This value is used to identify messages sent between the two devices that use this Session.
	CIP_EIP_SVC_CODE_REGISTER_SESSION = 0x0065 // Foward Open Request
	// The SendRRData Command is used for unconnected explicit messaging, and the SendUnitData Command is
	// used for connected explicit messaging. The device transmitting the SendRRData request creates a Sender
	// Context value that is returned with the reply.
	// The SendUnitData does not use the Sender Context field
	CIP_EIP_SVC_CODE_SEND_RR_DATA = 0x006F // Explicit Msg
	// =================================
	// CIP EthernetIP Data/Address Types
	// =================================
	CIP_EIP_DATA_TYPE_IO_MESSAGE       = 0x00B1
	CIP_EIP_DATA_TYPE_EXPLICIT_MESSAGE = 0x00B2
	CIP_EIP_ADDRESS_TYPE_SEQUENSED     = 0x8002
	// ==========================
	// CIP EthernetIP ResultCodes
	// ==========================
	OK = 0x00 // Success
	// CIP EthernetIP Gen Error ResultCodes
	ERR_BAD_PACKET_FORMAT           = 0x81
	ERR_BAD_PACKET_TOO_SHORT        = 0x82
	ERR_BAD_PACKET_TOO_LONG         = 0x83
	ERR_BAD_PACKET_DATA_LEN         = 0x84
	ERR_BAD_PACKET_DATA_LEN_TOO_BIG = 0x85
	ERR_BAD_CMD                     = 0x86
	// CIP EthernetIP Net Error ResultCodes
	ERR_NET_INTERFACES    = 0x92
	ERR_NET_IFC_ADDRS     = 0x93
	ERR_NET_PARSE_CIDR    = 0x94
	ERR_NET_LISTEN        = 0x95
	ERR_NET_LISTEN_ACCEPT = 0x96
	ERR_COMM_WRITE        = 0xA1
	ERR_COMM_READ         = 0xA2
	ERR_COMM_READ_TIMEPUT = 0xA3
)
