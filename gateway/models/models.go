package models

import "os/exec"

type MQTT struct {
	Server       string `json:"server"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	InputInfo    string `json:"inputInfo"`
	ErrorInfo    string `json:"errorInfo"`
	StatusDevice string `json:"satusDevice"`
}

type GatewayProtocol struct {
	Protocol string
	Exec     string
}

type Process struct {
	Name   string    `json:"name"`
	Pid    int       `json:"pid"`
	Cmd    *exec.Cmd `json:"cmd"`
	Active bool      `json:"active"`
	Exec   string    `json:"exec"`
}

type Configurations struct {
	KeepAlive struct {
		Interval int `json:"interval"`
	}
	Backend struct {
		Host string `json:"host"`
	}
	Protocols struct {
		Directory string `json:"directory"`
	}
	Log struct {
		Filename   string `json:"Filename"`
		MaxSize    int    `json:"MaxSize"`
		MaxBackups int    `json:"MaxBackups"`
		MaxAge     int    `json:"MaxAge"`
		Compress   bool   `json:"Compress"`
	}
}
