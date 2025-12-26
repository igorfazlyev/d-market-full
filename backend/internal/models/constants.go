package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Specialization represents dental specializations
type Specialization struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TreatmentStatus represents treatment plan statuses
type TreatmentStatus struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// OfferStatus represents offer statuses
type OfferStatus struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AppointmentStatus represents appointment statuses
type AppointmentStatus struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ScanStatus represents CT scan statuses
type ScanStatus struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UrgencyLevel represents urgency levels
type UrgencyLevel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Gender represents gender options
type Gender struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PriceSegment represents price segments
type PriceSegment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// City represents cities
type City struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// District represents districts within cities
type District struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CityID    uint           `gorm:"not null" json:"city_id"`
	City      City           `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Name      string         `gorm:"not null" json:"name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
