package endpoints

import (
	"fmt"
	"gwisi40server/database"
	"gwisi40server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get all devices
// @Description	Get all devices
// @Tags			Devices
// @Security		BasicAuth
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Device
// @Failure		401				"Request not authorized"
// @Failure		404				"No	devices	found"
// @Router			/api/v1/devices [get]
func GetAllDevices(c *fiber.Ctx) error {
	var devices []models.Device

	result := database.DB.Preload("PLC").Preload("Type").Preload("Protocol").Find(&devices)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(devices)
}

// @Summary		Get device by ID
// @Description	Get device by ID
// @Tags			Devices
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Device
// @Failure		401				"Request not authorized"
// @Failure		404				"No			device found"
// @Param			id				path	int	true	"Device ID"
// @Router			/api/v1/device/{id} [get]
func GetDevice(c *fiber.Ctx) error {
	var device models.Device

	result := database.DB.Preload("PLC").Preload("Type").Preload("Protocol").First(&device, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(device)
}

// @Summary		Creates a device
// @Description	Creates a device
// @Tags			Devices
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Device
// @Failure		400				Missing		or	invalid	request	body
// @Failure		401				"Request not authorized"
// @Failure		404				"No		protocol, type or PLC found or invalid values for name, port, memory block, class, instance or attribute"
// @Param			type			body	modelsswagger.Device	true	"Add device"
// @Router			/api/v1/device [post]
func CreateDevice(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	protocolID, _ := strconv.Atoi(data["protocolID"])
	plcID, _ := strconv.Atoi(data["plcid"])
	typeID, _ := strconv.Atoi(data["typeID"])

	memoryBlock, _ := strconv.Atoi(data["memoryBlock"])
	class, _ := strconv.Atoi(data["class"])
	instance, _ := strconv.Atoi(data["instance"])
	attribute, _ := strconv.Atoi(data["attribute"])
	address, _ := strconv.Atoi(data["address"])

	result := database.DB.Find(&models.Protocol{}, "id = ?", protocolID)
	if result.RowsAffected == 0 {
		fmt.Printf("Protocol %d not found\n", protocolID)
		return fiber.ErrNotFound
	}

	result = database.DB.Find(&models.PLC{}, "id = ?", plcID)
	if result.RowsAffected == 0 {
		fmt.Printf("PLC %d not found\n", plcID)
		return fiber.ErrNotFound
	}

	result = database.DB.Find(&models.Type{}, "id = ?", typeID)
	if result.RowsAffected == 0 {
		fmt.Printf("Type %d not found\n", typeID)
		return fiber.ErrNotFound
	}

	device := models.Device{
		Name:        data["name"],
		ProtocolID:  uint(protocolID),
		PLCID:       uint(plcID),
		TypeID:      uint(typeID),
		MemoryBlock: memoryBlock,
		Class:       class,
		Instance:    instance,
		Attribute:   attribute,
		Address:     address,
	}

	database.DB.Create(&device)

	database.DB.Preload("PLC").Preload("Type").Preload("Protocol").First(&device, "id = ?", device.ID)

	return c.JSON(device)
}

// @Summary		Deletes a device
// @Description	Deletes a device
// @Tags			Devices
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.Device
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	device
// @Param			id				path	int		true	"Device ID"
// @Router			/api/v1/device/{id} [delete]
func DeleteDevice(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.Device{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
