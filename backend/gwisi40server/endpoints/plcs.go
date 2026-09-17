package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get all PLCs
// @Description	Get all PLCs
// @Tags			PLCs
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.PLC
// @Failure		401				"Request not authorized"
// @Failure		404				No	PLCs	found
// @Router			/api/v1/plcs [get]
func GetAllPlcs(c *fiber.Ctx) error {

	var plcs []models.PLC

	result := database.DB.Find(&plcs)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(plcs)
}

// @Summary		Get PLC by ID
// @Description	Get PLC by ID
// @Tags			PLCs
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.PLC
// @Failure		401				"Request not authorized"
// @Failure		404				No		PLC	found
// @Param			id				path	int	true	"PLC ID"
// @Router			/api/v1/plc/{id} [get]
func GetPlc(c *fiber.Ctx) error {
	var plc models.PLC

	result := database.DB.First(&plc, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(plc)
}

// @Summary		Creates a PLC
// @Description	Creates a PLC
// @Tags			PLCs
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.PLC
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or					invalid	request	body
// @Param			type			body	modelsswagger.PLC	true	"Add PLC"
// @Router			/api/v1/plc [post]
func CreatePlc(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	plc := models.PLC{
		Name: data["name"],
	}

	database.DB.Create(&plc)

	database.DB.Find(&plc, "id = ?", plc.ID)

	return c.JSON(plc)
}

// @Summary		Deletes a PLC
// @Description	Deletes a PLC
// @Tags			PLCs
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.PLC
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	PLC
// @Param			id				path	int		true	"PLC ID"
// @Router			/api/v1/plc/{id} [delete]
func DeletePlc(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.PLC{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
