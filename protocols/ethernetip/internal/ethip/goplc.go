package ethip

import (
	"bytes"
	"ethernetip/protocols/ethernetip"
	"ethernetip/protocols/ethernetip/commonIndustrialProtocol"
	"ethernetip/protocols/ethernetip/commonIndustrialProtocol/segment"
	"ethernetip/protocols/ethernetip/commonIndustrialProtocol/segment/epath"
	"ethernetip/protocols/ethernetip/lib"
	_type "ethernetip/protocols/ethernetip/type"
	golog "ethernetip/utils/gologtofile"
	logPkg "ethernetip/utils/gologtofile"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"

	//"strings"
	"time"
	//mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 24/04/2022 removed 	config "ethernetip/utils/goconfigyml"

const ServiceGetAttributeSingle = 0x0E
const ServiceSetAttributeSingle = 0x10

type controller struct {
	VendorID     _type.UINT
	DeviceType   _type.UINT
	ProductCode  _type.UINT
	Major        _type.USINT
	Minor        _type.USINT
	Status       _type.UINT
	SerialNumber _type.UDINT
	Version      string
	Name         string
}

type EipPlc struct {
	tcpAddr    *net.TCPAddr
	tcpConn    *net.TCPConn
	config     *Config
	session    _type.UDINT
	Controller *controller

	writeRoute  bool
	sender      chan []byte
	bufferData  []byte
	TargetPath  []byte
	HandleMap   map[ethernetip.Command]func(*ethernetip.Encapsulation)
	OnConnected func()

	ContextPool map[_type.ULINT]func(*commonIndustrialProtocol.MessageRouterResponse)
}
type Plc_io struct {
	bit_memory       byte
	word_memory      uint32
	word16_memory    uint16
	word32_memory    uint32
	type_data_memory string
	newData          bool
	responseStatus   uint8
}

var Plc_data Plc_io

func DefaultConfig() *Config {
	_defaultConfig := &Config{}
	_defaultConfig.Port = 0xAF12
	_defaultConfig.ReconnectionInterval = 3
	_defaultConfig.Logger = nil

	return _defaultConfig
}
func (p *EipPlc) Connect() error {
	golog.CtsLog.Info("PLC EthernetIP TCP Connecting to Ip:Port", p.tcpAddr)
	_conn, err := net.DialTCP("tcp", nil, p.tcpAddr)
	if err != nil {
		return err
	}

	err2 := _conn.SetKeepAlive(true)
	if err2 != nil {
		return err2
	}

	p.tcpConn = _conn
	p.connected()
	return nil
}

func (p *EipPlc) connected() {
	golog.CtsLog.Info("PLC EthernetIP Connected with Ip:Port", p.tcpAddr)
	p.bufferData = []byte{}

	if !p.writeRoute {
		go p.write()
	}

	go p.read()

	p.config.Println("PLC EthernetIP Register Session with Ip:Port", p.tcpAddr)
	encapsulation := ethernetip.RequestRegisterSession(0)

	p.sender <- encapsulation.Buffer()
}

func (p *EipPlc) disconnected(err error) {
	if err == io.EOF {
		p.config.Println("PLC Disconnected from EthernetIP Ip:Port", p.tcpAddr)
		p.config.Println("EOF")
	} else {
		p.config.Println("PLC Disconnected from EthernetIP Ip:Port", p.tcpAddr)
		p.config.Println(err)
	}
	if p.tcpConn != nil {
		_ = p.tcpConn.Close()
		p.tcpConn = nil
	}

	p.tcpConn = nil

	if p.config.ReconnectionInterval != 0 {
		time.Sleep(p.config.ReconnectionInterval)
		golog.CtsLog.Info("PLC Reconnecting to EthernetIP Ip:Port", p.tcpAddr)
		err := p.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (p *EipPlc) read() {
	defer func() {
		if err := recover(); err != nil {
			go p.disconnected(err.(error))
		}
	}()

	buf := make([]byte, 1024*64)
	var err error
	for {
		var length int
		// Set a timeout for the read and write operations
		//timeout := 5 * time.Second
		//p.tcpConn.SetDeadline(time.Now().Add(timeout))
		length, err = p.tcpConn.Read(buf)
		if err != nil {
			break
		}
		//p.config.Printf("PLC EthernetIP <= Received %d bytes\n", length)

		p.bufferData = append(p.bufferData, buf[0:length]...)
		if len(p.bufferData) > 24 {
			read, encapsulations := ethernetip.Slice(p.bufferData)
			p.bufferData = p.bufferData[read:]

			for _, encapsulation := range encapsulations {
				if encapsulation.Status == ethernetip.StatusSuccess {
					if exec, ok := p.HandleMap[encapsulation.Command]; ok {
						exec(encapsulation)
					} else {
						golog.CtsLog.Error("PLC EthernetIP Received encapsulation Command: %#x ,but no registered handler!\n", encapsulation.Command)
					}
				}
			}
		}
	}
}

func (p *EipPlc) handleRegisterSession(encapsulation *ethernetip.Encapsulation) {
	p.session = encapsulation.SessionHandle
	golog.CtsLog.Info("PLC EthernetIP -- Cmd %#x\n", encapsulation.Command)
	golog.CtsLog.Info("PLC EthernetIP -- Session %#x\n", p.session)
	/*
		// get_attribute_all
		mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
		mr1.Service = 0x01
		mr1.RequestPath = segment.Paths(
			epath.LogicalBuild(epath.LogicalTypeClassID, 1, true),
			epath.LogicalBuild(epath.LogicalTypeInstanceID, 1, true),
		)

		p.ContextPool[math.MaxUint64] = p.getAttributeAll
		p.UcmmSend(3, 250, math.MaxUint64, mr1)
	*/
}

/*
	func (p *EipPlc) getAttributeAll(mr *commonIndustrialProtocol.MessageRouterResponse) {
		//p.config.Printf("%+v\n", mr)
		// CIP_Networks_Common-Industrial_Protocol_Family.pdf
		// 2.5.1. Identity Object (Class ID: 0x01)              Mandatory Attributes
		dataReader := bytes.NewReader(mr.ResponseData)       //      Idx  Len Header = 15 bytes
		lib.ReadByte(dataReader, &p.Controller.VendorID)     // HF1  0    2
		lib.ReadByte(dataReader, &p.Controller.DeviceType)   // HF2  2    2
		lib.ReadByte(dataReader, &p.Controller.ProductCode)  // HF3  4    2
		lib.ReadByte(dataReader, &p.Controller.Major)        // HF4  6    1 RevisionH
		lib.ReadByte(dataReader, &p.Controller.Minor)        // HF5  7    1 RevisionL
		lib.ReadByte(dataReader, &p.Controller.Status)       // HF6  8    2 Device Status
		lib.ReadByte(dataReader, &p.Controller.SerialNumber) // HF7 10    4 32 bits
		nameLen := uint8(0)
		lib.ReadByte(dataReader, &nameLen) //                   HF8 14    1 n=PayLoadLen ASCII up to 32 char
		nameBuf := make([]byte, nameLen)
		lib.ReadByte(dataReader, nameBuf) //                    PF1 15    n =PayLoad
		p.Controller.Name = string(nameBuf)
		p.Controller.Version = fmt.Sprintf("%d.%d", p.Controller.Major, p.Controller.Minor)
		// Optional Attributes
		// - State                                                        2
		// - Configuration Consistency Value
		// - Heartbeat Interval
		// - Active Language
		// - Supported Language List
		// - International Product Name
		if p.OnConnected != nil {
			p.OnConnected()
		}
	}
*/
func (p *EipPlc) getAttributeSinglefromMemoryType(mr *commonIndustrialProtocol.MessageRouterResponse) {
	// CIP_Networks_Common-Industrial_Protocol_Family.pdf
	// 2.5.1. Identity Object (Class ID: 0x04)
	//						  (Instance ID: 102)
	//						  (Object ID: 0x03)
	Plc_data.responseStatus = uint8(mr.GeneralStatus)
	if mr.GeneralStatus == 0 { // if not equal zero CIP message return erro
		dataReader := bytes.NewReader(mr.ResponseData)
		switch len(mr.ResponseData) {
		case 1:
			Plc_data.type_data_memory = "BYTE"
			lib.ReadByte(dataReader, &Plc_data.bit_memory)
		case 2:
			Plc_data.type_data_memory = "INT"
			lib.ReadByte(dataReader, &Plc_data.word16_memory)
		case 4:
			Plc_data.type_data_memory = "WORD"
			lib.ReadByte(dataReader, &Plc_data.word32_memory)
		default:
			Plc_data.type_data_memory = "INVALID"
		}
	}
	Plc_data.newData = true
}

func (p *EipPlc) getAttributeSinglefromWordMemory(mr *commonIndustrialProtocol.MessageRouterResponse) {

	// CIP_Networks_Common-Industrial_Protocol_Family.pdf
	// 2.5.1. Identity Object (Class ID: 0x04)
	//						  (Instance ID: 102)
	//						  (Object ID: 0x03)
	Plc_data.responseStatus = uint8(mr.GeneralStatus)
	if mr.GeneralStatus == 0 { // if not equal zero CIP message return erro
		dataReader := bytes.NewReader(mr.ResponseData)
		if len(mr.ResponseData) > 2 {
			lib.ReadByte(dataReader, &Plc_data.word32_memory)
			Plc_data.word_memory = Plc_data.word32_memory
		} else {
			lib.ReadByte(dataReader, &Plc_data.word16_memory)
			Plc_data.word_memory = uint32(Plc_data.word16_memory)
		}
	}
	Plc_data.newData = true
}

func (p *EipPlc) getAttributeSinglefromBitMemory(mr *commonIndustrialProtocol.MessageRouterResponse) {

	// CIP_Networks_Common-Industrial_Protocol_Family.pdf
	// 2.5.1. Identity Object (Class ID: 0x04)
	//						  (Instance ID: 101)
	//						  (Object ID: 0x03)
	Plc_data.responseStatus = uint8(mr.GeneralStatus)
	if mr.GeneralStatus == 0 { // if not equal zero CIP message return erro
		dataReader := bytes.NewReader(mr.ResponseData)
		lib.ReadByte(dataReader, &Plc_data.bit_memory)
	}
	Plc_data.newData = true
}

func (p *EipPlc) setAttributeSingleToBitMemory(mr *commonIndustrialProtocol.MessageRouterResponse) {

	// CIP_Networks_Common-Industrial_Protocol_Family.pdf
	// 2.5.1. Identity Object (Class ID: 0x04)
	//						  (Instance ID: 101)
	//						  (Object ID: 0x03)
	//dataWriter := bytes.New(mr.ResponseData)
	Plc_data.responseStatus = uint8(mr.GeneralStatus)
	//if mr.GeneralStatus == 0 { // if not equal zero CIP message return erro
	//	data := new(bytes.Buffer)
	//	lib.WriteByte(data, &Plc_data.bit_memory)
	//}
	Plc_data.newData = true
}

func (p *EipPlc) setAttributeSingleToWordMemory(mr *commonIndustrialProtocol.MessageRouterResponse) {

	// CIP_Networks_Common-Industrial_Protocol_Family.pdf
	// 2.5.1. Identity Object (Class ID: 0x04)
	//						  (Instance ID: 101)
	//						  (Object ID: 0x03)
	//dataWriter := bytes.New(mr.ResponseData)
	Plc_data.responseStatus = uint8(mr.GeneralStatus)
	//if mr.GeneralStatus == 0 { // if not equal zero CIP message return erro
	//	data := new(bytes.Buffer)
	//	lib.WriteByte(data, &Plc_data.word_memory)
	//}
	Plc_data.newData = true
}

// func (p *EipPlc) WriteAnalogOutputs(ClassID, InstanceID, Attribute uint32, dt uint16) {
func (p *EipPlc) WriteBitMemory(ClassID, InstanceID, Attribute uint32, dt byte) (resp byte) {

	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	data := new(bytes.Buffer)

	lib.WriteByte(data, dt)

	mr1.RequestData = data.Bytes()
	//fmt.Printf("PLC EthernetIP mr1.RequestData = %v\n", mr1.RequestData)

	Plc_data.newData = false
	mr1.Service = ServiceSetAttributeSingle
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, ClassID, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, InstanceID, true),
		epath.LogicalBuild(epath.LogicalTypeAttributeID, Attribute, true),
	)

	p.ContextPool[math.MaxUint32] = p.setAttributeSingleToBitMemory
	p.UcmmSend(3, 250, math.MaxUint32, mr1)

	for !Plc_data.newData {
		time.Sleep(time.Millisecond)
	}
	return Plc_data.responseStatus
}

func (p *EipPlc) WriteWordMemory(ClassID, InstanceID, Attribute, dt uint32) (resp byte) {

	typeMemory := p.CheckTypeMemory(ClassID, InstanceID, Attribute)

	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	data := new(bytes.Buffer)

	switch typeMemory {
	case "BYTE":
		if dt > 0xFF {
			Plc_data.responseStatus = 0x15 // Data to write must be <= 0xFF
			return Plc_data.responseStatus
		}
		lib.WriteByte(data, uint8(dt))
	case "INT":
		if dt > 0xFFFF {
			Plc_data.responseStatus = 0x15 // Data to write must be <= 0xFFFF
			return Plc_data.responseStatus
		}
		lib.WriteByte(data, uint16(dt))
	case "WORD":
		if dt > 0xFFFFFFFF {
			Plc_data.responseStatus = 0x15 // Data to write must be <= 0xFFFFFFFF
			return Plc_data.responseStatus
		}
		lib.WriteByte(data, dt)
	default:
		Plc_data.responseStatus = 0x18 // Data type not supported
	}

	mr1.RequestData = data.Bytes()
	//fmt.Printf("PLC EthernetIP mr1.RequestData = %v\n", mr1.RequestData)

	Plc_data.newData = false
	mr1.Service = ServiceSetAttributeSingle
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, ClassID, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, InstanceID, true),
		epath.LogicalBuild(epath.LogicalTypeAttributeID, Attribute, true),
	)

	p.ContextPool[math.MaxUint32] = p.setAttributeSingleToWordMemory
	p.UcmmSend(3, 250, math.MaxUint32, mr1)

	for !Plc_data.newData {
		time.Sleep(time.Millisecond)
	}
	return Plc_data.responseStatus
}

