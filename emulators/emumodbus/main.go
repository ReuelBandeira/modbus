package main

import (
	"crypto/tls"  // Modbus Secure TCP Server TLS
	"crypto/x509" // Modbus Secure TCP Server TLS
	utilsPkg "emumodbus/utils"
	cfgPkg "emumodbus/utils/config"
	logPkg "emumodbus/utils/gologtofile"
	"fmt"
	"net"
	"runtime"

	"os"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
)

var GoUnitTestRunning bool = false

const (
	MINUS_ONE int16 = -1
)

// To build for Linux on Windows type
// set GOOS=linux
// set GOARCH=amd64
// go build -o emumodbus

/* Modbus TCP+TLS (MBAPS or Modbus Security) server example.
 *
 * This file is intended to be a demo of the modbus server in a tcp+tls
 * configuration.
 * It shows how to configure and start a server, as well as how to use
 * client roles to perform authorization in the handler.
 * Feel free to use it as boilerplate for simple servers.
 *
 * This server simulates a simple wall clock device, exposing a 32-bit unix
 * timestamp in holding registers #0 and 1.
 * The timestamp is incremented every second by the main loop.
 *
 * Access control is done by way of Modbus Roles, which are encoded in the
 * client certificate as an X509 extension:
 * - any client can read the clock regardless of their role, provided that their
 *   certificate is accepted by the server,
 * - only clients with the "operator" role specified in their certificate can
 *   set the time.
 *
 * Certificates with no, invalid or multiple Modbus Role extensions will have
 * their role set to an empty string (req.ClientRole == "").
 *
 * Requests from clients with certificates not passing TLS verification are
 * rejected at the TLS layer (i.e. before reaching the Modbus layer).
 *
 *
 * The following commands can be used to create self-signed server and client
 * certificates:
 * $ mkdir certs
 *
 * on Linux: create the server key pair:
 * $ openssl req -x509 -newkey rsa:4096 -sha256 -days 360 -nodes \
 *   -keyout certs/server.key.pem -out certs/server.cert.pem \
 *   -subj "/CN=TEST SERVER CERT DO NOT USE/" -addext "subjectAltName=DNS:localhost" \
 *   -addext "keyUsage=keyCertSign,digitalSignature,keyEncipherment" \
 *   -addext "extendedKeyUsage=critical,serverAuth"
 *
 * on Linux:create a client certificate with the "user" role:
 * $ openssl req -x509 -newkey rsa:4096 -sha256 -days 360 -nodes \
 *   -keyout certs/client.key.pem -out certs/client.cert.pem \
 *   -subj "/CN=TEST CLIENT CERT DO NOT USE/" \
 *   -addext "keyUsage=keyCertSign,digitalSignature,keyEncipherment" \
 *   -addext "extendedKeyUsage=critical,clientAuth" \
 *   -addext "1.3.6.1.4.1.50316.802.1=ASN1:UTF8String:user"
 *
 * on Linux:create another client certificate with the "operator" role:
 * $ openssl req -x509 -newkey rsa:4096 -sha256 -days 360 -nodes \
 *   -keyout certs/operator-client.key.pem -out certs/operator-client.cert.pem \
 *   -subj "/CN=TEST CLIENT CERT DO NOT USE/" \
 *   -addext "keyUsage=keyCertSign,digitalSignature,keyEncipherment" \
 *   -addext "extendedKeyUsage=critical,clientAuth" \
 *   -addext "1.3.6.1.4.1.50316.802.1=ASN1:UTF8String:operator"
 *
 * on Linux:create a file containing both client certificates (for use by the server as an
 * 'allowed client list'):
 * $ cat certs/client.cert.pem certs/operator-client.cert.pem >certs/clients.cert.pem
 *
 * start the server:
 * $ go run examples/tls_server.go
 *
 * in another shell, read the clock with modbus-cli as the 'user' role:
 * $ go run cmd/modbus-cli.go --target tcp+tls://localhost:5802 --cert certs/client.cert.pem \
 *   --key certs/client.key.pem --ca certs/server.cert.pem rh:uint32:0
 *
 * attempting to set the clock as 'user' should fail with Illegal Function:
 * $ go run cmd/modbus-cli.go --target tcp+tls://localhost:5802 --cert certs/client.cert.pem \
 *   --key certs/client.key.pem --ca certs/server.cert.pem wr:uint32:0:1598692358
 *
 * setting the clock as 'operator' should succeed:
 * $ go run cmd/modbus-cli.go --target tcp+tls://localhost:5802 --cert certs/operator-client.cert.pem \
 *   --key certs/operator-client.key.pem --ca certs/server.cert.pem wr:uint32:0:1598692358
 *
 * reading the cock as 'operator' should also work:
 * $ go run cmd/modbus-cli.go --target tcp+tls://localhost:5802 --cert certs/operator-client.cert.pem \
 *   --key certs/operator-client.key.pem --ca certs/server.cert.pem rh:uint32:0
 */

/*
* Simple modbus server example.
*
* This file is intended to be a demo of the modbus server.
* It shows how to create and start a server, as well as how
* to write a handler object.
* Feel free to use it as boilerplate for simple servers.
 */

