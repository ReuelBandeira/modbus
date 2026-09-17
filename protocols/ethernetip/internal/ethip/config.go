package ethip

import (
	"log"
	"time"
)

var defaultConfig *Config

type Config struct {
	Port                 uint16
	ReconnectionInterval time.Duration
	Logger               *log.Logger
}

func (c *Config) Println(v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Println(v...)
	}
}

func (c *Config) Printf(format string, v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Printf(format, v...)
	}
}
