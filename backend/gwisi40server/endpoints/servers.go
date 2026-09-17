package endpoints

import (
	"gwisi40server/database"
	"gwisi40server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Get all servers
// @Description	Get all servers
// @Tags			Servers
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Server
// @Failure		401				"Request not authorized"
// @Failure		404				No	servers	found
// @Router			/api/v1/servers [get]
func GetAllServers(c *fiber.Ctx) error {

	var servers []models.Server

	result := database.DB.Model(servers).Preload("Protocol").Find(&servers)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(servers)
}

// @Summary		Get server by ID
// @Description	Get server by ID
// @Tags			Servers
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Server
// @Failure		401				"Request not authorized"
// @Failure		404				No		server	found
// @Param			id				path	int		true	"Server ID"
// @Router			/api/v1/server/{id} [get]
func GetServer(c *fiber.Ctx) error {
	var server models.Server

	result := database.DB.Preload("Protocol").First(&server, "id = ?", c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(server)
}

// @Summary		Creates a server
// @Description	Creates a server
// @Tags			Servers
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	modelsswagger.Server
// @Failure		400				Missing		or	invalid	request	body
// @Failure		401				"Request not authorized"
// @Failure		404				No		protocol				found	or	invalid	IP	or	Port	or	Name
// @Param			type			body	modelsswagger.Server	true	"Add server"
// @Router			/api/v1/server [post]
func CreateServer(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	port, _ := strconv.Atoi(data["port"])
	protocolID, _ := strconv.Atoi(data["protocolID"])

	result := database.DB.Find(&models.Protocol{}, "id = ?", protocolID)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	server := models.Server{
		IP:         data["ip"],
		Name:       data["name"],
		Port:       port,
		ProtocolID: uint(protocolID),
	}

	database.DB.Create(&server)

	database.DB.Preload("Protocol").First(&server, "id = ?", server.ID)

	return c.JSON(server)
}

// @Summary		Deletes a server
// @Description	Deletes a server
// @Tags			Servers
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.Server
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	server
// @Param			id				path	int		true	"Server ID"
// @Router			/api/v1/server/{id} [delete]
func DeleteServer(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.Server{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