var mutex sync.Mutex

// func (p *EipPlc) ReadAnalogInputs(ClassID, InstanceID, Attribute uint32) (rd_io uint16) {
func (p *EipPlc) ReadWordMemory(ClassID, InstanceID, Attribute uint32) (rd_io uint32, resp byte) {

	mutex.Lock()

	Plc_data.newData = false
	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	mr1.Service = ServiceGetAttributeSingle
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, ClassID, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, InstanceID, true),
		epath.LogicalBuild(epath.LogicalTypeAttributeID, Attribute, true),
	)

	context := _type.ULINT(rand.Uint64())
	p.ContextPool[context] = p.getAttributeSinglefromWordMemory
	p.UcmmSend(3, 250, context, mr1)

	defer mutex.Unlock()
	waitData := 5000
	for !Plc_data.newData {
		time.Sleep(time.Millisecond)
		waitData--
		if waitData == 0 {
			return Plc_data.word_memory, 0x01
		}
	}
	return Plc_data.word_memory, Plc_data.responseStatus

}

func (p *EipPlc) ReadBitMemory(ClassID, InstanceID, AttributeID uint32) (p_io byte, resp byte) {

	mutex.Lock()

	Plc_data.newData = false
	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	mr1.Service = ServiceGetAttributeSingle
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, ClassID, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, InstanceID, true),
		epath.LogicalBuild(epath.LogicalTypeAttributeID, AttributeID, true),
	)
	p.ContextPool[math.MaxUint32] = p.getAttributeSinglefromBitMemory
	p.UcmmSend(3, 250, math.MaxUint32, mr1)

	defer mutex.Unlock()
	waitData := 5000
	for !Plc_data.newData {
		time.Sleep(time.Millisecond)
		waitData--
		if waitData == 0 {
			return Plc_data.bit_memory, 0x01
		}
	}
	return Plc_data.bit_memory, Plc_data.responseStatus
}

