package endpoints

import (
	"encoding/json"
	"fmt"
	"gwisi40server/database"
	"gwisi40server/models"

	"gorm.io/gorm/clause"

	"github.com/gofiber/fiber/v2"
)

// GetAllProtocols @Summary		Get all protocols
// @Description	Get all protocols
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Protocol
// @Failure		401				"Request not authorized"
// @Failure		404				No	protocols	found
// @Router			/api/v1/protocols [get]
func GetAllProtocols(c *fiber.Ctx) error {
	var protocols []models.Protocol

	result := database.DB.Find(&protocols)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(protocols)
}

// GetProtocol @Summary		Get protocol by ID
// @Description	Get protocol by ID
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Protocol
// @Failure		401				"Request not authorized"
// @Failure		404				No		protocol	found
// @Param			id				path	int			true	"Protocol ID"
// @Router			/api/v1/protocol/{id} [get]
func GetProtocol(c *fiber.Ctx) error {
	var protocol models.Protocol

	result := database.DB.First(&protocol, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(protocol)
}

// GetProtocol @Summary		Get protocol by name
// @Description	Get protocol by name
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Protocol
// @Failure		404				No		protocol	found
// @Param			protocol				path	string			true	"Protocol ID"
// @Router			/api/v1/protocol/{protocol} [get]
func GetProtocolByName(c *fiber.Ctx) error {

	var protocol models.Protocol

	result := database.DB.First(&protocol, "protocol = ?", c.Params("protocol"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(protocol)
}

// CreateProtocol @Summary		Creates a protocol
// @Description	Creates a protocol
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Protocol
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or						invalid	request	body
// @Param			type			body	modelsswagger.Protocol	true	"Add protocol"
// @Router			/api/v1/protocol [post]
func CreateProtocol(c *fiber.Ctx) error {

	type body struct {
		Protocol string `json:"protocol"`
		Version  uint   `json:"version"`
		Alias    string `json:"alias"`
		Md5      string `json:"md5"`
		Config   string `json:"config"`
	}

	var data body

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	protocol := models.Protocol{
		Protocol: data.Protocol,
		Version:  data.Version,
		Alias:    data.Alias,
		MD5:      data.Md5,
		Config:   data.Config,
	}

	database.DB.Create(&protocol)

	return c.JSON(protocol)
}

// DeleteProtocol @Summary		Deletes a protocol
// @Description	Deletes a protocol
// @Tags			Protocols
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Protocol
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	protocol
// @Param			id				path	int		true	"Protocol ID"
// @Router			/api/v1/protocol/{id} [delete]
func DeleteProtocol(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.Protocol{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}

// CreateProtocols @Summary		Replaces all protocols
// @Description	Creates a protocol
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Protocol
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or invalid request body
// @Param			type			body	modelsswagger.Protocol	true	"Add protocol"
// @Router			/api/v1/protocols [post]
func CreateProtocols(c *fiber.Ctx) error {

	data := struct {
		Version   uint `json:"version"`
		Protocols []struct {
			Protocol string `json:"protocol"`
			Alias    string `json:"alias"`
			MD5      string `json:"md5"`
			Config   any    `json:"config"`
		} `json:"protocols"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	for _, value := range data.Protocols {
		var protocol models.Protocol
		protocol.Protocol = value.Protocol
		out, _ := json.Marshal(&value.Config)
		protocol.Config = string(out)
		protocol.Version = data.Version

		database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "protocol"}},
			DoUpdates: clause.AssignmentColumns([]string{"config", "version"}),
		}).Create(&protocol)
	}

	//Update Data Version
	var version []models.DataVersion
	database.DB.Model(&version).
		Where("data_type = ?", "Protocol").
		Update("version", data.Version)

	return c.JSON(data)
}

// UpdateProtocol @Summary		Updates a protocol
// @Description	Updates a protocol
// @Tags			Protocols
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.Protocol
// @Failure		400				Missing	or						invalid	request	body
// @Param			type			body	modelsswagger.Protocol	true	"Add protocol"
// @Router			/api/v1/protocol/{protocol} [patch]
func UpdateProtocol(c *fiber.Ctx) error {

	type body struct {
		Protocol string `json:"protocol"`
		Version  uint   `json:"version"`
		Alias    string `json:"alias"`
		Md5      string `json:"md5"`
		Config   string `json:"config"`
	}

	var data body

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Error on BODY PARSER: %v\n", err)
		return fiber.ErrBadRequest
	}

	var protocol models.Protocol

	result := database.DB.First(&protocol, "protocol = ?", c.Params("protocol"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	protocol.Protocol = data.Protocol
	protocol.MD5 = data.Md5
	protocol.Alias = data.Alias
	protocol.Version = data.Version
	protocol.Config = data.Config

	database.DB.Debug().Save(&protocol)

	return c.JSON(protocol)
}
