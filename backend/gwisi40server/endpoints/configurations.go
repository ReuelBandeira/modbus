package endpoints

import (
	"encoding/json"
	"fmt"
	"gwisi40server/database"
	"gwisi40server/models"
	"gwisi40server/mqtt"
	"gwisi40server/status"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/gofiber/fiber/v2"
)

// GetConfigurations @Summary		Get all Configurations
// @Description	Get all Configurations
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Configurations
// @Failure		401				"Request not authorized"
// @Failure		404				No	Configurations	found
// @Router			/api/v1/configurations [get]
func GetConfigurations(c *fiber.Ctx) error {

	type ReturnConf struct {
		ID          uint          `json:"id"`
		Address     string        `json:"address"`
		Port        string        `json:"port"`
		Name        string        `json:"name"`
		Status      bool          `json:"status"`
		Protocol    string        `json:"protocol"`
		ReadingTime uint          `json:"readingTime"`
		Topics      []string      `json:"topics"`
		Data        []interface{} `json:"data"`
	}

	var configuration models.Configurations
	var configurations []models.Configurations

	returnString := "[{\"devices\" : [\n"

	//---------------------------------------------------------------------------------------------------------------------------
	// Proccessing of the query parameters

	page, _ := strconv.Atoi(c.Query("page", "1"))
	items, _ := strconv.Atoi(c.Query("items", "10000"))

	// Calculate offset based on page number and items per page
	start := (page - 1) * items

	var rowCount int64
	database.DB.Find(&configurations).Count(&rowCount)

	pages := int(math.Ceil(float64(rowCount) / float64(items)))

	returnString = fmt.Sprintf("{\"pages\":%d, \"total\":%d, \"items\":", pages, rowCount) + returnString

	//---------------------------------------------------------------------------------------------------------------------------
	// Query the database for configurations with pagination

	result := database.DB.Offset(start).Limit(items).Find(&configurations)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	rows, _ := result.Rows()

	for rows.Next() {
		var returnConf ReturnConf

		err := database.DB.ScanRows(rows, &configuration)
		if err != nil {
			return err
		}

		_, dev, ok := status.GetDeviceStatus(configuration.ID)

		if ok {
			returnConf.Status = dev.Status
		} else {
			returnConf.Status = false
		}

		returnConf.ID = configuration.ID
		returnConf.Address = configuration.Address
		returnConf.Port = configuration.Port
		returnConf.Name = configuration.Name
		returnConf.Protocol = configuration.Protocol
		returnConf.ReadingTime = configuration.ReadingTime

		topics := strings.Split(configuration.Topics, ",")
		returnConf.Topics = make([]string, len(topics))
		copy(returnConf.Topics, topics)

		// Get devices (came as a string) and converts them to a slice
		var result interface{}
		err = json.Unmarshal([]byte(configuration.Devices), &result)
		if err != nil {
			return err
		}

		// Convert the slice to a map
		dmap := result.(map[string]interface{})

		// Append the map to the returnConf.Data
		returnConf.Data = append(returnConf.Data, dmap)

		// Convert all the structure to a JSON string
		x, err := json.Marshal(returnConf)
		if err != nil {
			return err
		}

		returnString = returnString + string(x) + ","
	}

	c.Set("Content-Type", "application/json; charset=utf-8")

	return c.SendString(strings.TrimSuffix(returnString, ",") + "]}]}")
}

