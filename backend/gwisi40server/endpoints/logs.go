package endpoints

import (
	"gwisi40server/models"
	"gwisi40server/services"
	"math"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Gets a list of log files
// @Description	Gets a list of log files
// @Tags			Logs
// @Success		200				{object} []models.Logfile
// @Router			/api/v1/logs [get]
func GetAllLogs(c *fiber.Ctx) error {

	type ReturnStruct struct {
		Pages int              `json:"pages"`
		Total int              `json:"total"`
		Items []models.Logfile `json:"items"`
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	items, _ := strconv.Atoi(c.Query("items", "10000"))

	logs := services.LoadLogFiles()

	pages := int(math.Ceil(float64(len(logs)) / float64(items)))

	start := (page - 1) * items
	end := page * items

	if start > len(logs) {
		return c.JSON([]string{})
	}

	if end > len(logs) {
		end = len(logs)
	}

	return c.JSON(ReturnStruct{Pages: int(pages), Total: len(logs), Items: logs[start:end]})
}

// @Summary		Download a log file
// @Description	Download a log file
// @Tags			Logs
// @Router			/api/v1/logfile/{id} [get]
func DownloadLog(c *fiber.Ctx) error {

	logs := services.LoadLogFiles()

	id, _ := strconv.Atoi(c.Params("id"))

	fname := logs[id].FullName

	return c.Download(fname, filepath.Base(fname))
}

// @Summary		Download all log files
// @Description	Download all log files
// @Tags			Logs
// @Router			/api/v1/logfiles [get]
func DownloadAllLogs(c *fiber.Ctx) error {

	fname := services.ZipLogFiles()

	return c.Download(fname)
}
