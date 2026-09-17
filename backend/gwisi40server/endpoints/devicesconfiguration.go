package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type IDevice struct {
	ID           uint        `json:"-"`
	Manufacturer string      `json:"fabricante"`
	Serie        string      `json:"serie"`
	Memories     []IMemories `json:"memorias"`
}

type IMemories struct {
	ID     uint      `json:"-"`
	Label  string    `json:"label"`
	Type   string    `json:"tipo"`
	Format string    `json:"formato"`
	Values []IValues `json:"valores"`
}

type IValues struct {
	ID    uint   `json:"-"`
	Key   string `json:"name"`
	Value string `json:"address"`
}

func GetDeviceConfiguration(c *fiber.Ctx) error {

	var deviceConfiguration models.DeviceConfiguration
	var memories []models.Memories
	var values []models.Values

	var idevice IDevice
	var imemories []IMemories
	var ivalues []IValues

	var deviceID uint
	var memoryID uint

	// Read the device configuration
	result := database.DB.Find(&deviceConfiguration)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	deviceID = deviceConfiguration.ID

	// Read the memories of the device configuration
	rows, _ := database.DB.Model(&memories).Where("device_id = ?", deviceID).Rows()

	var imemory models.Memories

	for rows.Next() {
		database.DB.ScanRows(rows, &imemory)
		memory := IMemories{Label: imemory.Label, Type: imemory.Type, Format: imemory.Format}

		memoryID = imemory.ID

		vals, err := database.DB.Model(&values).Where("memory_id = ?", memoryID).Rows()

		if err != nil {
			log.Fatalln(err)
		}

		var ivalue models.Values

		for vals.Next() {
			database.DB.ScanRows(vals, &ivalue)
			value := IValues{Key: ivalue.Key, Value: ivalue.Value}
			ivalues = append(ivalues, value)
		}

		memory.Values = ivalues
		imemories = append(imemories, memory)
	}

	idevice = IDevice{
		Manufacturer: deviceConfiguration.Manufacturer,
		Serie:        deviceConfiguration.Serie,
		Memories:     imemories,
	}

	return c.JSON(idevice)
}
