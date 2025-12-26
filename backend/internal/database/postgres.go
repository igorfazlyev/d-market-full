package database

import (
	"dental-marketplace/backend/internal/models"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// Connect establishes database connection
func Connect(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established")

	return &Database{DB: db}, nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	err := d.DB.AutoMigrate(
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

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// Seed populates database with sample data
// Seed populates database with sample data
func (d *Database) Seed() error {
	log.Println("🌱 Seeding database with sample data...")

	// Check if data already exists
	var userCount int64
	d.DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("⏭️  Database already seeded, skipping...")
		return nil
	}

	// Hash password for demo accounts
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 1. CREATE PATIENT
	patientUser := &models.User{
		Username:     "patient",
		PasswordHash: string(hashedPassword),
		Role:         models.RolePatient,
		Email:        "anna.petrova@example.com",
		Phone:        "+7 916 555-1234",
		IsActive:     true,
	}
	if err := d.DB.Create(patientUser).Error; err != nil {
		return fmt.Errorf("failed to create patient user: %w", err)
	}

	patient := &models.Patient{
		UserID:       patientUser.ID,
		FirstName:    "Анна",
		LastName:     "Петрова",
		DateOfBirth:  time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
		Gender:       "female",
		City:         "Москва",
		District:     "Центральный",
		PriceSegment: "средний",
	}
	if err := d.DB.Create(patient).Error; err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}

	// 2. CREATE CLINICS
	// Clinic 1: StomaPro
	clinic1User := &models.User{
		Username:     "clinic1",
		PasswordHash: string(hashedPassword),
		Role:         models.RoleClinic,
		Email:        "info@stomapro.ru",
		Phone:        "+7 495 123-4567",
		IsActive:     true,
	}
	if err := d.DB.Create(clinic1User).Error; err != nil {
		return fmt.Errorf("failed to create clinic1 user: %w", err)
	}

	clinic1 := &models.Clinic{
		UserID:            clinic1User.ID,
		Name:              "СтомаПро",
		LicenseNumber:     "ЛО-77-01-012345",
		YearEstablished:   2015,
		Rating:            4.8,
		ReviewCount:       156,
		City:              "Москва",
		District:          "Центральный",
		Address:           "ул. Тверская, д. 15",
		HasTherapy:        true,
		HasOrthopedics:    true,
		HasSurgery:        true,
		HasHygiene:        true,
		HasPeriodontics:   true,
		OffersInstallment: true,
		OffersInsurance:   false,
	}
	if err := d.DB.Create(clinic1).Error; err != nil {
		return fmt.Errorf("failed to create clinic1: %w", err)
	}

	// Clinic 2: DentalPlus
	clinic2User := &models.User{
		Username:     "clinic2",
		PasswordHash: string(hashedPassword),
		Role:         models.RoleClinic,
		Email:        "contact@dentalplus.ru",
		Phone:        "+7 495 987-6543",
		IsActive:     true,
	}
	if err := d.DB.Create(clinic2User).Error; err != nil {
		return fmt.Errorf("failed to create clinic2 user: %w", err)
	}

	clinic2 := &models.Clinic{
		UserID:            clinic2User.ID,
		Name:              "ДентаПлюс",
		LicenseNumber:     "ЛО-77-01-067890",
		YearEstablished:   2018,
		Rating:            4.5,
		ReviewCount:       98,
		City:              "Москва",
		District:          "Северный",
		Address:           "Дмитровское шоссе, д. 89",
		HasTherapy:        true,
		HasOrthopedics:    true,
		HasSurgery:        true,
		HasHygiene:        true,
		HasPeriodontics:   false,
		OffersInstallment: true,
		OffersInsurance:   true,
	}
	if err := d.DB.Create(clinic2).Error; err != nil {
		return fmt.Errorf("failed to create clinic2: %w", err)
	}

	// 3. CREATE REGULATOR
	regulatorUser := &models.User{
		Username:     "regulator",
		PasswordHash: string(hashedPassword),
		Role:         models.RoleRegulator,
		Email:        "regulator@health.gov.ru",
		Phone:        "+7 495 777-8888",
		IsActive:     true,
	}
	if err := d.DB.Create(regulatorUser).Error; err != nil {
		return fmt.Errorf("failed to create regulator user: %w", err)
	}

	regulator := &models.Regulator{
		UserID:       regulatorUser.ID,
		Organization: "Департамент здравоохранения города Москвы",
		Region:       "Москва",
		Position:     "Старший инспектор",
	}
	if err := d.DB.Create(regulator).Error; err != nil {
		return fmt.Errorf("failed to create regulator: %w", err)
	}

	// 4. CREATE PRICE LISTS FOR CLINICS
	priceListClinic1 := []models.PriceList{
		// Терапия
		{ClinicID: clinic1.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение кариеса", Price: 5000, WarrantyYears: 1},
		{ClinicID: clinic1.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение пульпита (лечение каналов)", Price: 15000, WarrantyYears: 2},
		{ClinicID: clinic1.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение периодонтита", Price: 12000, WarrantyYears: 1},
		{ClinicID: clinic1.ID, Specialization: models.SpecTherapy, ServiceName: "Пломба светоотверждаемая", Price: 4500, WarrantyYears: 1},

		// Ортопедия
		{ClinicID: clinic1.ID, Specialization: models.SpecOrthopedics, ServiceName: "Коронка металлокерамическая", Price: 30000, WarrantyYears: 3},
		{ClinicID: clinic1.ID, Specialization: models.SpecOrthopedics, ServiceName: "Коронка циркониевая", Price: 45000, WarrantyYears: 5},
		{ClinicID: clinic1.ID, Specialization: models.SpecOrthopedics, ServiceName: "Мостовидный протез (3 единицы)", Price: 85000, WarrantyYears: 3},

		// Хирургия
		{ClinicID: clinic1.ID, Specialization: models.SpecSurgery, ServiceName: "Удаление зуба (простое)", Price: 3000, WarrantyYears: 0},
		{ClinicID: clinic1.ID, Specialization: models.SpecSurgery, ServiceName: "Удаление зуба (сложное)", Price: 8000, WarrantyYears: 0},
		{ClinicID: clinic1.ID, Specialization: models.SpecSurgery, ServiceName: "Имплант (Nobel Biocare)", Price: 95000, WarrantyYears: 10},
		{ClinicID: clinic1.ID, Specialization: models.SpecSurgery, ServiceName: "Имплант (Straumann)", Price: 120000, WarrantyYears: 15},
		{ClinicID: clinic1.ID, Specialization: models.SpecSurgery, ServiceName: "Костная пластика", Price: 45000, WarrantyYears: 0},

		// Гигиена
		{ClinicID: clinic1.ID, Specialization: models.SpecHygiene, ServiceName: "Профессиональная чистка зубов", Price: 5000, WarrantyYears: 0},
		{ClinicID: clinic1.ID, Specialization: models.SpecHygiene, ServiceName: "Отбеливание зубов", Price: 18000, WarrantyYears: 0},
		{ClinicID: clinic1.ID, Specialization: models.SpecHygiene, ServiceName: "Чистка Air Flow", Price: 4000, WarrantyYears: 0},

		// Пародонтология
		{ClinicID: clinic1.ID, Specialization: models.SpecPeriodontics, ServiceName: "Лечение пародонтита (за квадрант)", Price: 15000, WarrantyYears: 1},
		{ClinicID: clinic1.ID, Specialization: models.SpecPeriodontics, ServiceName: "Пластика десны", Price: 35000, WarrantyYears: 2},
	}

	priceListClinic2 := []models.PriceList{
		// Терапия - немного дешевле
		{ClinicID: clinic2.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение кариеса", Price: 4000, WarrantyYears: 1},
		{ClinicID: clinic2.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение пульпита (лечение каналов)", Price: 12000, WarrantyYears: 2},
		{ClinicID: clinic2.ID, Specialization: models.SpecTherapy, ServiceName: "Лечение периодонтита", Price: 10000, WarrantyYears: 1},
		{ClinicID: clinic2.ID, Specialization: models.SpecTherapy, ServiceName: "Пломба светоотверждаемая", Price: 3500, WarrantyYears: 1},

		// Ортопедия
		{ClinicID: clinic2.ID, Specialization: models.SpecOrthopedics, ServiceName: "Коронка металлокерамическая", Price: 25000, WarrantyYears: 2},
		{ClinicID: clinic2.ID, Specialization: models.SpecOrthopedics, ServiceName: "Коронка циркониевая", Price: 38000, WarrantyYears: 4},
		{ClinicID: clinic2.ID, Specialization: models.SpecOrthopedics, ServiceName: "Мостовидный протез (3 единицы)", Price: 70000, WarrantyYears: 2},

		// Хирургия
		{ClinicID: clinic2.ID, Specialization: models.SpecSurgery, ServiceName: "Удаление зуба (простое)", Price: 2500, WarrantyYears: 0},
		{ClinicID: clinic2.ID, Specialization: models.SpecSurgery, ServiceName: "Удаление зуба (сложное)", Price: 7000, WarrantyYears: 0},
		{ClinicID: clinic2.ID, Specialization: models.SpecSurgery, ServiceName: "Имплант (Osstem)", Price: 75000, WarrantyYears: 7},
		{ClinicID: clinic2.ID, Specialization: models.SpecSurgery, ServiceName: "Имплант (Nobel Biocare)", Price: 85000, WarrantyYears: 10},
		{ClinicID: clinic2.ID, Specialization: models.SpecSurgery, ServiceName: "Костная пластика", Price: 38000, WarrantyYears: 0},

		// Гигиена
		{ClinicID: clinic2.ID, Specialization: models.SpecHygiene, ServiceName: "Профессиональная чистка зубов", Price: 4000, WarrantyYears: 0},
		{ClinicID: clinic2.ID, Specialization: models.SpecHygiene, ServiceName: "Отбеливание зубов", Price: 15000, WarrantyYears: 0},
		{ClinicID: clinic2.ID, Specialization: models.SpecHygiene, ServiceName: "Чистка Air Flow", Price: 3500, WarrantyYears: 0},
	}

	if err := d.DB.Create(&priceListClinic1).Error; err != nil {
		return fmt.Errorf("failed to create price list for clinic1: %w", err)
	}
	if err := d.DB.Create(&priceListClinic2).Error; err != nil {
		return fmt.Errorf("failed to create price list for clinic2: %w", err)
	}

	// 5. CREATE CT SCANS FOR PATIENT
	scan1 := &models.CTScan{
		PatientID:   patient.ID,
		UploadDate:  time.Date(2024, 11, 15, 10, 30, 0, 0, time.UTC),
		FileURL:     "/uploads/scans/scan_001_20241115.dcm",
		Status:      models.ScanStatusCompleted,
		AIProcessed: true,
	}
	if err := d.DB.Create(scan1).Error; err != nil {
		return fmt.Errorf("failed to create scan1: %w", err)
	}

	scan2 := &models.CTScan{
		PatientID:   patient.ID,
		UploadDate:  time.Date(2024, 12, 10, 14, 15, 0, 0, time.UTC),
		FileURL:     "/uploads/scans/scan_002_20241210.dcm",
		Status:      models.ScanStatusCompleted,
		AIProcessed: true,
	}
	if err := d.DB.Create(scan2).Error; err != nil {
		return fmt.Errorf("failed to create scan2: %w", err)
	}

	// 6. CREATE TREATMENT PLAN FOR SCAN 2 (most recent)
	treatmentPlan := &models.TreatmentPlan{
		PatientID:            patient.ID,
		CTScanID:             scan2.ID,
		Status:               models.PlanStatusGenerated,
		RequiresTherapy:      true,
		RequiresOrthopedics:  true,
		RequiresSurgery:      true,
		RequiresHygiene:      true,
		RequiresPeriodontics: false,
		TherapyMinCost:       25000,
		TherapyMaxCost:       35000,
		OrthopedicsMinCost:   55000,
		OrthopedicsMaxCost:   75000,
		SurgeryMinCost:       95000,
		SurgeryMaxCost:       120000,
		HygieneMinCost:       4000,
		HygieneMaxCost:       6000,
		PeriodonticsMinCost:  0,
		PeriodonticsMaxCost:  0,
	}
	if err := d.DB.Create(treatmentPlan).Error; err != nil {
		return fmt.Errorf("failed to create treatment plan: %w", err)
	}

	// 7. CREATE TREATMENT ITEMS
	treatmentItems := []models.TreatmentItem{
		// Терапия
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecTherapy, ToothNumber: "16", Diagnosis: "Глубокий кариес", Procedure: "Лечение кариеса + пломба", Urgency: "high", EstimatedCost: 9500},
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecTherapy, ToothNumber: "25", Diagnosis: "Острый пульпит", Procedure: "Лечение каналов", Urgency: "high", EstimatedCost: 15000},
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecTherapy, ToothNumber: "14", Diagnosis: "Поверхностный кариес", Procedure: "Лечение кариеса + пломба", Urgency: "medium", EstimatedCost: 8000},

		// Ортопедия
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecOrthopedics, ToothNumber: "25", Diagnosis: "Восстановление после пульпита", Procedure: "Коронка циркониевая", Urgency: "medium", EstimatedCost: 45000},
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecOrthopedics, ToothNumber: "21", Diagnosis: "Скол коронки", Procedure: "Коронка циркониевая", Urgency: "high", EstimatedCost: 45000},

		// Хирургия
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecSurgery, ToothNumber: "37", Diagnosis: "Отсутствует зуб", Procedure: "Имплант (Nobel Biocare)", Urgency: "medium", EstimatedCost: 95000},

		// Гигиена
		{TreatmentPlanID: treatmentPlan.ID, Specialization: models.SpecHygiene, ToothNumber: "Все", Diagnosis: "Налет и зубной камень", Procedure: "Профессиональная чистка", Urgency: "medium", EstimatedCost: 5000},
	}
	if err := d.DB.Create(&treatmentItems).Error; err != nil {
		return fmt.Errorf("failed to create treatment items: %w", err)
	}

	// 8. CREATE CLINIC OFFERS
	offer1 := &models.ClinicOffer{
		TreatmentPlanID:   treatmentPlan.ID,
		ClinicID:          clinic1.ID,
		Status:            models.OfferStatusSent,
		TherapyCost:       32500,
		OrthopedicsCost:   90000,
		SurgeryCost:       95000,
		HygieneCost:       5000,
		PeriodonticsCost:  0,
		TotalCost:         222500,
		EstimatedDuration: "3-4 месяца",
		InstallmentMonths: 12,
		WarrantyDetails:   "10 лет на имплант, 5 лет на коронки, 1-2 года на пломбы",
		Notes:             "Премиальные материалы, опытные хирурги, индивидуальный подход",
	}
	if err := d.DB.Create(offer1).Error; err != nil {
		return fmt.Errorf("failed to create offer1: %w", err)
	}

	offer2 := &models.ClinicOffer{
		TreatmentPlanID:   treatmentPlan.ID,
		ClinicID:          clinic2.ID,
		Status:            models.OfferStatusSent,
		TherapyCost:       26500,
		OrthopedicsCost:   76000,
		SurgeryCost:       85000,
		HygieneCost:       4000,
		PeriodonticsCost:  0,
		TotalCost:         191500,
		EstimatedDuration: "2-3 месяца",
		InstallmentMonths: 12,
		WarrantyDetails:   "10 лет на имплант, 4 года на коронки, 1-2 года на пломбы",
		Notes:             "Хорошее соотношение цена-качество, принимаем страховки, гибкий график",
	}
	if err := d.DB.Create(offer2).Error; err != nil {
		return fmt.Errorf("failed to create offer2: %w", err)
	}

	// 9. CREATE STATISTICS FOR LAST 90 DAYS
	now := time.Now().UTC()
	for i := 90; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)

		// Regional aggregates (ClinicID = nil)
		regionalStats := &models.Statistics{
			Date:                    date,
			ClinicID:                nil,
			TreatmentPlansGenerated: 15 + (i % 10),
			AppointmentsScheduled:   12 + (i % 8),
			AppointmentsCompleted:   10 + (i % 7),
			TotalRevenue:            850000 + (i * 5000),
			PatientCount:            25 + (i % 15),
			CariesCount:             20 + (i % 5),
			PulpitisCount:           8 + (i % 3),
			PeriodontitisCount:      5 + (i % 2),
			GingivitisCount:         6 + (i % 3),
			ParodontitisCount:       4 + (i % 2),
			AverageWaitDays:         3.5 + float64(i%3),
			AverageTreatmentCost:    175000 + (i * 1000),
		}
		if err := d.DB.Create(regionalStats).Error; err != nil {
			return fmt.Errorf("failed to create regional stats: %w", err)
		}

		// Clinic 1 stats
		clinic1Stats := &models.Statistics{
			Date:                    date,
			ClinicID:                &clinic1.ID,
			TreatmentPlansGenerated: 8 + (i % 5),
			AppointmentsScheduled:   7 + (i % 4),
			AppointmentsCompleted:   6 + (i % 4),
			TotalRevenue:            450000 + (i * 3000),
			PatientCount:            14 + (i % 8),
			CariesCount:             11 + (i % 3),
			PulpitisCount:           4 + (i % 2),
			PeriodontitisCount:      3 + (i % 2),
			GingivitisCount:         3 + (i % 2),
			ParodontitisCount:       2 + (i % 1),
			AverageWaitDays:         2.5 + float64(i%2),
			AverageTreatmentCost:    195000 + (i * 800),
		}
		if err := d.DB.Create(clinic1Stats).Error; err != nil {
			return fmt.Errorf("failed to create clinic1 stats: %w", err)
		}

		// Clinic 2 stats
		clinic2Stats := &models.Statistics{
			Date:                    date,
			ClinicID:                &clinic2.ID,
			TreatmentPlansGenerated: 7 + (i % 5),
			AppointmentsScheduled:   5 + (i % 4),
			AppointmentsCompleted:   4 + (i % 3),
			TotalRevenue:            400000 + (i * 2000),
			PatientCount:            11 + (i % 7),
			CariesCount:             9 + (i % 2),
			PulpitisCount:           4 + (i % 1),
			PeriodontitisCount:      2 + (i % 1),
			GingivitisCount:         3 + (i % 1),
			ParodontitisCount:       2 + (i % 1),
			AverageWaitDays:         4.5 + float64(i%3),
			AverageTreatmentCost:    155000 + (i * 600),
		}
		if err := d.DB.Create(clinic2Stats).Error; err != nil {
			return fmt.Errorf("failed to create clinic2 stats: %w", err)
		}
	}

	log.Println("✅ База данных успешно заполнена!")
	log.Println("")
	log.Println("📋 Учетные данные для входа:")
	log.Println("   Пациент:   username: patient   | password: password")
	log.Println("   Клиника 1: username: clinic1   | password: password")
	log.Println("   Клиника 2: username: clinic2   | password: password")
	log.Println("   Регулятор: username: regulator | password: password")
	log.Println("")

	return nil
}

// Close closes database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