// run this with go run examples/tcp_server.go
func main() {

	// Get Host Name
	utilsPkg.Hostname, _ = os.Hostname()

	// Initialize CtsLogs with default parameters
	rc, err := logPkg.InitCtsLogs(
		cfgPkg.DefaultLogFilePath,
		cfgPkg.DefaultLogFileName,
		cfgPkg.DeleteExistingLogFiles,
		cfgPkg.DefaultLogLevel,
	)
	if (err != nil) || (rc < 0) {
		/// Fail to setup Logs
		logPkg.CtsLog.Error("InitCtsLogs:Parameters rc[%d]\n                         OS[%s:%s] %d bits Host[%s]\n              LogFilePath[%s]\n              LogFileName[%s]\n                 LogLevel[%d]\n                   err[%s]\n\n",
			rc,
			runtime.GOOS,
			runtime.GOARCH,
			utilsPkg.OsBits,
			utilsPkg.Hostname,
			cfgPkg.DefaultLogFilePath,
			cfgPkg.DefaultLogFileName,
			cfgPkg.DefaultLogLevel,
			err)
		return
	}

	// Logs successfull Initialized
	logPkg.CtsLog.Debug("InitCtsLogs:Parameters\n Running               OS[%s:%s] %d bits Host[%s]\n              LogFilePath[%s]\n              LogFileName[%s]\n                 LogLevel[%d]\n\n",
		runtime.GOOS,
		runtime.GOARCH,
		utilsPkg.OsBits,
		utilsPkg.Hostname,
		cfgPkg.DefaultLogFilePath,
		cfgPkg.DefaultLogFileName,
		cfgPkg.DefaultLogLevel)

	// Set Log Level 3 = Warning  level 5 = Debug
	rc = logPkg.CtsLog.SetLogLevel(3)
	if rc < 0 {
		logPkg.CtsLog.Error("SetLogLevel(3)=Warning Faled\n")
	} else {
		logPkg.CtsLog.Debug("SetLogLevel(3)=Warning Success\n")
	}

	// ===================================================
	// BEGIN Modbus : TCP Server Optaining Host Ip Address
	// ===================================================
	var tcpSvrErr error
	var mbTcpSever *modbus.ModbusServer
	var mbTcpSvrEh *mbTcpSvrEventHandler
	var ticker *time.Ticker

	// create the handler object
	mbTcpSvrEh = &mbTcpSvrEventHandler{}
	logPkg.CtsLog.Debug("mbTcpSvrEventHandler:eh.uptime[%d]\n", mbTcpSvrEh.uptime)

	// BEGIN: GET HOST IP ADDRESS
	var MbSvrIpAddress string = "localhost" // Default Host Name localhost(127.0.0.1)
	var MbSvrPort int = 502                 // Default Modbus Port
	var MbTlsSvrPort int = 802              // Default Modbus TLS Port

	// Use net.Interfaces to get a list of all Host network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		logPkg.CtsLog.Error(" err[%v]", err)
		os.Exit(2)
	}

	logPkg.CtsLog.Debug("Connected IPv4 addresses:\n")

	for _, iface := range interfaces {
		// Check if the interface is up and not a loopback interface
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// Get the interface's addresses
			addrs, err := iface.Addrs()
			if err != nil {
				logPkg.CtsLog.Error(" err[%v]", err)
				continue
			}

			// Iterate through the addresses associated with the interface
			for _, addr := range addrs {
				// Check if the address is an IPv4 address
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					logPkg.CtsLog.Debug(" err[%v]", err)
					continue
				}

				// List all IPv4 address
				if ip.To4() != nil {
					logPkg.CtsLog.Debug("iface Name[%s]: ip[%s]\n", iface.Name, ip)
					if iface.Name == "eth0" { // linux
						MbSvrIpAddress = ip.String()
						break
					} else if iface.Name == "wlan0" { // linux
						MbSvrIpAddress = ip.String()
						break
					} else if iface.Name == "Ethernet" { // windows
						MbSvrIpAddress = ip.String()
						break
					} else if iface.Name == "Wi-Fi" { // windows
						MbSvrIpAddress = ip.String()
						break
					}
				}
			}
		}
	}

	// Start TCP Console Server
	go ConsoleGoRotine()

	// Start listening on port 502
	strMbTcpServerURL := fmt.Sprintf("tcp://%s:%d", MbSvrIpAddress, MbSvrPort)

	// --END: GET HOST IP ADDRESS
	// var strMbTcpServerURL string = "tcp://localhost:5502"

	var maxClients uint = 5

	var mbTcpServerConf = modbus.ServerConfiguration{
		URL: strMbTcpServerURL,
		// close idle connections after 30s of inactivity
		Timeout: 30 * time.Second,
		// accept 5 concurrent connections max.
		MaxClients: maxClients,
	}

	// create the mbTcpSever object
	mbTcpSever, tcpSvrErr = modbus.NewServer(&mbTcpServerConf, mbTcpSvrEh)
	if tcpSvrErr != nil {
		logPkg.CtsLog.Error("modbus.NewServer:mbTcpServerConf\n\t  tcpSvrErr[%s]\n\t        URL[%s]\n\t              Timeout[%d]seconds\n\t           MaxClients[%d]\n", tcpSvrErr, mbTcpServerConf.URL, mbTcpServerConf.Timeout/1000000000, mbTcpServerConf.MaxClients)
		os.Exit(1)
	}
	logPkg.CtsLog.Debug("modbus.NewServer:mbTcpServerConf\n\t                  URL[%s]\n\t              Timeout[%d]seconds\n\t           MaxClients[%d]\n", mbTcpServerConf.URL, mbTcpServerConf.Timeout/1000000000, mbTcpServerConf.MaxClients)

	// start accepting client connections
	// note that Start() returns as soon as the mbTcpSever is started
	tcpSvrErr = mbTcpSever.Start()
	if tcpSvrErr != nil {
		logPkg.CtsLog.Error("mbTcpSever.Start(): mbTcpServerConf.URL[%s] mbTcpServerConf.Timeout[%d]seconds mbTcpServerConf.MaxClients[%d] err[%s]\n", mbTcpServerConf.URL, mbTcpServerConf.Timeout/1000000000, mbTcpServerConf.MaxClients, tcpSvrErr)
		os.Exit(1)
	}
	logPkg.CtsLog.Debug(" server.Start(): Started URL[%s]Clean Server Listenning Timeout[%d]seconds MaxClients[%d]\n", mbTcpServerConf.URL, mbTcpServerConf.Timeout/1000000000, mbTcpServerConf.MaxClients)
	// =======================
	// END Modbus : TCP Server
	// =======================

	// =====================================
	// BEGIN Modbus : Secure TCP Server(TLS)
	// =====================================
	var mbTlsSvrErr error
	var mbTlsSvrEh *mbTlsServerEventHandler
	var mbTlsServer *modbus.ModbusServer
	var mbTlsServerKeyPair tls.Certificate
	var mbTlsClientCertPool *x509.CertPool

	var strMbTlsServerURL = fmt.Sprintf("tcp+tls://%s:%d", MbSvrIpAddress, MbTlsSvrPort)
	// var strMbTlsServerURL string = "tcp+tls://localhost:5802"
	var strMbTlsClientCertificate string = "certs/client.cert.pem"
	var strMbTlsServerCertificate string = "certs/server.cert.pem"
	var strMbTlsServerPrivateKey string = "certs/server.key.pem"
	// create the handler object
	mbTlsSvrEh = &mbTlsServerEventHandler{}

	// load the tlsServer certificate and its associated private key, which
	// are used to authenticate the tlsServer to the client.
	// note that a tls.Certificate object can contain both the cert and its key,
	// which is the case here.
	mbTlsServerKeyPair, mbTlsSvrErr = tls.LoadX509KeyPair(
		strMbTlsServerCertificate, strMbTlsServerPrivateKey)
	if mbTlsSvrErr != nil {
		logPkg.CtsLog.Error("tls.LoadX509KeyPair:\n\t tlsServerCertificate[%s]\n\t  tlsServerprivateKey[%s]\n\t                Error[%s]\n", strMbTlsServerCertificate, strMbTlsServerPrivateKey, mbTlsSvrErr)
	} else {
		logPkg.CtsLog.Debug(" tlsServerCertificate[%s]SVR_CERT_PUK\n", strMbTlsServerCertificate)
		logPkg.CtsLog.Debug(" tlsServerprivateKey [%s]SVR_PRIVK\n", strMbTlsServerPrivateKey)
		// load TLS client authentication material, which could either be:
		// - the CA (Certificate Authority) certificate(s) used to sign client certs,
		// - the list of allowed client certs, if client certificates are self-signed or
		//   if client certificate pinning is required.
		mbTlsClientCertPool, mbTlsSvrErr = modbus.LoadCertPool(strMbTlsClientCertificate)
		if mbTlsSvrErr != nil {
			logPkg.CtsLog.Error("modbus.LoadCertPool CA/client\n\t   Clientcertificates[%s]\n\t Error[%s]\n", strMbTlsClientCertificate, mbTlsSvrErr)
		} else {
			logPkg.CtsLog.Debug("   Clientcertificates[%s]CLI_CERT_PUKS\n\n", strMbTlsClientCertificate)
			// create the tlsSegorver object
			mbTlsServer, mbTlsSvrErr = modbus.NewServer(&modbus.ServerConfiguration{
				// listen on localhost port 5802
				URL: strMbTlsServerURL,
				// accept 10 concurrent connections max.
				MaxClients: 10,
				// close idle connections after 1min of inactivity
				Timeout: 60 * time.Second,
				// use tlsServerKeyPair as tlsServer certificate + tlsServer private key
				TLSServerCert: &mbTlsServerKeyPair,
				// use the client cert/CA pool to verify client certificates
				TLSClientCAs: mbTlsClientCertPool,
			}, mbTlsSvrEh)
			if mbTlsSvrErr != nil {
				logPkg.CtsLog.Error("modbus.NewSever:FAIL to create\n                tlsServer[%s] Error[%s]\n", strMbTlsServerURL, mbTlsSvrErr)
			} else {
				// start accepting client connections
				// note that Start() returns as soon as the tlsServer is started
				mbTlsSvrErr = mbTlsServer.Start()
				if mbTlsSvrErr != nil {
					logPkg.CtsLog.Error("modbus.NewSever:FAIL to start\n                tlsServer[%s] Error[%s]\n", strMbTlsServerURL, mbTlsSvrErr)
				}
				logPkg.CtsLog.Debug("modbus>NewServer:PASS Started  tlsServer[%s]Secure Server Listenning\n\n", strMbTlsServerURL)
			}
		}
	}
	logPkg.CtsLog.Warn(" Modbus:Server[%s](TLS)Secure Communication Listenning\n", strMbTlsServerURL)
	logPkg.CtsLog.Warn(" Modbus:Server[    %s] Non-Secure Communication Listenning\n", mbTcpServerConf.URL)

	// increment a 32-bit uptime counter every second.
	// (this counter is exposed as input registers 200-201 for demo purposes)
	ticker = time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		// since the handler methods are called from multiple goroutines,
		// use locking where appropriate to avoid concurrency issues.
		mbTcpSvrEh.lock.Lock()
		mbTcpSvrEh.uptime++
		mbTcpSvrEh.lock.Unlock()
	}

	// never reached
	// return
}