func (p *EipPlc) CheckTypeMemory(ClassID, InstanceID, AttributeID uint32) string {
	Plc_data.newData = false
	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	mr1.Service = ServiceGetAttributeSingle
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, ClassID, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, InstanceID, true),
		epath.LogicalBuild(epath.LogicalTypeAttributeID, AttributeID, true),
	)
	p.ContextPool[math.MaxUint32] = p.getAttributeSinglefromMemoryType
	p.UcmmSend(3, 250, math.MaxUint32, mr1)
	for !Plc_data.newData {
		time.Sleep(time.Millisecond)
	}

	return Plc_data.type_data_memory
}

func (p *EipPlc) UcmmSend(timeTicks _type.USINT, timeoutTicks _type.USINT, context _type.ULINT, mr1 *commonIndustrialProtocol.MessageRouterRequest) {
	ucmm := &commonIndustrialProtocol.UnconnectedSend{}
	ucmm.TimeTick = timeTicks
	ucmm.TimeOutTicks = timeoutTicks
	ucmm.MessageRequest = mr1
	ucmm.RouterPath = p.TargetPath

	cpf := &ethernetip.CommonPacketFormat{}

	cpf.UnconnectedData(mr1.Buffer())
	pkg := ethernetip.RequestSendRRData(p.session, context, 10, cpf)
	p.sender <- pkg.Buffer()
}

