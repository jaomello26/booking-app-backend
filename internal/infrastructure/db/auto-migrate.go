package db

import (
	"log"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.UserGroup{},
		&models.Apartment{},
		&models.Booking{},
	)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
		return err
	}

	return nil
}