// ==============================
// Modbus TCP Server Event Hander
// ==============================
// Handler object, passed to the NewServer() constructor above.
type mbTcpSvrEventHandler struct {
	// this lock is used to avoid concurrency issues between goroutines, as
	// handler methods are called from different goroutines
	// (1 goroutine per client)
	lock sync.RWMutex

	// simple uptime counter, incremented in the main() above and exposed
	// as a 32-bit input register (2 consecutive 16-bit modbus registers).
	uptime uint32

	// these are here to hold client-provided (written) values, for both coils and
	// holding registers

	// Emulating Modbus Registers
	coils            [65535]bool   // DO: Digital Output
	discreteInputs   [65535]bool   // DI: Digital Input  Read Only
	holdingRegisters [65535]uint16 // AO: Anoloog Output
	inputRegisters   [65535]uint16 // AI: Analog Input   Read Only
	holdingReg1      uint16        // this is a 16-bit unsigned integer
	holdingReg2      uint16        // this is a 16-bit unsigned integer
	holdingReg3      int16         // this is a 16-bit signed integer
	holdingReg4      uint32        // this is a 32-bit unsigned integer

}

// Modbus:Coils:RW  DO(DigitalOutputs) true/false
func (eh *mbTcpSvrEventHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	// since we're manipulating variables shared between multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	eh.lock.Lock()
	// release the lock upon return
	defer eh.lock.Unlock()

	if req.UnitId > 31 {
		// note: we're merely filtering here, but we could as well use the unit
		// ID field to support multiple register maps in a single server.
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleCoils: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(uint16(len(eh.coils))) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleCoils: req.Addr[%d] >=uint16(uint16(len(eh.coils))=[%d] err[%s]\n", req.Addr, uint16(len(eh.coils)), err)
		return
	}

	// make sure that all registers covered by this request actually exist
	if int(req.Addr)+int(req.Quantity) > len(eh.coils) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleCoils: int(req.Addr)+int(req.Quantity)[%d] > len(eh.coils)[%d] err[%s]\n", int(req.Addr)+int(req.Quantity), len(eh.coils), err)
		return
	}

	// loop through `req.Quantity` registers, from address `req.Addr` to
	// `req.Addr + req.Quantity - 1`, which here is conveniently `req.Addr + i`
	for i := 0; i < int(req.Quantity); i++ {
		// ignore the write if the current register address is 80
		if req.IsWrite && int(req.Addr)+i != 80 {
			// assign the value
			if len(req.Args) > i {
				eh.coils[int(req.Addr)+i] = req.Args[i]
			}
		}
		// append the value of the requested register to res so they can be and interlock coil
		// sent back to the client
		res = append(res, eh.coils[int(req.Addr)+i])
		if eh.coils[int(req.Addr)+i] {
			eh.coils[int(req.Addr)+i] = false
		} else {
			eh.coils[int(req.Addr)+i] = true
		}
		switch req.Addr {
		case 0x100:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		case 0x101:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		case 0x102:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		case 0x103:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		case 0x104:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		case 0x105:
			logPkg.CtsLog.Warn("HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, eh.coils[int(req.Addr)+i])
		}

	}
	return
}

// Modbus:DiscreteInputs:RO  DI(DigitalInputs) true/false
func (eh *mbTcpSvrEventHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {

	// acquire a lock to avoid concurrency issues.
	eh.lock.Lock()
	// release the lock upon return
	defer eh.lock.Unlock()

	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleDiscreteInputs: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(uint16(len(eh.discreteInputs))) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleDiscreteInputs: req.Addr[%d] >=uint16(uint16(len(eh.discreteInputs))=[%d] err[%s]\n", req.Addr, uint16(len(eh.discreteInputs)), err)
		return
	}

	// make sure that all registers covered by this request actually exist
	if int(req.Addr)+int(req.Quantity) > len(eh.discreteInputs) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleDiscreteInputs: int(req.Addr)+int(req.Quantity)[%d] > len(eh.discreteInputs)[%d] err[%s]\n", int(req.Addr)+int(req.Quantity), len(eh.discreteInputs), err)
		return
	}

	// loop through `req.Quantity` registers, from address `req.Addr` to
	// `req.Addr + req.Quantity - 1`, which here is conveniently `req.Addr + i`
	for i := 0; i < int(req.Quantity); i++ {
		// sent back to the client and interlock it discrete input value
		res = append(res, eh.discreteInputs[int(req.Addr)+i])
		if eh.discreteInputs[int(req.Addr)+i] {
			eh.discreteInputs[int(req.Addr)+i] = false
		} else {
			eh.discreteInputs[int(req.Addr)+i] = true
		}
		switch req.Addr {
		case 0x100:
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x101:
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x102:
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x103:
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x104:
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleDiscreteInputs: res.res       [%v]\n", res)
		}
	}
	return
}

// Modbus:HoldingRegisters:RW  AO(AnaloOutputs) - DataTypes: uint16/32/64, int16/32/64, float16/32/64
func (eh *mbTcpSvrEventHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	var regAddr uint16
	// since we're manipulating variables shared between multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	eh.lock.Lock()
	// release the lock upon return
	defer eh.lock.Unlock()

	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleHoldingRegisters: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(uint16(len(eh.holdingRegisters))) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleHoldingRegisters: req.Addr[%d] >=uint16(uint16(len(eh.holdingRegisters))=[%d] err[%s]\n", req.Addr, uint16(len(eh.holdingRegisters)), err)
		return
	}

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		if regAddr >= uint16(uint16(len(eh.holdingRegisters))) {
			err = modbus.ErrIllegalDataAddress
			logPkg.CtsLog.Error("HandleHoldingRegisters: regAddr[%d] > uint16(65535)[%d] err[%s]\n", regAddr, uint16(65535), err)
			return
		} else if len(req.Args) > i && req.Args[i] >= uint16(len(eh.holdingRegisters)) {
			err = modbus.ErrIllegalDataAddress
			logPkg.CtsLog.Error("HandleHoldingRegisters: req.Args[%d]=[%d] > uint16(len(eh.holdingRegisters)[%d] err[%s]\n", i, req.Args[i], uint16(len(eh.holdingRegisters)), err)
			return
		}
		switch regAddr {
		// expose the static, read-only value of 0xff00 in register 100
		case 0x100:
			res = append(res, 0xff00)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose holdingReg1 in register 101 (RW)
		case 0x101:
			if req.IsWrite {
				if len(req.Args) > i {
					eh.holdingReg1 = req.Args[i]
				}
			}
			res = append(res, eh.holdingReg1)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose holdingReg2 in register 102 (RW)
		case 0x102:
			if req.IsWrite {
				// only accept values 2 and 4
				if len(req.Args) > i {
					switch req.Args[i] {
					case 2, 4:
						eh.holdingReg2 = req.Args[i]
						// make note of the change (e.g. for auditing purposes)
						fmt.Printf("%s set reg#102 to %v\n", req.ClientAddr, eh.holdingReg2)
					default:
						// if the written value is neither 2 nor 4,
						// return a modbus "illegal data value" to
						// let the client know that the value is
						// not acceptable.
						err = modbus.ErrIllegalDataValue
						return
					}
				}
			}
			res = append(res, eh.holdingReg2)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose eh.holdingReg3 in register 103 (RW)
		// note: eh.holdingReg3 is a signed 16-bit integer
		case 0x103:
			if req.IsWrite {
				// cast the 16-bit unsigned integer passed by the server
				// to a 16-bit signed integer when writing
				if len(req.Args) > i {
					eh.holdingReg3 = int16(req.Args[i])
				}
			}
			// cast the 16-bit signed integer from the handler to a 16-bit unsigned
			// integer so that we can append it to `res`.
			res = append(res, uint16(eh.holdingReg3))
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose the 16 most-significant bits of eh.holdingReg4 in register 200
		case 0x200:
			if req.IsWrite {
				if len(req.Args) > i {
					eh.holdingReg4 =
						((uint32(req.Args[i])<<16)&0xffff0000 |
							(eh.holdingReg4 & 0x0000ffff))
				}

			}
			res = append(res, uint16((eh.holdingReg4>>16)&0x0000ffff))
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose the 16 least-significant bits of eh.holdingReg4 in register 201
		case 0x201:
			if req.IsWrite {
				if len(req.Args) > i {
					eh.holdingReg4 =
						(uint32(req.Args[i])&0x0000ffff |
							(eh.holdingReg4 & 0xffff0000))
				}
			}
			res = append(res, uint16(eh.holdingReg4&0x0000ffff))
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		case 0x202:
			if req.IsWrite {
				if len(req.Args) > i {
					eh.holdingReg4 =
						(uint32(req.Args[i])&0x0000ffff |
							(eh.holdingReg4 & 0xffff0000))
				}
			}
			res = append(res, uint16(eh.holdingReg4&0x0000ffff))
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("HandleHoldingRegisters: res.res       [%v]\n", res)

		// any other address is unknown
		default:
			// Append and update HoldingRegister
			if len(req.Args) > i {
				res = append(res, eh.holdingRegisters[req.Args[i]])
				eh.holdingRegisters[req.Args[i]]++
			} else {
				res = append(res, eh.holdingRegisters[regAddr])
				eh.holdingRegisters[regAddr]++
			}
		}
	}
	return
}