// GetConfigurations @Summary		Get all Configurations
// @Description	Get all Configurations
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Configurations
// @Failure		401				"Request not authorized"
// @Failure		404				No	Configurations	found
// @Router			/api/v1/configurationbyprotocol/{protocol} [get]
func GetConfigurationsByProtocol(c *fiber.Ctx) error {

	type ReturnConf struct {
		ID          uint          `json:"id"`
		Address     string        `json:"address"`
		Port        string        `json:"port"`
		Name        string        `json:"name"`
		Status      bool          `json:"status"`
		Protocol    string        `json:"protocol"`
		ReadingTime uint          `json:"readingTime"`
		Topics      []string      `json:"topics"`
		Data        []interface{} `json:"data"`
	}

	returnString := "[{\"devices\" : [\n"

	var configurations []models.Configurations

	protocol := c.Params("protocol")

	result := database.DB.Find(&configurations).Where("protocol = ?", protocol)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	var configuration models.Configurations

	rows, _ := database.DB.Model(&configuration).Where("protocol = ?", protocol).Rows()

	for rows.Next() {
		var returnConf ReturnConf

		err := database.DB.ScanRows(rows, &configuration)
		if err != nil {
			return err
		}

		_, dev, ok := status.GetDeviceStatus(configuration.ID)

		if ok {
			returnConf.Status = dev.Status
		} else {
			returnConf.Status = false
		}

		returnConf.ID = configuration.ID
		returnConf.Address = configuration.Address
		returnConf.Port = configuration.Port
		returnConf.Name = configuration.Name
		returnConf.Protocol = configuration.Protocol
		returnConf.ReadingTime = configuration.ReadingTime

		topics := strings.Split(configuration.Topics, ",")
		returnConf.Topics = make([]string, len(topics))
		copy(returnConf.Topics, topics)

		// Get devices (came as a string) and converts them to a slice
		var result interface{}
		err = json.Unmarshal([]byte(configuration.Devices), &result)
		if err != nil {
			return err
		}

		// Convert the slice to a map
		dmap := result.(map[string]interface{})

		// Append the map to the returnConf.Data
		returnConf.Data = append(returnConf.Data, dmap)

		// Convert all the structure to a JSON string
		x, err := json.Marshal(returnConf)
		if err != nil {
			return err
		}

		returnString = returnString + string(x) + ","
	}

	c.Set("Content-Type", "application/json; charset=utf-8")

	return c.SendString(strings.TrimSuffix(returnString, ",") + "]}]")
}

// CreateConfiguration @Summary		Creates a Configurations
// @Description	Creates a Configurations
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Configurations
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or					invalid	request	body
// @Param			type			body	models.Configurations	true	"Add Configurations"
// @Router			/api/v1/configurations [post]
func CreateConfiguration(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	var protocolVersion models.DataVersion

	result := database.DB.First(&protocolVersion, "data_type = ?", "Protocol")
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	rt, _ := strconv.Atoi(data["readingTime"])

	topics := strings.Replace(data["topics"], "[", "", -1)
	topics = strings.Replace(topics, "]", "", -1)
	topics = strings.Replace(topics, "\"", "", -1)

	devices := data["data"][1 : len(data["data"])-1]

	config := models.Configurations{
		Version:     protocolVersion.Version,
		Address:     data["address"],
		Port:        data["port"],
		Name:        data["name"],
		Protocol:    data["protocol"],
		ReadingTime: uint(rt),
		Topics:      topics,
		Devices:     devices,
	}

	database.DB.Create(&config)

	database.DB.Find(&config, "id = ?", config.ID)

	//Update Data Version
	var version []models.DataVersion
	database.DB.Model(&version).
		Where("data_type = ?", "Configurations").
		Update("version", gorm.Expr("version + ?", 1))

	type ReturnConf struct {
		ID          uint          `json:"id"`
		Address     string        `json:"address"`
		Port        string        `json:"port"`
		Name        string        `json:"name"`
		Status      bool          `json:"status"`
		Protocol    string        `json:"protocol"`
		ReadingTime uint          `json:"readingTime"`
		Topics      []string      `json:"topics"`
		Data        []interface{} `json:"data"`
	}

	var returnConf ReturnConf

	returnConf.ID = config.ID
	returnConf.Address = config.Address
	returnConf.Port = config.Port
	returnConf.Name = config.Name
	returnConf.Status = false
	returnConf.Protocol = config.Protocol
	returnConf.ReadingTime = config.ReadingTime

	ntopics := strings.Split(config.Topics, ",")
	returnConf.Topics = make([]string, len(ntopics))
	copy(returnConf.Topics, ntopics)

	// Get devices (came as a string) and converts them to a slice
	var nresult interface{}
	err := json.Unmarshal([]byte(config.Devices), &nresult)
	if err != nil {
		return err
	}

	// Convert the slice to a map
	var dmap map[string]interface{}

	if config.Devices != "" {
		// Convert the slice to a map
		dmap = nresult.(map[string]interface{})
	}

	// Append the map to the returnConf.Data
	returnConf.Data = append(returnConf.Data, dmap)

	type msgGateway struct {
		Protocol string `json:"protocol"`
	}

	m := msgGateway{
		Protocol: config.Protocol,
	}

	j, _ := json.Marshal(m)

	mqtt.Publish(string(j))

	return c.JSON(returnConf)
}

