package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/C9b3rD3vi1/Burnr/models"
	"log"
)


// initialize the database
var DB *gorm.DB

// ConnectDB function to connect to the database
// and migrate the schema
func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("links.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Link{}, &models.User{}, &models.LinkClick{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
}