// Modbus:InputRegisters:RO  AI(AnalogInputs) - DataTypes: uint16/32/64, int16/32/64, float16/32/64
func (eh *mbTcpSvrEventHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	var regAddr uint16
	// since we're manipulating variables shared between multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	eh.lock.Lock()
	// release the lock upon return
	defer eh.lock.Unlock()

	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("InputRegisters: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(uint16(len(eh.inputRegisters))) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("InputRegisters: req.Addr[%d] >= uint16(65535)[%d] err[%s]\n", req.Addr, uint16(65535), err)
		return
	}

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		if regAddr >= uint16(uint16(len(eh.inputRegisters))) {
			err = modbus.ErrIllegalDataAddress
			logPkg.CtsLog.Error("InputRegisters: regAddr[%d] >= uint16(65535)[%d] err[%s]\n", regAddr, uint16(65535), err)
			return
		}
		switch regAddr {
		// expose the static, read-only value of 0xff00 in register 100
		case 0x100:
			res = append(res, 0xff00)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// expose holdingReg1 in register 101 (RW)
		case 0x101:
			res = append(res, eh.holdingReg1)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// expose holdingReg2 in register 102 (RW)
		case 0x102:
			res = append(res, eh.holdingReg2)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// expose eh.holdingReg3 in register 103 (RW)
		// note: eh.holdingReg3 is a signed 16-bit integer
		case 0x103:
			// cast the 16-bit signed integer from the handler to a 16-bit unsigned
			// integer so that we can append it to `res`.
			res = append(res, uint16(eh.holdingReg3))
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// expose the 16 most-significant bits of eh.holdingReg4 in register 200
		case 0x200:
			res = append(res, uint16((eh.holdingReg4>>16)&0x0000ffff))
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// expose the 16 least-significant bits of eh.holdingReg4 in register 201
		case 0x201:
			res = append(res, uint16(eh.holdingReg4&0x0000ffff))
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		case 0x202:
			res = append(res, uint16(eh.holdingReg4&0x0000ffff))
			logPkg.CtsLog.Warn("InputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("InputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("InputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("InputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("InputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("InputRegisters: res.res       [%v]\n", res)

		// any other address is unknown
		default:
			// FAKE: Input Regiter Changing Append and update HoldingRegister
			res = append(res, eh.inputRegisters[regAddr])
			eh.inputRegisters[regAddr]++
		}
	}
	return
}

// ============================
// Modbus Secure Server TLS+TCP
// ============================
// Handler object, passed to the NewServer() constructor above.
type mbTlsServerEventHandler struct {
	// this lock is used to avoid concurrency issues between goroutines, as
	// handler methods are called from different goroutines
	// (1 goroutine per client)
	lock sync.RWMutex

	// unix timestamp register, incremented in the main() function above and exposed
	// as a 32-bit holding register (2 consecutive 16-bit modbus registers).
	clock uint32

	// Emulating Secure TLS Modbus Registers
	coils            [65535]bool   // DO:(DigitalOutputs)
	discreteInputs   [65535]bool   // DI:(DigitalInputs)
	holdingRegisters [65535]uint16 //AO:(DigitalOutputs)
	InpputRegisters  [65535]uint16 //AI:(DigitalInputs)
}

// Modbus:Secure(TLS):Coils:RW  DO(DigitalOutputs) true/false
func (mbTlsSvrEh *mbTlsServerEventHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	// since we're manipulating variables accessed from multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	mbTlsSvrEh.lock.Lock()
	// release the lock upon return
	defer mbTlsSvrEh.lock.Unlock()

	if req.UnitId > 31 {
		// note: we're merely filtering here, but we could as well use the unit
		// ID field to support multiple register maps in a single server.
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleCoils: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(len(mbTlsSvrEh.coils)) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleCoils: req.Addr[%d] >= uint16(len(mbTlsSvrEh.coils)[%d]\n", req.Addr, uint16(len(mbTlsSvrEh.coils)))
		return nil, err
	}

	// make sure that all registers covered by this request actually exist
	if int(req.Addr)+int(req.Quantity) > len(mbTlsSvrEh.coils) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleCoils: int(req.Addr)+int(req.Quantity)[%d] > len(eh.coils)[%d] err[%s]\n", int(req.Addr)+int(req.Quantity), len(mbTlsSvrEh.coils), err)
		return
	}

	// loop through `req.Quantity` registers, from address `req.Addr` to
	// `req.Addr + req.Quantity - 1`, which here is conveniently `req.Addr + i`
	for i := 0; i < int(req.Quantity); i++ {
		// ignore the write if the current register address is 80
		if req.IsWrite && int(req.Addr)+i != 80 {
			// assign the value
			if len(req.Args) > i {
				mbTlsSvrEh.coils[int(req.Addr)+i] = req.Args[i]
			}
		}
		// append the value of the requested register to res so they can be and interlock coil
		// sent back to the client
		res = append(res, mbTlsSvrEh.coils[int(req.Addr)+i])
		if mbTlsSvrEh.coils[int(req.Addr)+i] {
			mbTlsSvrEh.coils[int(req.Addr)+i] = false
		} else {
			mbTlsSvrEh.coils[int(req.Addr)+i] = true
		}
		switch req.Addr {
		case 0x100:
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, mbTlsSvrEh.coils[int(req.Addr)+i])
			logPkg.CtsLog.Warn("mbTls:HandleCoils: res           [%v]\n", res)
		case 0x101:
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, mbTlsSvrEh.coils[int(req.Addr)+i])
			logPkg.CtsLog.Warn("mbTls:HandleCoils: res           [%v]\n", res)
		case 0x102:
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, mbTlsSvrEh.coils[int(req.Addr)+i])
			logPkg.CtsLog.Warn("mbTls:HandleCoils: res           [%v]\n", res)
		case 0x103:
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, mbTlsSvrEh.coils[int(req.Addr)+i])
			logPkg.CtsLog.Warn("mbTls:HandleCoils: res           [%v]\n", res)
		case 0x104:
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleCoils: eh.coils[int(req.Addr[0x%X])+i(%d)]=%v\n", int(req.Addr), i, mbTlsSvrEh.coils[int(req.Addr)+i])
			logPkg.CtsLog.Warn("mbTls:HandleCoils: res           [%v]\n", res)
		}
	}
	return
}