// @Summary		Get device count
// @Description	Get device count
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Configurations
// @Router			/api/v1/configurationcount [get]
func GetConfigurationCount(c *fiber.Ctx) error {

	//TODO: this endpoint will be removed in the future. It was superseeded by the /api/v1/deviceinfo endpoint

	var count int64 = 0

	database.DB.Model(&models.Configurations{}).Count(&count)

	c.Set("Content-Type", "application/json; charset=utf-8")

	return c.SendString("{\n\t\"count\":" + fmt.Sprintf("%d", count) + "\n}")
}

// @Summary		Get device count
// @Description	Get device count
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Configurations
// @Router			/api/v1/deviceinfo [get]
func GetDeviceInfo(c *fiber.Ctx) error {

	var count int64 = 0

	var diskPercent float64 = 0.0
	var memoryPercent float64 = 0.0
	var cpuPercent float64 = 0.0

	// Get info from the host
	hostStat, _ := host.Info()

	// Get info from the host disks
	if hostStat.OS == "windows" {
		var total, free uint64 = 0.0, 0.0
		partitions, _ := disk.Partitions(false)
		for _, partition := range partitions {
			v3, _ := disk.Usage(partition.Device)
			total = total + v3.Total
			free = free + v3.Free
		}
		diskPercent = (float64(total-free) / float64(total)) * 100
	} else {
		diskstat, _ := disk.Usage("/")
		diskPercent = diskstat.UsedPercent
	}

	// Get info from the host memory
	v, _ := mem.VirtualMemory()
	memoryPercent = v.UsedPercent

	// Get info from the host CPUs
	cpus := 0
	v4, _ := cpu.Percent(0, true)
	for _, cpu := range v4 {
		cpus++
		cpuPercent = cpuPercent + cpu
	}
	cpuPercent = cpuPercent / float64(cpus)

	// Get info from the configured devics
	database.DB.Model(&models.Configurations{}).Count(&count)

	c.Set("Content-Type", "application/json; charset=utf-8")

	type result struct {
		Count  int64   `json:"count"`
		Memory float64 `json:"memory"`
		CPU    float64 `json:"cpu"`
		Disk   float64 `json:"disk"`
	}

	return c.JSON(result{Count: count, Memory: memoryPercent, CPU: cpuPercent, Disk: diskPercent})
}

// GetConfigurations @Summary		Get all Configurations
// @Description	Get all Configurations with ID
// @Tags			Configurations
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Configurations
// @Failure		404				No	Configurations	found
// @Router			/api/v1/configurationwithid [get]
func GetConfigurationsWithID(c *fiber.Ctx) error {

	type Conf struct {
		ID      string `json:"id"`
		Devices string `json:"devices"`
	}

	var returnData []Conf

	var configurations []models.Configurations

	result := database.DB.Find(&configurations)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	var configuration models.Configurations
	var cf Conf

	rows, _ := database.DB.Model(&configuration).Rows()

	for rows.Next() {
		err := database.DB.ScanRows(rows, &configuration)
		if err != nil {
			return err
		}

		cf.ID = fmt.Sprintf("%d", configuration.ID)
		cf.Devices = configuration.Devices

		returnData = append(returnData, cf)
	}

	c.Set("Content-Type", "application/json; charset=utf-8")

	// return c.SendString(strings.TrimSuffix(returnString, ",") + "]")
	return c.JSON(returnData)
}

