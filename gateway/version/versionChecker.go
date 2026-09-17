package version

import (
	"fmt"
	"gateway/globals"
	"gateway/models"
	"gateway/protocols"
	"runtime"

	"bytes"
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Protocol struct {
	ID       uint   `json:"-"`
	Protocol string `json:"protocol"`
	Alias    string `json:"alias"`
	Version  uint   `json:"version"`
	MD5      string `json:"md5"`
	Config   string `json:"config"`
}

var (
	allprotocols   []Protocol
	protocolName   string
	executableName string
)

// ================================================================================================
// Calculate MD5 checksum of the protocol configuration file
// ================================================================================================
func calculateMD5(protocol string, filename string) (string, error) {

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", err
	}

	md5 := fmt.Sprintf("%x", string(hash.Sum(nil)))

	return md5, nil

}

// ================================================================================================
// Search on the protocols folder for protocols
// ================================================================================================
func searchProtocols() {

	err := filepath.Walk(globals.Conf.Protocols.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("%v\n", err)
			}
			if info.IsDir() && strings.Count(path, globals.PathSeparator) == 2 {
				searchExecutables(path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

}

// ================================================================================================
// Search on each protocol folder for executables
// ================================================================================================
func searchExecutables(path string) {

	protocolName = path[strings.LastIndex(path, globals.PathSeparator)+1:]

	switch runtime.GOOS {
	case "windows":
		executableName = path + globals.PathSeparator + protocolName + ".exe"
	case "linux":
		executableName = path + globals.PathSeparator + protocolName
	}

	if _, err := os.Stat(executableName); err != nil {
		return
	}

	searchProtocolConfigurations(protocolName, executableName)

}

// ================================================================================================
// Search on each protocol folder for protocol_config.json
// ================================================================================================
func searchProtocolConfigurations(protocolname string, executablename string) {

	configuration := globals.Conf.Protocols.Directory + "/" + protocolname + "/protocol_config.json"

	if _, err := os.Stat(configuration); err != nil {
		return
	}

	md5, _ := calculateMD5(protocolname, configuration)

	b, err := os.ReadFile(configuration)
	if err != nil {
		log.Printf("FAIL reading protocol configuration file [%s]\n", configuration)
		return
	}

	config := string(b)

	// Parses the JSON-encoded data and stores the result in the new Protocol struct
	type configJSON struct {
		Protocol string
		Alias    string
		Config   map[string]interface{} `json:"config"`
	}

	newconfig := new(configJSON)

	err = json.Unmarshal([]byte(config), &newconfig)
	if err != nil {
		log.Printf("FAIL unmarshaling JSON for protocol [%s]: %v\n", protocolname, err)
		return
	}

	jsonStr, err := json.Marshal(newconfig.Config)
	if err != nil {
		log.Printf("FAIL marshaling Extra Data to JSON for protocol [%s]: %v\n", protocolname, err)
		return
	}

	p := Protocol{Protocol: protocolname, Alias: newconfig.Alias, Version: 1, MD5: md5, Config: string(jsonStr)}
	allprotocols = append(allprotocols, p)

	p2 := models.GatewayProtocol{
		Protocol: protocolname,
		Exec:     executablename,
	}

	globals.GatewayProtocols = append(globals.GatewayProtocols, p2)

}

// ================================================================================================
// Get protocol from Backend
// ================================================================================================
func getProtocol(proto Protocol) (Protocol, int, error) {

	var currentProtocol Protocol

	client := &http.Client{}

	req, _ := http.NewRequest("GET", globals.Conf.Backend.Host+"/protocols/"+proto.Protocol, nil)

	req.Header.Set("Api-key", globals.GenerateApiKey())

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("FAIL to read protocols from Backend: %s\n", err)
		return Protocol{}, -1, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("FAIL to read already existent protocol from Backend.\n")
		return Protocol{}, resp.StatusCode, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FAIL to read response body.\n")
	}

	err = json.Unmarshal(body, &currentProtocol)
	if err != nil {
		log.Printf("FAIL unmarshaling JSON for Protocols: %v\n", err)
		return proto, resp.StatusCode, err
	}

	return currentProtocol, resp.StatusCode, nil

}

// ================================================================================================
// Send protocol to Backend
// ================================================================================================
func sendToBackend(protocol Protocol, insert bool) {

	p, err := json.Marshal(protocol)
	if err != nil {
		log.Printf("FAIL marshaling JSON for Protocols: %v\n", err)
		return
	}

	var req *http.Request

	if !insert {
		log.Printf("Updating protocol [%s]\n", protocol.Protocol)
		req, err = http.NewRequest("PATCH", globals.Conf.Backend.Host+"/protocol/"+protocol.Protocol, bytes.NewReader(p))
	} else {
		log.Printf("Creating protocol [%s]\n", protocol.Protocol)
		req, err = http.NewRequest("POST", globals.Conf.Backend.Host+"/protocol", bytes.NewReader(p))
	}

	if err != nil {
		log.Printf("FAIL to create protocol request to Backend.\n")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-key", globals.GenerateApiKey())

	client := http.Client{Timeout: 10 * time.Second}

	_, err = client.Do(req)

	if err != nil {
		log.Printf("FAIL to create protocol on Backend.\n")
	}

}

// ================================================================================================
// Integrate protocol on Backend if necessary
// ================================================================================================
func integrateProtocol(protocol Protocol) {

	existentProtocolConfig, status, err := getProtocol(protocol)

	if err != nil {
		log.Printf("FAIL to get protocol from Backend.\n")
		return
	}

	switch status {

	case 200:

		if existentProtocolConfig.MD5 != protocol.MD5 {
			log.Printf("ProtocolConfig for Protocol [%s] updated\n", protocol.Protocol)
			sendToBackend(protocol, false)
		}

	case 404:

		log.Printf("ProtocolConfig for Protocol [%s] cerated on Backend\n", protocol.Protocol)

		sendToBackend(protocol, true)

	default:

		log.Printf("Protocol [%s] has an unexpected status code [%d]\n", protocol.Protocol, status)
	}

	protocols.GetBackendDevSettings(protocol.Protocol)

}

func CheckVersions() {
	globals.Dictionary = make(map[string]string)

	searchProtocols()
	for _, protocol := range allprotocols {
		globals.Dictionary[protocol.Protocol] = protocol.Protocol
		integrateProtocol(protocol)
	}
}
