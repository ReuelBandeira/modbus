package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"

	"github.com/gofiber/fiber/v2"
)

// Getlanguage @Summary		Get all language
// @Description	Get all language
// @Tags			language
// @Produce		json
//
// @Success		200				{object}	models.Language
// @Failure		401				"Request not authorized"
// @Failure		404				No	language	found
// @Router			/api/v1/language [get]
func GetLanguage(c *fiber.Ctx) error {

	var language []models.Language

	result := database.DB.Find(&language)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(language)
}
