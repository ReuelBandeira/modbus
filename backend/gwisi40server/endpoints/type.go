package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get all types
// @Description	Get all types
// @Tags			Types
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Type
// @Failure		401				"Request not authorized"
// @Failure		404				No	types	found
// @Router			/api/v1/types [get]
func GetAllTypes(c *fiber.Ctx) error {
	var types []models.Type

	result := database.DB.Find(&types)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(types)
}

// @Summary		Get type by ID
// @Description	Get type by ID
// @Tags			Types
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Type
// @Failure		401				"Request not authorized"
// @Failure		404				No		type	found
// @Param			id				path	int		true	"Type ID"
// @Router			/api/v1/type/{id} [get]
func GetType(c *fiber.Ctx) error {
	var types models.Type

	result := database.DB.First(&types, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(types)
}

// @Summary		Creates a type
// @Description	Creates a type
// @Tags			Types
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Type
// @Failure		400				Missing		or	invalid	request	body
// @Failure		401				"Request not authorized"
// @Param			type			body	modelsswagger.Type	true	"Add type"
// @Router			/api/v1/type [post]
func CreateType(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	types := models.Type{
		Type: data["type"],
	}

	database.DB.Create(&types)

	return c.JSON(types)
}

// @Summary		Deletes a type
// @Description	Deletes a type
// @Tags			Types
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Type
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	type
// @Param			id				path	int		true	"Type ID"
// @Router			/api/v1/type/{id} [delete]
func DeleteType(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.Type{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
