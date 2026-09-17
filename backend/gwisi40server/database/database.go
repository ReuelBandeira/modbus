package database

import (
	"gwisi40server/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Our database
var DB *gorm.DB

// ConnectDB - returns a pointer to a new database connection
func ConnectDB(initialize bool) error {
	log.Println("Opening connection to database")

	db, err := gorm.Open(sqlite.Open("gwisi40.db?_pragma=foreign_keys(1)"), &gorm.Config{})

	if err != nil {
		return err
	}

	log.Println("Migrating database tables")

	if result := db.AutoMigrate(
		&models.Protocol{},
		&models.Server{},
		&models.Device{},
		&models.Type{},
		&models.PLC{},
		&models.User{},
		&models.DeviceConfiguration{},
		&models.Memories{},
		&models.MQTTSettings{},
		&models.Configurations{},
		&models.DataVersion{},
		&models.Values{},
		&models.Resources{},
		&models.Language{}); result != nil {
		return result
	}

	DB = db

	if initialize {
		InitializeDB()
	}

	return nil
}
