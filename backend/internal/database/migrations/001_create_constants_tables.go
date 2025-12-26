package migrations

import (
	"dental-marketplace/backend/internal/models"

	"gorm.io/gorm"
)

func CreateConstantsTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Role{},
		&models.Specialization{},
		&models.TreatmentStatus{},
		&models.OfferStatus{},
		&models.AppointmentStatus{},
		&models.ScanStatus{},
		&models.UrgencyLevel{},
		&models.Gender{},
		&models.PriceSegment{},
		&models.City{},
		&models.District{},
	)
}
