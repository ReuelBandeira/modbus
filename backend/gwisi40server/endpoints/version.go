package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetDataVersion @Summary		Get data version
// @Description	Get data version
// @Tags			DataVersion
// @Security		BasicAuth
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Authentication header"
// @Success		200				{object}	models.DataVersion
// @Failure		401				"Request not authorized"
// @Failure		404				"No	devices	found"
// @Router			/api/v1/version [get]
func GetDataVersion(c *fiber.Ctx) error {

	var version []models.DataVersion

	result := database.DB.Find(&version)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	rows, _ := database.DB.Model(&version).Rows()
	var item models.DataVersion

	for rows.Next() {
		err := database.DB.ScanRows(rows, &item)
		if err != nil {
			return err
		}
		switch item.DataType {
		case "Protocol":
			c.Set("X-Custom-Protocol-Version", strconv.Itoa(int(item.Version)))
		case "MQTT":
			c.Set("X-Custom-MQTT-Version", strconv.Itoa(int(item.Version)))
		case "Configurations":
			c.Set("X-Custom-Configurations-Version", strconv.Itoa(int(item.Version)))
		}
	}

	return c.JSON(version)
}