// Modbus:Secure(TLS)-DiscreteInputs:RO  DI(DigitalInputs) true/false
func (mbTlsSvrEh *mbTlsServerEventHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	// since we're manipulating variables accessed from multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	mbTlsSvrEh.lock.Lock()
	// release the lock upon return
	defer mbTlsSvrEh.lock.Unlock()

	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleDiscreteInputs: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}
	if req.Addr >= uint16(len(mbTlsSvrEh.discreteInputs)) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleDiscreteInputs: req.Addr[%d] >= uint16(len(mbTlsSvrEh.discreteInputs)[%d]\n", req.Addr, uint16(len(mbTlsSvrEh.discreteInputs)))
		return nil, err
	}

	// make sure that all registers covered by this request actually exist
	if int(req.Addr)+int(req.Quantity) > len(mbTlsSvrEh.discreteInputs) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("HandleDiscreteInputs: int(req.Addr)+int(req.Quantity)[%d] > len(eh.discreteInputs)[%d] err[%s]\n", int(req.Addr)+int(req.Quantity), len(mbTlsSvrEh.discreteInputs), err)
		return
	}

	// loop through `req.Quantity` registers, from address `req.Addr` to
	// `req.Addr + req.Quantity - 1`, which here is conveniently `req.Addr + i`
	for i := 0; i < int(req.Quantity); i++ {
		// sent back to the client and interlock it discrete input value
		res = append(res, mbTlsSvrEh.discreteInputs[int(req.Addr)+i])
		if mbTlsSvrEh.discreteInputs[int(req.Addr)+i] {
			mbTlsSvrEh.discreteInputs[int(req.Addr)+i] = false
		} else {
			mbTlsSvrEh.discreteInputs[int(req.Addr)+i] = true
		}
		switch req.Addr {
		case 0x100:
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x101:
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x102:
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x103:
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: res.res       [%v]\n", res)
		case 0x104:
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleDiscreteInputs: res.res       [%v]\n", res)
		}
	}
	return
}