func (p *EipPlc) handleSendData(encapsulation *ethernetip.Encapsulation) {
	cpf := ethernetip.SendRRDataParser(encapsulation.Data)
	mr := commonIndustrialProtocol.MRParser(cpf.DataItem.Data)
	if mr.GeneralStatus != 0 {
		golog.CtsLog.Error("PLC EthernetIP target error => Service Code: %#x | Status: %#x | Addtional: %s\n", mr.ReplyService, mr.GeneralStatus, mr.AdditionalStatus)
	}
	p.ContextPool[encapsulation.SenderContext](mr)
	delete(p.ContextPool, encapsulation.SenderContext)

}

func (p *EipPlc) write() {
	p.writeRoute = true
	for {
		data := <-p.sender
		_, _ = p.tcpConn.Write(data)
		//p.config.Printf("PLC EthernetIP => Sent %d bytes\n", len(data))
	}
}

func NewOriginator(addr string, slot uint8, cfg *Config) (*EipPlc, error) {
	_plc := &EipPlc{}
	_plc.config = cfg
	if _plc.config == nil {
		_plc.config = defaultConfig
	}

	_tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, _plc.config.Port))
	if err != nil {
		return nil, err
	}

	_plc.tcpAddr = _tcp
	_plc.Controller = &controller{}
	_plc.sender = make(chan []byte)
	_plc.HandleMap = make(map[ethernetip.Command]func(*ethernetip.Encapsulation))
	_plc.TargetPath = epath.PortBuild([]byte{slot}, 1, true)
	_plc.ContextPool = make(map[_type.ULINT]func(*commonIndustrialProtocol.MessageRouterResponse))

	_plc.HandleMap[ethernetip.CommandNOP] = ethernetip.HandleNop
	_plc.HandleMap[ethernetip.CommandListIdentity] = ethernetip.HandleListIdentity
	_plc.HandleMap[ethernetip.CommandListInterfaces] = ethernetip.HandleListInterfaces
	_plc.HandleMap[ethernetip.CommandRegisterSession] = _plc.handleRegisterSession
	_plc.HandleMap[ethernetip.CommandSendRRData] = _plc.handleSendData
	_plc.HandleMap[ethernetip.CommandSendUnitData] = _plc.handleSendData

	return _plc, nil
}

