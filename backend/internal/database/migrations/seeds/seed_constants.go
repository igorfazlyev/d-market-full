package seeds

import (
	"dental-marketplace/backend/internal/models"
	"log"

	"gorm.io/gorm"
)

func SeedConstants(db *gorm.DB) error {
	log.Println("🌱 Seeding constants...")

	// Seed Roles
	roles := []models.Role{
		{Code: "patient", Name: "Пациент", SortOrder: 1},
		{Code: "clinic", Name: "Клиника", SortOrder: 2},
		{Code: "regulator", Name: "Регулятор", SortOrder: 3},
	}
	for _, role := range roles {
		db.Where(models.Role{Code: role.Code}).FirstOrCreate(&role)
	}

	// Seed Specializations
	specializations := []models.Specialization{
		{Code: "therapy", Name: "Терапия", SortOrder: 1},
		{Code: "orthopedics", Name: "Ортопедия", SortOrder: 2},
		{Code: "surgery", Name: "Хирургия", SortOrder: 3},
		{Code: "hygiene", Name: "Гигиена", SortOrder: 4},
		{Code: "periodontics", Name: "Пародонтология", SortOrder: 5},
	}
	for _, spec := range specializations {
		db.Where(models.Specialization{Code: spec.Code}).FirstOrCreate(&spec)
	}

	// Seed Treatment Statuses
	treatmentStatuses := []models.TreatmentStatus{
		{Code: "generated", Name: "Сформирован", SortOrder: 1},
		{Code: "offers_received", Name: "Получены предложения", SortOrder: 2},
		{Code: "offer_accepted", Name: "Предложение принято", SortOrder: 3},
		{Code: "in_progress", Name: "В процессе", SortOrder: 4},
		{Code: "completed", Name: "Завершен", SortOrder: 5},
	}
	for _, status := range treatmentStatuses {
		db.Where(models.TreatmentStatus{Code: status.Code}).FirstOrCreate(&status)
	}

	// Seed Offer Statuses
	offerStatuses := []models.OfferStatus{
		{Code: "pending", Name: "В ожидании", SortOrder: 1},
		{Code: "sent", Name: "Отправлено", SortOrder: 2},
		{Code: "accepted", Name: "Принято", SortOrder: 3},
		{Code: "rejected", Name: "Отклонено", SortOrder: 4},
	}
	for _, status := range offerStatuses {
		db.Where(models.OfferStatus{Code: status.Code}).FirstOrCreate(&status)
	}

	// Seed Appointment Statuses
	appointmentStatuses := []models.AppointmentStatus{
		{Code: "scheduled", Name: "Запланирован", SortOrder: 1},
		{Code: "confirmed", Name: "Подтвержден", SortOrder: 2},
		{Code: "completed", Name: "Завершен", SortOrder: 3},
		{Code: "cancelled", Name: "Отменен", SortOrder: 4},
	}
	for _, status := range appointmentStatuses {
		db.Where(models.AppointmentStatus{Code: status.Code}).FirstOrCreate(&status)
	}

	// Seed Scan Statuses
	scanStatuses := []models.ScanStatus{
		{Code: "uploaded", Name: "Загружен", SortOrder: 1},
		{Code: "processing", Name: "Обрабатывается", SortOrder: 2},
		{Code: "completed", Name: "Завершен", SortOrder: 3},
		{Code: "failed", Name: "Ошибка", SortOrder: 4},
	}
	for _, status := range scanStatuses {
		db.Where(models.ScanStatus{Code: status.Code}).FirstOrCreate(&status)
	}

	// Seed Urgency Levels
	urgencyLevels := []models.UrgencyLevel{
		{Code: "low", Name: "Низкая", SortOrder: 1},
		{Code: "medium", Name: "Средняя", SortOrder: 2},
		{Code: "high", Name: "Высокая", SortOrder: 3},
	}
	for _, level := range urgencyLevels {
		db.Where(models.UrgencyLevel{Code: level.Code}).FirstOrCreate(&level)
	}

	// Seed Genders
	genders := []models.Gender{
		{Code: "male", Name: "Мужской", SortOrder: 1},
		{Code: "female", Name: "Женский", SortOrder: 2},
		{Code: "other", Name: "Другой", SortOrder: 3},
	}
	for _, gender := range genders {
		db.Where(models.Gender{Code: gender.Code}).FirstOrCreate(&gender)
	}

	// Seed Price Segments
	priceSegments := []models.PriceSegment{
		{Code: "economy", Name: "эконом", SortOrder: 1},
		{Code: "medium", Name: "средний", SortOrder: 2},
		{Code: "premium", Name: "премиум", SortOrder: 3},
	}
	for _, segment := range priceSegments {
		db.Where(models.PriceSegment{Code: segment.Code}).FirstOrCreate(&segment)
	}

	// Seed Cities
	cities := []models.City{
		{Code: "moscow", Name: "Москва", SortOrder: 1},
		{Code: "spb", Name: "Санкт-Петербург", SortOrder: 2},
		{Code: "kazan", Name: "Казань", SortOrder: 3},
		{Code: "ekb", Name: "Екатеринбург", SortOrder: 4},
		{Code: "nsk", Name: "Новосибирск", SortOrder: 5},
	}
	for _, city := range cities {
		db.Where(models.City{Code: city.Code}).FirstOrCreate(&city)
	}

	// Seed Districts for Moscow
	var moscow models.City
	db.Where("code = ?", "moscow").First(&moscow)

	moscowDistricts := []models.District{
		{CityID: moscow.ID, Code: "moscow_central", Name: "Центральный", SortOrder: 1},
		{CityID: moscow.ID, Code: "moscow_northern", Name: "Северный", SortOrder: 2},
		{CityID: moscow.ID, Code: "moscow_southern", Name: "Южный", SortOrder: 3},
		{CityID: moscow.ID, Code: "moscow_eastern", Name: "Восточный", SortOrder: 4},
		{CityID: moscow.ID, Code: "moscow_western", Name: "Западный", SortOrder: 5},
	}
	for _, district := range moscowDistricts {
		db.Where(models.District{Code: district.Code}).FirstOrCreate(&district)
	}

	// Seed Districts for SPb
	var spb models.City
	db.Where("code = ?", "spb").First(&spb)

	spbDistricts := []models.District{
		{CityID: spb.ID, Code: "spb_central", Name: "Центральный", SortOrder: 1},
		{CityID: spb.ID, Code: "spb_nevsky", Name: "Невский", SortOrder: 2},
		{CityID: spb.ID, Code: "spb_vasileostrovsky", Name: "Василеостровский", SortOrder: 3},
		{CityID: spb.ID, Code: "spb_admiralteysky", Name: "Адмиралтейский", SortOrder: 4},
	}
	for _, district := range spbDistricts {
		db.Where(models.District{Code: district.Code}).FirstOrCreate(&district)
	}

	log.Println("✅ Constants seeded successfully")
	return nil
}