// Modbus:Secure(TLS):HoldingRegisters:RW  AO(AnaloOutputs) - DataTypes: uint16/32/64, int16/32/64, float16/32/64
func (mbTlsSvrEh *mbTlsServerEventHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	var regAddr uint16
	// since we're manipulating variables accessed from multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	mbTlsSvrEh.lock.Lock()
	// release the lock upon return
	defer mbTlsSvrEh.lock.Unlock()

	// require the "operator" role for write operations (i.e. set the clock).
	if req.IsWrite && req.ClientRole != "operator" {
		logPkg.CtsLog.Error("write access denied: client %s missing the 'operator' role (role: '%s')\n",
			req.ClientAddr, req.ClientRole)
		err = modbus.ErrIllegalFunction
		return
	}
	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleHoldingRegisters: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}

	if req.Addr >= uint16(len(mbTlsSvrEh.holdingRegisters)) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleHoldingRegisters: req.Addr[%d] >= uint16(len(mbTlsSvrEh.holdingRegisters)[%d]\n", req.Addr, uint16(len(mbTlsSvrEh.holdingRegisters)))
		return nil, err
	}

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		switch regAddr {
		// expose the 16 most-significant bits of the clock in register #0
		case 0x100:
			if req.IsWrite {
				mbTlsSvrEh.clock =
					((uint32(req.Args[i])<<16)&0xffff0000 |
						(mbTlsSvrEh.clock & 0x0000ffff))
			}
			mbTlsSvrEh.holdingRegisters[uint16(req.Addr)] = uint16(mbTlsSvrEh.clock & 0x0000ffff)
			res = append(res, uint16((mbTlsSvrEh.clock>>16)&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: res.res       [%v]\n", res)

		// expose the 16 least-significant bits of the clock in register #1
		case 0x101:
			if req.IsWrite {
				mbTlsSvrEh.clock =
					(uint32(req.Args[i])&0x0000ffff |
						(mbTlsSvrEh.clock & 0xffff0000))
			}
			mbTlsSvrEh.holdingRegisters[uint16(req.Addr)] = uint16(mbTlsSvrEh.clock & 0x0000ffff)
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: res.res       [%v]\n", res)

		case 0x102:
			if req.IsWrite {
				mbTlsSvrEh.clock =
					(uint32(req.Args[i])&0x0000ffff |
						(mbTlsSvrEh.clock & 0xffff0000))
			}
			mbTlsSvrEh.holdingRegisters[uint16(req.Addr)] = uint16(mbTlsSvrEh.clock & 0x0000ffff)
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: res.res       [%v]\n", res)

		case 0x103:
			if req.IsWrite {
				mbTlsSvrEh.clock =
					(uint32(req.Args[i])&0x0000ffff |
						(mbTlsSvrEh.clock & 0xffff0000))
			}
			mbTlsSvrEh.holdingRegisters[uint16(req.Addr)] = uint16(mbTlsSvrEh.clock & 0x0000ffff)
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: res.res       [%v]\n", res)

		case 0x104:
			if req.IsWrite {
				mbTlsSvrEh.clock =
					(uint32(req.Args[i])&0x0000ffff |
						(mbTlsSvrEh.clock & 0xffff0000))
			} else {
				mbTlsSvrEh.coils[uint16(req.Addr)] = true
			}
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleHoldingRegisters: res.res       [%v]\n", res)

		// any other address is unknown
		default:
			mbTlsSvrEh.coils[uint16(req.Addr)] = true
			mbTlsSvrEh.holdingRegisters[uint16(req.Addr)] = uint16(mbTlsSvrEh.clock & 0x0000ffff)
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
		}
	}
	return
}

