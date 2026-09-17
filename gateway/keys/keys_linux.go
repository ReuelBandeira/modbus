package keys

import (
	"encoding/json"
	"fmt"
	"gateway/globals"
	"gateway/mqtt"
	"os/exec"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const low = 0
const high = 1

var resetDetected bool = false
var rebootDetected bool = false
var resetTimeCount int = 0
var rebootTimeCount int = 0

func ScanKeys() {
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		fmt.Sprint("unable to open gpio", err.Error())
		return
	}

	defer rpio.Close()
	pin := rpio.Pin(25) // VENTILADOR
	pin.Output()
	pin.Write(high) // Liga o ventilador

	pin = rpio.Pin(7) // LED1, Blue
	pin.Output()
	pin.Write(low) // Desliga o LED Blue

	pin = rpio.Pin(16) // LED SYS
	pin.Output()
	pin.Write(low) // Desliga o LED SYS

	pin = rpio.Pin(6) // Tecla Reset
	pin.Input()

	pin = rpio.Pin(5) // Tecla Reboot
	pin.Input()

	for {
		if !resetDetected && !rebootDetected {
			pin = rpio.Pin(16)
			pin.Write(low) // LED SYS Off
			time.Sleep(1000 * time.Millisecond)
		}
		pin = rpio.Pin(6) // Tecla Reset
		if pin.Read() == low {
			if resetTimeCount == 0 {
				resetDetected = true
				resetTimeCount = 40
				pin = rpio.Pin(16) // LED SYS
				pin.Toggle()       // LED SYS --> Low --> High --> Low ...
			} else {
				time.Sleep(250 * time.Millisecond)
				pin = rpio.Pin(16) // LED SYS
				pin.Toggle()
				resetTimeCount--
				if resetTimeCount == 0 {
					// Run the reboot command using the 'shutdown' utility
					cmd := exec.Command("sudo", "shutdown", "-r", "now")
					// Run the command and capture its output and error streams
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println("Error:", err)
					}
					// Print the command output
					fmt.Println(string(output))
				}
			}
		} else {
			resetTimeCount = 0
			resetDetected = false
		}
		// Verifica se tecla Reboot está ativada
		pin = rpio.Pin(5) // Tecla Reboot
		if pin.Read() == low {
			if rebootTimeCount == 0 {
				rebootDetected = true
				rebootTimeCount = 10
				pin = rpio.Pin(16) // LED SYS
				pin.Toggle()
			} else {
				time.Sleep(500 * time.Millisecond)
				pin = rpio.Pin(16) // LED SYS
				pin.Toggle()
				rebootTimeCount--
				if rebootTimeCount == 0 {
					SendRebootProtocol()
					pin = rpio.Pin(16)
					pin.Write(low) // LED SYS Off
				}
			}
		} else {
			rebootTimeCount = 0
			rebootDetected = false

		}
	}
}

func SendRebootProtocol() {
	type MessageRebootProtocol struct {
		MessageType string `json:"messageType"`
		Data        struct {
			Action string `json:"action"`
			Name   string `json:"protocol"`
		} `json:"data"`
	}

	for _, protocol := range globals.GatewayProtocols {
		var data MessageRebootProtocol

		data.MessageType = "action"
		data.Data.Action = "restart"

		data.Data.Name = protocol.Protocol

		jsonPayload, err := json.Marshal(data)
		if err != nil {
			panic(err) //Check - não pode parar o programa
		}

		mqtt.SendMsg(globals.LocalClient, protocol.Protocol, string(jsonPayload))
	}
}
