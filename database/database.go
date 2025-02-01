package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lokesh2201013/email-service/models"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=9910994194lokesh dbname=EMAILID port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables
	db.AutoMigrate(&models.Sender{},&models.Template{})

	DB = db
}