// Modbus:Secure(TLS)-InputRegisters:RO  AI(AnalogInputs) - DataTypes: uint16/32/64, int16/32/64, float16/32/64
func (mbTlsSvrEh *mbTlsServerEventHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	var regAddr uint16

	// since we're manipulating variables accessed from multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	mbTlsSvrEh.lock.Lock()
	// release the lock upon return
	defer mbTlsSvrEh.lock.Unlock()

	if req.UnitId > 31 {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleInputRegisters: req.UnitId[%d] > 31 err[%s]\n", req.UnitId, err)
		return
	}

	if req.Addr >= uint16(len(mbTlsSvrEh.InpputRegisters)) {
		err = modbus.ErrIllegalDataAddress
		logPkg.CtsLog.Error("mbTls:HandleInputRegisters: req.Addr[%d] >= uint16(len(mbTlsSvrEh.InpputRegisters)[%d]\n", req.Addr, uint16(len(mbTlsSvrEh.InpputRegisters)))
		return nil, err
	}

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		mbTlsSvrEh.InpputRegisters[uint16(req.Addr)] = uint16((mbTlsSvrEh.clock >> 16) & 0x0000ffff)
		switch regAddr {
		// expose the 16 most-significant bits of the clock in register #0
		case 0x100:
			res = append(res, uint16((mbTlsSvrEh.clock>>16)&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: res.res       [%v]\n", res)

		// expose the 16 least-significant bits of the clock in register #1
		case 0x101:
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: res.res       [%v]\n", res)
			// expose the 16 least-significant bits of the clock in register #1
		case 0x102:
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: res.res       [%v]\n", res)
			// expose the 16 least-significant bits of the clock in register #1
		case 0x103:
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: res.res       [%v]\n", res)
			// expose the 16 least-significant bits of the clock in register #1
		case 0x104:
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientAddr[%s]\n", req.ClientAddr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.UnitId    [%d]\n", req.UnitId)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.ClientRole[%s]\n", req.ClientRole)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Addr      [0x%X(%d)]\n", req.Addr, req.Addr)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: req.Quantity  [%d]\n", req.Quantity)
			logPkg.CtsLog.Warn("mbTls:HandleInputRegisters: res.res       [%v]\n", res)

		// any other address is unknown
		default:
			res = append(res, uint16(mbTlsSvrEh.clock&0x0000ffff))
		}
	}
	return
}