// @Summary		Deletes a configuration
// @Description	Deletes a configuration
// @Tags			Configurations
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Configurations
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	condfiguration
// @Param			id				path	int		true	"Device ID"
// @Router			/api/v1/configurations/{id} [delete]
func DeleteConfiguration(c *fiber.Ctx) error {

	var configuration models.Configurations

	result := database.DB.First(&configuration, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	result = database.DB.Delete(&models.Configurations{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	type msgGateway struct {
		Protocol string `json:"protocol"`
	}

	m := msgGateway{
		Protocol: configuration.Protocol,
	}

	j, _ := json.Marshal(m)

	mqtt.Publish(string(j))

	return nil
}

// @Summary		Updates a configuration
// @Description	Updates a configuration
// @Tags			Configurations
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Configurations
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	update	condfiguration
// @Param			id				path	int		true	"Device ID"
// @Router			/api/v1/condifgurations/{id} [patch]
func UpdateConfiguration(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	var configuration models.Configurations

	result := database.DB.First(&configuration, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	var protocolVersion models.DataVersion

	rt, _ := strconv.Atoi(data["readingTime"])

	topics := strings.Replace(data["topics"], "[", "", -1)
	topics = strings.Replace(topics, "]", "", -1)
	topics = strings.Replace(topics, "\"", "", -1)

	ndata := data["data"][1 : len(data["data"])-1]

	config := models.Configurations{
		Version:     protocolVersion.Version,
		Address:     data["address"],
		Port:        data["port"],
		Name:        data["name"],
		Protocol:    data["protocol"],
		ReadingTime: uint(rt),
		Topics:      topics,
		Devices:     ndata,
	}

	database.DB.Find(&config, "id = ?", config.ID)

	type ReturnConf struct {
		ID          uint          `json:"id"`
		Address     string        `json:"address"`
		Port        string        `json:"port"`
		Name        string        `json:"name"`
		Status      bool          `json:"status"`
		Protocol    string        `json:"protocol"`
		ReadingTime uint          `json:"readingTime"`
		Topics      []string      `json:"topics"`
		Data        []interface{} `json:"data"`
	}

	var returnConf ReturnConf

	_, dev, ok := status.GetDeviceStatus(configuration.ID)

	if ok {
		returnConf.Status = dev.Status
	} else {
		returnConf.Status = false
	}

	returnConf.ID = config.ID
	returnConf.Address = config.Address
	returnConf.Port = config.Port
	returnConf.Name = config.Name
	returnConf.Protocol = config.Protocol
	returnConf.ReadingTime = config.ReadingTime

	ntopics := strings.Split(config.Topics, ",")
	returnConf.Topics = make([]string, len(ntopics))
	copy(returnConf.Topics, ntopics)

	// Get devices (came as a string) and converts them to a slice
	var nresult interface{}
	err := json.Unmarshal([]byte(config.Devices), &nresult)
	if err != nil {
		return err
	}

	var dmap map[string]interface{}

	if config.Devices != "" {
		// Convert the slice to a map
		dmap = nresult.(map[string]interface{})
	}

	// Append the map to the returnConf.Data
	returnConf.Data = append(returnConf.Data, dmap)

	result = database.DB.Model(&configuration).Updates(config)

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	//Update Data Version
	var version []models.DataVersion
	database.DB.Model(&version).
		Where("data_type = ?", "Configurations").
		Update("version", gorm.Expr("version + ?", 1))

	type msgGateway struct {
		Protocol string `json:"protocol"`
	}

	m := msgGateway{
		Protocol: config.Protocol,
	}

	j, _ := json.Marshal(m)

	mqtt.Publish(string(j))

	return c.JSON(returnConf)
}
