package routes

import (
	"encoding/base64"
	"gwisi40server/database"
	"gwisi40server/endpoints"
	"gwisi40server/models"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Setup sets up the routes
func Setup(app *fiber.App) {

	app.Post("/api/v1/login", endpoints.Login)
	app.Post("/api/v1/register", endpoints.CreateUser)
	app.Get("/api/v1/logs", endpoints.GetAllLogs)
	app.Get("/api/v1/logfile/:id", endpoints.DownloadLog)
	app.Get("/api/v1/logfiles", endpoints.DownloadAllLogs)

	protected := app.Group("/api/v1/")

	// Define a JWT configuration
	jwtConfig := jwtware.Config{
		SigningKey: []byte(os.Getenv("API_SECRET")),
	}

	// JWT middleware
	jwtMiddleware := jwtware.New(jwtConfig)

	// Middleware chain: API key validation or JWT authentication
	protected.Use(func(c *fiber.Ctx) error {
		// Attempt API key validation
		apik, _ := base64.StdEncoding.DecodeString(c.Get("API-key"))
		if strings.HasPrefix(string(apik), "gateway-") {
			return c.Next()
		}
		// If API key validation fails, proceed to JWT validation
		return jwtMiddleware(c)
	})

	// Add data version
	protected.Use(func(c *fiber.Ctx) error {
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

		return c.Next()
	})

	protected.Put("user/:id", endpoints.UpdateUser)
	protected.Patch("userdata/:id", endpoints.UpdateUserData)

	protected.Get("language", endpoints.GetLanguage)

	protected.Get("mqtt", endpoints.GetMqttSettings)
	protected.Get("mqttstatus", endpoints.GetMqttStatus)

	protected.Get("deviceconfiguration", endpoints.GetDeviceConfiguration)

	protected.Get("configurations", endpoints.GetConfigurations)
	protected.Get("configurationsbyprotocol/:protocol", endpoints.GetConfigurationsByProtocol)
	protected.Get("configurationwithid", endpoints.GetConfigurationsWithID)
	protected.Get("deviceinfo", endpoints.GetDeviceInfo)
	protected.Get("configurationcount", endpoints.GetConfigurationCount)

	protected.Get("resources", endpoints.GetResources)
	protected.Post("resources", endpoints.CreateResources)

	protected.Get("version", endpoints.GetDataVersion)

	protected.Get("protocols/:protocol", endpoints.GetProtocolByName)
	protected.Post("protocol", endpoints.CreateProtocol)
	protected.Patch("protocol/:protocol", endpoints.UpdateProtocol)
	protected.Post("protocols", endpoints.CreateProtocols)

	protected.Post("mqtt", endpoints.CreateMqttSetting)
	protected.Delete("mqtt/:id", endpoints.DeleteMqttSetting)

	protected.Post("configurations", endpoints.CreateConfiguration)
	protected.Delete("configurations/:id", endpoints.DeleteConfiguration)
	protected.Patch("configurations/:id", endpoints.UpdateConfiguration)

	protected.Get("servers", endpoints.GetAllServers)
	protected.Get("server/:id", endpoints.GetServer)
	protected.Post("server", endpoints.CreateServer)
	protected.Delete("server/:id", endpoints.DeleteServer)

	protected.Get("types", endpoints.GetAllTypes)
	protected.Get("type/:id", endpoints.GetType)
	protected.Post("type", endpoints.CreateType)
	protected.Delete("type/:id", endpoints.DeleteType)

	protected.Get("protocols", endpoints.GetAllProtocols)
	protected.Get("protocol/:id", endpoints.GetProtocol)
	// protected.Post("protocol", endpoints.CreateProtocol)
	protected.Delete("protocol/:id", endpoints.DeleteProtocol)

	protected.Get("plcs", endpoints.GetAllPlcs)
	protected.Get("plc/:id", endpoints.GetPlc)
	protected.Post("plc", endpoints.CreatePlc)
	protected.Delete("plc/:id", endpoints.DeletePlc)

	protected.Get("devices", endpoints.GetAllDevices)
	protected.Get("device/:id", endpoints.GetDevice)
	protected.Post("device", endpoints.CreateDevice)
	protected.Delete("device/:id", endpoints.DeleteDevice)
}
