package repository

import (
	"dental-marketplace/backend/internal/models"

	"gorm.io/gorm"
)

type ConstantsRepository struct {
	db *gorm.DB
}

func NewConstantsRepository(db *gorm.DB) *ConstantsRepository {
	return &ConstantsRepository{db: db}
}

func (r *ConstantsRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&roles).Error
	return roles, err
}

func (r *ConstantsRepository) GetSpecializations() ([]models.Specialization, error) {
	var specs []models.Specialization
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&specs).Error
	return specs, err
}

func (r *ConstantsRepository) GetTreatmentStatuses() ([]models.TreatmentStatus, error) {
	var statuses []models.TreatmentStatus
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&statuses).Error
	return statuses, err
}

func (r *ConstantsRepository) GetOfferStatuses() ([]models.OfferStatus, error) {
	var statuses []models.OfferStatus
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&statuses).Error
	return statuses, err
}

func (r *ConstantsRepository) GetAppointmentStatuses() ([]models.AppointmentStatus, error) {
	var statuses []models.AppointmentStatus
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&statuses).Error
	return statuses, err
}

func (r *ConstantsRepository) GetScanStatuses() ([]models.ScanStatus, error) {
	var statuses []models.ScanStatus
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&statuses).Error
	return statuses, err
}

func (r *ConstantsRepository) GetUrgencyLevels() ([]models.UrgencyLevel, error) {
	var levels []models.UrgencyLevel
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&levels).Error
	return levels, err
}

func (r *ConstantsRepository) GetGenders() ([]models.Gender, error) {
	var genders []models.Gender
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&genders).Error
	return genders, err
}

func (r *ConstantsRepository) GetPriceSegments() ([]models.PriceSegment, error) {
	var segments []models.PriceSegment
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&segments).Error
	return segments, err
}

func (r *ConstantsRepository) GetCities() ([]models.City, error) {
	var cities []models.City
	err := r.db.Where("is_active = ?", true).Order("sort_order, name").Find(&cities).Error
	return cities, err
}

func (r *ConstantsRepository) GetDistricts() ([]models.District, error) {
	var districts []models.District
	err := r.db.Preload("City").Where("is_active = ?", true).Order("sort_order, name").Find(&districts).Error
	return districts, err
}

func (r *ConstantsRepository) GetDistrictsByCity(cityID uint) ([]models.District, error) {
	var districts []models.District
	err := r.db.Where("city_id = ? AND is_active = ?", cityID, true).Order("sort_order, name").Find(&districts).Error
	return districts, err
}
