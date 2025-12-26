package repository

import (
	"dental-marketplace/backend/internal/models"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRecordNotFound    = errors.New("record not found")
)

// Repository handles database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ==================== User Operations ====================

// AuthenticateUser validates user credentials and returns user with profile
func (r *Repository) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	
	// Find user by username with related profile data
	err := r.db.Preload("Patient").Preload("Clinic").Preload("Regulator").
		Where("username = ? AND is_active = ?", username, true).
		First(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// GetUserByID retrieves user by ID with profile
func (r *Repository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Patient").Preload("Clinic").Preload("Regulator").
		First(&user, userID).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(user *models.User) error {
	// Check if username already exists
	var count int64
	r.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return ErrUserAlreadyExists
	}

	return r.db.Create(user).Error
}

// ==================== Patient Operations ====================

// GetPatientByUserID retrieves patient profile by user ID
func (r *Repository) GetPatientByUserID(userID uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Where("user_id = ?", userID).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &patient, nil
}

// UpdatePatientSearchCriteria updates patient's clinic search preferences
func (r *Repository) UpdatePatientSearchCriteria(patientID uint, city, district, priceSegment string) error {
	return r.db.Model(&models.Patient{}).Where("id = ?", patientID).
		Updates(map[string]interface{}{
			"city":          city,
			"district":      district,
			"price_segment": priceSegment,
		}).Error
}

// GetPatientCTScans retrieves all CT scans for a patient
func (r *Repository) GetPatientCTScans(patientID uint) ([]models.CTScan, error) {
	var scans []models.CTScan
	err := r.db.Where("patient_id = ?", patientID).
		Order("upload_date DESC").
		Find(&scans).Error
	return scans, err
}

// GetCTScanByID retrieves a CT scan by ID
func (r *Repository) GetCTScanByID(scanID uint) (*models.CTScan, error) {
	var scan models.CTScan
	err := r.db.First(&scan, scanID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &scan, nil
}

// CreateCTScan creates a new CT scan
func (r *Repository) CreateCTScan(scan *models.CTScan) error {
	return r.db.Create(scan).Error
}

// GetTreatmentPlanByScanID retrieves treatment plan for a CT scan
func (r *Repository) GetTreatmentPlanByScanID(scanID uint) (*models.TreatmentPlan, error) {
	var plan models.TreatmentPlan
	err := r.db.Preload("Items").Where("ct_scan_id = ?", scanID).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &plan, nil
}

// GetTreatmentPlanByID retrieves treatment plan by ID with items and offers
func (r *Repository) GetTreatmentPlanByID(planID uint) (*models.TreatmentPlan, error) {
	var plan models.TreatmentPlan
	err := r.db.Preload("Items").Preload("Offers.Clinic").
		First(&plan, planID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &plan, nil
}

// GetPatientTreatmentPlans retrieves all treatment plans for a patient
func (r *Repository) GetPatientTreatmentPlans(patientID uint) ([]models.TreatmentPlan, error) {
	var plans []models.TreatmentPlan
	err := r.db.Preload("Items").Preload("Offers.Clinic").
		Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

// GetOffersForTreatmentPlan retrieves all clinic offers for a treatment plan
func (r *Repository) GetOffersForTreatmentPlan(planID uint) ([]models.ClinicOffer, error) {
	var offers []models.ClinicOffer
	err := r.db.Preload("Clinic").
		Where("treatment_plan_id = ? AND status != ?", planID, models.OfferStatusPending).
		Order("total_cost ASC").
		Find(&offers).Error
	return offers, err
}

// UpdateTreatmentPlanStatus updates the status of a treatment plan
func (r *Repository) UpdateTreatmentPlanStatus(planID uint, status string) error {
	return r.db.Model(&models.TreatmentPlan{}).
		Where("id = ?", planID).
		Update("status", status).Error
}

// GetPatientAppointments retrieves all appointments for a patient
func (r *Repository) GetPatientAppointments(patientID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Preload("Clinic").
		Where("patient_id = ?", patientID).
		Order("appointment_date DESC").
		Find(&appointments).Error
	return appointments, err
}

// CreateAppointment creates a new appointment
func (r *Repository) CreateAppointment(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

// CreateReview creates a new review
func (r *Repository) CreateReview(review *models.Review) error {
	// Create review
	if err := r.db.Create(review).Error; err != nil {
		return err
	}

	// Update clinic rating and review count
	return r.updateClinicRating(review.ClinicID)
}

// CreateComplaint creates a new complaint
func (r *Repository) CreateComplaint(complaint *models.Complaint) error {
	return r.db.Create(complaint).Error
}

// ==================== Clinic Operations ====================

// GetClinicByUserID retrieves clinic profile by user ID
func (r *Repository) GetClinicByUserID(userID uint) (*models.Clinic, error) {
	var clinic models.Clinic
	err := r.db.Where("user_id = ?", userID).First(&clinic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &clinic, nil
}

// GetClinicByID retrieves clinic by ID
func (r *Repository) GetClinicByID(clinicID uint) (*models.Clinic, error) {
	var clinic models.Clinic
	err := r.db.First(&clinic, clinicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &clinic, nil
}

// GetClinics retrieves all clinics with optional filters
func (r *Repository) GetClinics(city, district, priceSegment string) ([]models.Clinic, error) {
	query := r.db.Model(&models.Clinic{})

	if city != "" {
		query = query.Where("city = ?", city)
	}
	if district != "" {
		query = query.Where("district = ?", district)
	}

	var clinics []models.Clinic
	err := query.Order("rating DESC").Find(&clinics).Error
	return clinics, err
}

// GetClinicPriceList retrieves price list for a clinic
func (r *Repository) GetClinicPriceList(clinicID uint, specialization string) ([]models.PriceList, error) {
	query := r.db.Where("clinic_id = ?", clinicID)
	
	if specialization != "" {
		query = query.Where("specialization = ?", specialization)
	}

	var priceList []models.PriceList
	err := query.Order("specialization, service_name").Find(&priceList).Error
	return priceList, err
}

// UpdatePriceList updates clinic price list
func (r *Repository) UpdatePriceList(items []models.PriceList) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item.ID == 0 {
				// Create new
				if err := tx.Create(&item).Error; err != nil {
					return err
				}
			} else {
				// Update existing
				if err := tx.Save(&item).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// DeletePriceListItem deletes a price list item
func (r *Repository) DeletePriceListItem(itemID uint) error {
	return r.db.Delete(&models.PriceList{}, itemID).Error
}

// GetIncomingTreatmentPlans retrieves treatment plans for clinic to review
func (r *Repository) GetIncomingTreatmentPlans(clinicID uint, status string) ([]models.TreatmentPlan, error) {
	// Get treatment plans that match clinic capabilities and don't have an offer from this clinic yet
	var plans []models.TreatmentPlan
	
	query := r.db.Preload("Items").Preload("Offers", "clinic_id = ?", clinicID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Where("id NOT IN (?)", 
		r.db.Table("clinic_offers").
			Select("treatment_plan_id").
			Where("clinic_id = ?", clinicID),
	).Order("created_at DESC").Find(&plans).Error

	return plans, err
}

// CreateClinicOffer creates a new clinic offer
func (r *Repository) CreateClinicOffer(offer *models.ClinicOffer) error {
	return r.db.Create(offer).Error
}

// UpdateClinicOffer updates an existing clinic offer
func (r *Repository) UpdateClinicOffer(offer *models.ClinicOffer) error {
	return r.db.Save(offer).Error
}

// GetClinicOffers retrieves all offers made by a clinic
func (r *Repository) GetClinicOffers(clinicID uint, status string) ([]models.ClinicOffer, error) {
	query := r.db.Preload("TreatmentPlan").Where("clinic_id = ?", clinicID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var offers []models.ClinicOffer
	err := query.Order("created_at DESC").Find(&offers).Error
	return offers, err
}

// GetClinicLeads retrieves accepted offers (leads) for a clinic
func (r *Repository) GetClinicLeads(clinicID uint) ([]models.ClinicOffer, error) {
	var offers []models.ClinicOffer
	err := r.db.Preload("TreatmentPlan.Patient").
		Where("clinic_id = ? AND status = ?", clinicID, models.OfferStatusAccepted).
		Order("created_at DESC").
		Find(&offers).Error
	return offers, err
}

// GetClinicAppointments retrieves all appointments for a clinic
func (r *Repository) GetClinicAppointments(clinicID uint, status string) ([]models.Appointment, error) {
	query := r.db.Preload("Patient").Where("clinic_id = ?", clinicID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var appointments []models.Appointment
	err := query.Order("appointment_date DESC").Find(&appointments).Error
	return appointments, err
}

// UpdateAppointmentStatus updates appointment status
func (r *Repository) UpdateAppointmentStatus(appointmentID uint, status, notes string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if notes != "" {
		updates["notes"] = notes
	}
	
	return r.db.Model(&models.Appointment{}).
		Where("id = ?", appointmentID).
		Updates(updates).Error
}

// AcceptClinicOffer marks an offer as accepted and creates appointment
func (r *Repository) AcceptClinicOffer(offerID, patientID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get the offer
		var offer models.ClinicOffer
		if err := tx.First(&offer, offerID).Error; err != nil {
			return err
		}

		// Update offer status
		if err := tx.Model(&offer).Update("status", models.OfferStatusAccepted).Error; err != nil {
			return err
		}

		// Reject other offers for the same treatment plan
		if err := tx.Model(&models.ClinicOffer{}).
			Where("treatment_plan_id = ? AND id != ?", offer.TreatmentPlanID, offerID).
			Update("status", models.OfferStatusRejected).Error; err != nil {
			return err
		}

		// Update treatment plan status
		if err := tx.Model(&models.TreatmentPlan{}).
			Where("id = ?", offer.TreatmentPlanID).
			Update("status", models.PlanStatusOfferSelected).Error; err != nil {
			return err
		}

		// Create appointment
		appointment := &models.Appointment{
			PatientID:       patientID,
			ClinicID:        offer.ClinicID,
			TreatmentPlanID: offer.TreatmentPlanID,
			ClinicOfferID:   offerID,
			AppointmentDate: time.Now().Add(7 * 24 * time.Hour), // Default to 1 week from now
			Status:          models.AppointmentStatusScheduled,
			Notes:           "Initial consultation",
		}
		
		return tx.Create(appointment).Error
	})
}

// GetClinicDashboardMetrics retrieves key metrics for clinic dashboard
func (r *Repository) GetClinicDashboardMetrics(clinicID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Count new treatment plans
	var newPlansCount int64
	r.db.Model(&models.TreatmentPlan{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&newPlansCount)
	metrics["new_plans"] = newPlansCount

	// Count offers sent
	var offersSentCount int64
	r.db.Model(&models.ClinicOffer{}).
		Where("clinic_id = ? AND created_at BETWEEN ? AND ?", clinicID, startDate, endDate).
		Count(&offersSentCount)
	metrics["offers_sent"] = offersSentCount

	// Count accepted offers (leads)
	var leadsCount int64
	r.db.Model(&models.ClinicOffer{}).
		Where("clinic_id = ? AND status = ?", clinicID, models.OfferStatusAccepted).
		Count(&leadsCount)
	metrics["leads"] = leadsCount

	// Calculate potential revenue
	var totalRevenue int64
	r.db.Model(&models.ClinicOffer{}).
		Where("clinic_id = ? AND status = ?", clinicID, models.OfferStatusAccepted).
		Select("COALESCE(SUM(total_cost), 0)").
		Scan(&totalRevenue)
	metrics["potential_revenue"] = totalRevenue

	// Calculate conversion rate
	if offersSentCount > 0 {
		conversionRate := float64(leadsCount) / float64(offersSentCount) * 100
		metrics["conversion_rate"] = fmt.Sprintf("%.1f%%", conversionRate)
	} else {
		metrics["conversion_rate"] = "0%"
	}

	return metrics, nil
}

// ==================== Regulator Operations ====================

// GetRegulatorByUserID retrieves regulator profile by user ID
func (r *Repository) GetRegulatorByUserID(userID uint) (*models.Regulator, error) {
	var regulator models.Regulator
	err := r.db.Where("user_id = ?", userID).First(&regulator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &regulator, nil
}

// GetStatistics retrieves statistics for a date range
func (r *Repository) GetStatistics(startDate, endDate time.Time, clinicID *uint) ([]models.Statistics, error) {
	query := r.db.Where("date BETWEEN ? AND ?", startDate, endDate)
	
	if clinicID != nil {
		query = query.Where("clinic_id = ?", *clinicID)
	} else {
		query = query.Where("clinic_id IS NULL") // Regional aggregates
	}

	var stats []models.Statistics
	err := query.Order("date ASC").Find(&stats).Error
	return stats, err
}

// GetAllClinicsWithStats retrieves all clinics with their statistics
func (r *Repository) GetAllClinicsWithStats() ([]map[string]interface{}, error) {
	var clinics []models.Clinic
	err := r.db.Find(&clinics).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(clinics))
	for i, clinic := range clinics {
		// Get recent stats for this clinic
		var stats models.Statistics
		r.db.Where("clinic_id = ?", clinic.ID).
			Order("date DESC").
			First(&stats)

		result[i] = map[string]interface{}{
			"id":                clinic.ID,
			"name":              clinic.Name,
			"license_number":    clinic.LicenseNumber,
			"rating":            clinic.Rating,
			"review_count":      clinic.ReviewCount,
			"city":              clinic.City,
			"district":          clinic.District,
			"year_established":  clinic.YearEstablished,
			"patient_count":     stats.PatientCount,
			"total_revenue":     stats.TotalRevenue,
			"average_wait_days": stats.AverageWaitDays,
		}
	}

	return result, nil
}

// GetAllComplaints retrieves all complaints with optional filters
func (r *Repository) GetAllComplaints(status string) ([]models.Complaint, error) {
	query := r.db.Preload("Patient").Preload("Clinic")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var complaints []models.Complaint
	err := query.Order("created_at DESC").Find(&complaints).Error
	return complaints, err
}

// ==================== Helper Functions ====================

// updateClinicRating recalculates and updates clinic rating
func (r *Repository) updateClinicRating(clinicID uint) error {
	var avgRating float64
	var reviewCount int64

	r.db.Model(&models.Review{}).
		Where("clinic_id = ?", clinicID).
		Count(&reviewCount)

	r.db.Model(&models.Review{}).
		Where("clinic_id = ?", clinicID).
		Select("AVG(rating)").
		Scan(&avgRating)

	return r.db.Model(&models.Clinic{}).
		Where("id = ?", clinicID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": reviewCount,
		}).Error
}
