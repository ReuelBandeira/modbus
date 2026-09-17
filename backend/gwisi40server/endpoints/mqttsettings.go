package endpoints

import (
	"encoding/json"
	"gwisi40server/database"
	"gwisi40server/models"
	"gwisi40server/mqtt"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// GetMqttSettings @Summary		Get all MQTT Settings
// @Description	Get all MQTT Settings
// @Tags			MQTT
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.MQTTSettings
// @Failure		401				"Request not authorized"
// @Failure		404				No	MQQT Settings	found
// @Router			/api/v1/mqtt [get]
func GetMqttSettings(c *fiber.Ctx) error {

	var mqtt models.MQTTSettings

	result := database.DB.Find(&mqtt)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(mqtt)
}

// GetMqttStatus @Summary		Get MQTT Status
// @Description	Get MQTT Status
// @Tags			MQTT
// @Accept			json
// @Produce		json
//
// @Success		200				{object}	models.MQTTSettings
// @Failure		401				"Request not authorized"
// @Failure		404				No	MQQT Settings	found
// @Router			/api/v1/mqttstatus [get]
func GetMqttStatus(c *fiber.Ctx) error {

	var mqttx models.MQTTSettings

	result := database.DB.Find(&mqttx)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	mqttx.Connected = mqtt.VerifyBrokerConnection()

	database.DB.Model(&mqttx).Where("id = ?", 1).Updates(&mqttx)

	return c.JSON(mqttx)
}

// CreateMqttSetting @Summary		Creates a MQTT Setting
// @Description	Creates a MQTT Setting
// @Tags			MQTT
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.MQTTSettings
// @Failure		401				"Request not authorized"
// @Failure		400				Missing	or					invalid	request	body
// @Param			type			body	models.MQTTSettings	true	"Add MQQT Setting"
// @Router			/api/v1/mqtt [post]
func CreateMqttSetting(c *fiber.Ctx) error {

	var mqttl []models.MQTTSettings

	var mqttnew models.MQTTSettings
	var data map[string]string

	result := database.DB.Find(&mqttl)

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	if result.RowsAffected == 0 {

		mqttnew = models.MQTTSettings{
			ID:           1,
			Server:       data["server"],
			Port:         data["port"],
			Connected:    false,
			Username:     data["username"],
			Password:     data["password"],
			InputInfo:    data["inputInfo"],
			ErrorInfo:    data["errorInfo"],
			StatusDevice: data["statusDevice"],
		}

		database.DB.Create(&mqttnew)
	} else {
		mqttnew = models.MQTTSettings{
			ID:           1,
			Server:       data["server"],
			Port:         data["port"],
			Connected:    false,
			Username:     data["username"],
			Password:     data["password"],
			InputInfo:    data["inputInfo"],
			ErrorInfo:    data["errorInfo"],
			StatusDevice: data["statusDevice"],
		}
		database.DB.Model(&mqttl).Where("id = ?", 1).Updates(&mqttnew)
	}

	type msgGateway struct {
		Protocol string `json:"protocol"`
	}

	m := msgGateway{
		Protocol: "mqtt",
	}

	j, _ := json.Marshal(m)

	mqtt.Publish(string(j))

	//Update Data Version
	var version []models.DataVersion
	database.DB.Model(&version).
		Where("data_type = ?", "MQTT").
		Update("version", gorm.Expr("version + ?", 1))

	mqtt.InitializeExternalBrokerConnection()

	return c.JSON(mqttnew)
}

// DeleteMqttSetting @Summary		Deletes a MQTT Setting
// @Description	Deletes a MQTT Setting
// @Tags			MQTT
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.MQTTSettings
// @Failure		401				"Request not authorized"
// @Failure		422				Cannot	delete	MQTT Setting
// @Param			id				path	int		true	"MQTT Setting ID"
// @Router			/api/v1/mqtt/{id} [delete]
func DeleteMqttSetting(c *fiber.Ctx) error {

	result := database.DB.Delete(&models.MQTTSettings{}, c.Params("id"))

	if result.RowsAffected == 0 {
		return fiber.ErrUnprocessableEntity
	}

	return nil
}
