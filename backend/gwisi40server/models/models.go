package models

type User struct {
	Id       uint
	Name     string
	Email    string
	Password []byte
	Language string
}

type Login struct {
	Email    string
	Password string
}

type Language struct {
	ID       uint   `json:"id"`
	Language string `json:"language"`
}

type Token struct {
	Token string
}

type Protocol struct {
	ID       uint   `json:"id"`
	Protocol string `json:"protocol" gorm:"unique"`
	Alias    string `json:"alias"`
	MD5      string `json:"md5"`
	Version  uint   `json:"version"`
	Config   string `json:"config"`
}

type Type struct {
	ID   uint
	Type string
}

type Server struct {
	ID         uint
	Name       string
	IP         string
	Port       int
	ProtocolID uint `json:"-"`
	Protocol   Protocol
}

type PLC struct {
	ID   uint
	Name string
}

type Device struct {
	ID          uint
	Name        string
	ProtocolID  uint `json:"-"`
	Protocol    Protocol
	PLCID       uint `json:"-"`
	PLC         PLC
	TypeID      uint `json:"-"`
	Type        Type
	MemoryBlock int
	Class       int
	Instance    int
	Attribute   int
	Address     int
}

type DeviceConfiguration struct {
	ID           uint   `json:"-"`
	Manufacturer string `json:"fabricante"`
	Serie        string `json:"serie"`
}

type Memories struct {
	ID       uint   `json:"id"`
	DeviceID uint   `json:"-"`
	Label    string `json:"label"`
	Type     string `json:"tipo"`
	Format   string `json:"formato"`
}

type Values struct {
	ID       uint `json:"-"`
	MemoryID uint `json:"-"`
	Key      string
	Value    string
}

type MQTTSettings struct {
	ID           uint   `json:"-"`
	Server       string `json:"server"`
	Port         string `json:"port"`
	Connected    bool   `json:"connected"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	InputInfo    string `json:"inputInfo"`
	ErrorInfo    string `json:"errorInfo"`
	StatusDevice string `json:"statusDevice"`
}

type Configurations struct {
	ID          uint   `json:"id"`
	Version     uint   `json:"version"`
	Address     string `json:"address"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	ReadingTime uint   `json:"readingTime"`
	Topics      string `json:"topics"`
	Devices     string `json:"devices"`
}

type DataVersion struct {
	ID       uint   `json:"-"`
	DataType string `json:"data-type"`
	Version  uint   `json:"version"`
}

type Resources struct {
	ID               uint    `json:"-"`
	MemoryUsage      float64 `json:"memoryUsage"`
	DiskUsage        float64 `json:"diskUsage"`
	ConnectedDevices string  `json:"connectedDevices"`
}

type Logfile struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}
