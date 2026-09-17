package process

import (
	"gateway/globals"
	"gateway/mqtt"
	"path/filepath"

	"log"
	"os"
	"os/exec"
	"strconv"
)

func Initialize() {

	for _, protocol := range globals.GatewayProtocols {
		go Start(protocol.Exec, protocol.Protocol)
	}
}

func Start(process string, name string) {

	workingDir := globals.Conf.Protocols.Directory + "/" + name

	executable, err := filepath.Abs(process)
	if err != nil {
		log.Printf("Error getting absolute path: %s\n", err)
	}

	log.Printf("Starting process: %s, Name: %s\n", executable, name)

	cmd := exec.Command(executable)
	cmd.Dir = workingDir

	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting process: %s\n", err)
		return
	}

	mqtt.SendMsg(globals.LocalClient, "gateway", "Protocol "+name+" started - PID: "+strconv.Itoa(cmd.Process.Pid))
	log.Printf("Protocol %s started with PID: %d\n", name, cmd.Process.Pid)

	// When received a message from the protocol on the protocolname topic, send it to the external MQTT broker
	// And when I receive a message from the external MQTT broker on the protocolname topic, send it to the local MQTT broker
	mqtt.MQTTProtocolMessages(name)
}

func Stop(process string) {

	temp := globals.Processes[:0]

	for _, p := range globals.Processes {
		log.Printf("Required to STOP %s; Process: %s\n", process, p.Name)
		if p.Name == process {
			p.Cmd.Process.Kill()
			p.Active = false
			mqtt.SendMsg(globals.LocalClient, "gateway", "Protocol "+process+" stopped")
			log.Printf("Protocol %s stopped\n", process)
		} else {
			temp = append(temp, p)
		}
	}

	globals.Processes = temp
}
