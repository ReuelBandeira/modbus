package modelsswagger

type Protocol struct {
	Protocol string
}

type Type struct {
	Type string
}

type Server struct {
	Name       string
	IP         string
	Port       string
	ProtocolID string
}

type PLC struct {
	Name string
}

type Device struct {
	Name        string
	ProtocolID  string
	PLCID       string
	TypeID      string
	MemoryBlock string
	Class       string
	Instance    string
	Attribute   string
	Address     string
}
