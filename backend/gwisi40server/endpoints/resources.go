package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Creates resources info
// @Description	Creates resources
// @Tags			Resources
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Resources
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or					invalid	request	body
// @Param			type			body	models.Resources	true	"Add resource info"
// @Router			/api/v1/resources [post]
func CreateResources(c *fiber.Ctx) error {

	var resources []models.Resources

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	mem, err := strconv.ParseFloat(data["memoryUsage"], 64)

	if err != nil {
		return fiber.ErrBadRequest
	}

	disk, err := strconv.ParseFloat(data["diskUsage"], 64)

	if err != nil {
		return fiber.ErrBadRequest
	}

	res := models.Resources{
		ID:               1,
		MemoryUsage:      mem,
		DiskUsage:        disk,
		ConnectedDevices: data["connectedDevices"],
	}

	result := database.DB.Find(&resources)
	if result.RowsAffected == 0 {
		database.DB.Create(&res)
	} else {
		database.DB.Model(&resources).Where("id = ?", 1).Updates(res)
	}

	return c.JSON(res)
}

// @Summary		Get resources used
// @Description	Get resources used
// @Tags			Resources
// @Produce		json
//
// @Success		200				{object}	models.Resources
// @Router			/api/v1/resources [get]
func GetResources(c *fiber.Ctx) error {

	var resources []models.Resources

	result := database.DB.Find(&resources)

	if result.RowsAffected == 0 {
		res := models.Resources{
			ID:               1,
			MemoryUsage:      0.0,
			DiskUsage:        0.0,
			ConnectedDevices: "",
		}
		return c.JSON(res)
	}

	return c.JSON(resources)
}