func ConnEthernetIP(EipPlcIpAddr, EipPlcPort string) (*EipPlc, error) {
	EipPlc := &EipPlc{}
	golog.CtsLog.Info("STEP04 Initializing EthernetIp\n")

	eipcfg := DefaultConfig()
	//Reconnection Interval  default will not automatically reconnect
	eipcfg.ReconnectionInterval = time.Second * 3

	//EthernetIp Logger  default no log
	eipcfg.Logger = log.New(os.Stdout, "", log.LstdFlags)
	//EthernetIp Port
	port, _ := strconv.Atoi(EipPlcPort)
	eipcfg.Port = uint16(port)

	strtmp := fmt.Sprintf("EthernetIp PLC Ip[%s] Port[%s] ", EipPlcIpAddr, EipPlcPort)

	if !isTCPConnected(EipPlcIpAddr, EipPlcPort) {
		return EipPlc, fmt.Errorf("IP not connected")
	}

	Eipplc, err := NewOriginator(EipPlcIpAddr, 1, eipcfg)
	if err != nil {
		logPkg.CtsLog.Error("STEP04b: FAIL to Monitor to " + strtmp)
		logPkg.CtsLog.Error("%v", err)
		return nil, err
	} else {
		logPkg.CtsLog.Info("STEP04a: PASS Monitoring " + strtmp)
	}

	strtmp = fmt.Sprintf("EthernetIp Ip[%s] Port[%s] ", EipPlcIpAddr, EipPlcPort)
	err = Eipplc.Connect()
	if err != nil {
		logPkg.CtsLog.Error("STEP04d: FAIL to Connect to " + strtmp)
		logPkg.CtsLog.Error("%v", err)
	} else {
		logPkg.CtsLog.Info("STEP05a: PASS Connected " + strtmp)

	}
	return Eipplc, err
}
