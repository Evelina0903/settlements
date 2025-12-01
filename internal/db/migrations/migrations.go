package migrations

import (
	"settlements/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Type{}, &models.District{}, &models.City{})
	return err
}
