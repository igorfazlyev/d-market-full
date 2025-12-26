package migrations

import (
	"dental-marketplace/backend/internal/models"

	"gorm.io/gorm"
)

func CreateBusinessTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Patient{},
		&models.Clinic{},
		&models.Regulator{},
		&models.CTScan{},
		&models.TreatmentPlan{},
		&models.TreatmentItem{},
		&models.PriceList{},
		&models.ClinicOffer{},
		&models.Appointment{},
		&models.Review{},
		&models.Complaint{},
		&models.Statistics{},
	)
}
