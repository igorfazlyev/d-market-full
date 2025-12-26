package models

import (
	"time"

	"gorm.io/gorm"
)

// User roles
const (
	RolePatient   = "patient"
	RoleClinic    = "clinic"
	RoleRegulator = "regulator"
)

// Specializations
const (
	SpecTherapy      = "therapy"       // Терапия
	SpecOrthopedics  = "orthopedics"   // Ортопедия
	SpecSurgery      = "surgery"       // Хирургия
	SpecHygiene      = "hygiene"       // Гигиена
	SpecPeriodontics = "periodontics"  // Пародонтология
)

// Scan statuses
const (
	ScanStatusUploaded   = "uploaded"
	ScanStatusProcessing = "processing"
	ScanStatusCompleted  = "completed"
	ScanStatusError      = "error"
)

// Treatment plan statuses
const (
	PlanStatusGenerated = "generated"
	PlanStatusOffersRequested = "offers_requested"
	PlanStatusOffersReceived = "offers_received"
	PlanStatusOfferSelected = "offer_selected"
)

// Offer statuses
const (
	OfferStatusPending = "pending"
	OfferStatusSent = "sent"
	OfferStatusAccepted = "accepted"
	OfferStatusRejected = "rejected"
)

// Appointment statuses
const (
	AppointmentStatusScheduled = "scheduled"
	AppointmentStatusConfirmed = "confirmed"
	AppointmentStatusCompleted = "completed"
	AppointmentStatusCancelled = "cancelled"
	AppointmentStatusNoShow = "no_show"
)

// Lead statuses
const (
	LeadStatusNew = "new"
	LeadStatusContacted = "contacted"
	LeadStatusConsultationScheduled = "consultation_scheduled"
	LeadStatusTreatmentStarted = "treatment_started"
	LeadStatusTreatmentCompleted = "treatment_completed"
	LeadStatusRejected = "rejected"
)

// User represents all system users (patients, clinics, regulators)
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	Role         string `gorm:"not null;index" json:"role"` // patient, clinic, regulator
	Email        string `gorm:"uniqueIndex" json:"email"`
	Phone        string `json:"phone"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Patient   *Patient   `gorm:"foreignKey:UserID" json:"patient,omitempty"`
	Clinic    *Clinic    `gorm:"foreignKey:UserID" json:"clinic,omitempty"`
	Regulator *Regulator `gorm:"foreignKey:UserID" json:"regulator,omitempty"`
}

// Patient profile
type Patient struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID     uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender     string    `json:"gender"` // male, female, other
	
	// Search criteria
	City         string `json:"city"`
	District     string `json:"district"`
	PriceSegment string `json:"price_segment"` // economy, medium, premium
	
	// Relationships
	CTScans        []CTScan        `gorm:"foreignKey:PatientID" json:"ct_scans,omitempty"`
	TreatmentPlans []TreatmentPlan `gorm:"foreignKey:PatientID" json:"treatment_plans,omitempty"`
	Appointments   []Appointment   `gorm:"foreignKey:PatientID" json:"appointments,omitempty"`
	Reviews        []Review        `gorm:"foreignKey:PatientID" json:"reviews,omitempty"`
}

// Clinic profile
type Clinic struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID          uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	Name            string  `gorm:"not null" json:"name"`
	LicenseNumber   string  `gorm:"uniqueIndex" json:"license_number"`
	YearEstablished int     `json:"year_established"`
	Rating          float64 `gorm:"default:0" json:"rating"`
	ReviewCount     int     `gorm:"default:0" json:"review_count"`
	
	// Location
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
	
	// Capabilities
	HasTherapy      bool `gorm:"default:true" json:"has_therapy"`
	HasOrthopedics  bool `gorm:"default:true" json:"has_orthopedics"`
	HasSurgery      bool `gorm:"default:true" json:"has_surgery"`
	HasHygiene      bool `gorm:"default:true" json:"has_hygiene"`
	HasPeriodontics bool `gorm:"default:true" json:"has_periodontics"`
	
	// Services
	OffersInstallment bool `gorm:"default:false" json:"offers_installment"`
	OffersInsurance   bool `gorm:"default:false" json:"offers_insurance"`
	
	// Relationships
	PriceLists   []PriceList   `gorm:"foreignKey:ClinicID" json:"price_lists,omitempty"`
	Offers       []ClinicOffer `gorm:"foreignKey:ClinicID" json:"offers,omitempty"`
	Appointments []Appointment `gorm:"foreignKey:ClinicID" json:"appointments,omitempty"`
	Reviews      []Review      `gorm:"foreignKey:ClinicID" json:"reviews,omitempty"`
}

// Regulator profile
type Regulator struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID       uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Organization string `json:"organization"`
	Region       string `json:"region"`
	Position     string `json:"position"`
}

// CTScan represents uploaded CT scans
type CTScan struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	PatientID   uint      `gorm:"not null;index" json:"patient_id"`
	UploadDate  time.Time `json:"upload_date"`
	FileURL     string    `json:"file_url"`
	Status      string    `gorm:"default:'uploaded'" json:"status"` // uploaded, processing, completed, error
	AIProcessed bool      `gorm:"default:false" json:"ai_processed"`
	
	// Relationships
	TreatmentPlan *TreatmentPlan `gorm:"foreignKey:CTScanID" json:"treatment_plan,omitempty"`
}

// TreatmentPlan generated by AI
type TreatmentPlan struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	PatientID uint   `gorm:"not null;index" json:"patient_id"`
	CTScanID  uint   `gorm:"uniqueIndex" json:"ct_scan_id"`
	Status    string `gorm:"default:'generated'" json:"status"`
	
	// Summary by specialization
	RequiresTherapy      bool `json:"requires_therapy"`
	RequiresOrthopedics  bool `json:"requires_orthopedics"`
	RequiresSurgery      bool `json:"requires_surgery"`
	RequiresHygiene      bool `json:"requires_hygiene"`
	RequiresPeriodontics bool `json:"requires_periodontics"`
	
	// Estimated costs (ranges)
	TherapyMinCost      int `json:"therapy_min_cost"`
	TherapyMaxCost      int `json:"therapy_max_cost"`
	OrthopedicsMinCost  int `json:"orthopedics_min_cost"`
	OrthopedicsMaxCost  int `json:"orthopedics_max_cost"`
	SurgeryMinCost      int `json:"surgery_min_cost"`
	SurgeryMaxCost      int `json:"surgery_max_cost"`
	HygieneMinCost      int `json:"hygiene_min_cost"`
	HygieneMaxCost      int `json:"hygiene_max_cost"`
	PeriodonticsMinCost int `json:"periodontics_min_cost"`
	PeriodonticsMaxCost int `json:"periodontics_max_cost"`
	
	// Relationships
	Items  []TreatmentItem `gorm:"foreignKey:TreatmentPlanID" json:"items,omitempty"`
	Offers []ClinicOffer   `gorm:"foreignKey:TreatmentPlanID" json:"offers,omitempty"`
}

// TreatmentItem individual procedure in treatment plan
type TreatmentItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	TreatmentPlanID uint   `gorm:"not null;index" json:"treatment_plan_id"`
	Specialization  string `gorm:"not null" json:"specialization"` // therapy, orthopedics, surgery, hygiene, periodontics
	ToothNumber     string `json:"tooth_number"` // International notation: 11-48
	Diagnosis       string `json:"diagnosis"`
	Procedure       string `json:"procedure"`
	Urgency         string `json:"urgency"` // high, medium, low
	EstimatedCost   int    `json:"estimated_cost"`
}

// PriceList clinic pricing
type PriceList struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	ClinicID       uint   `gorm:"not null;index" json:"clinic_id"`
	Specialization string `gorm:"not null" json:"specialization"`
	ServiceName    string `gorm:"not null" json:"service_name"`
	Price          int    `gorm:"not null" json:"price"`
	WarrantyYears  int    `json:"warranty_years"`
}

// ClinicOffer from clinic to patient
type ClinicOffer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	TreatmentPlanID uint   `gorm:"not null;index" json:"treatment_plan_id"`
	ClinicID        uint   `gorm:"not null;index" json:"clinic_id"`
	Status          string `gorm:"default:'pending'" json:"status"` // pending, sent, accepted, rejected
	
	// Costs by specialization
	TherapyCost      int `json:"therapy_cost"`
	OrthopedicsCost  int `json:"orthopedics_cost"`
	SurgeryCost      int `json:"surgery_cost"`
	HygieneCost      int `json:"hygiene_cost"`
	PeriodonticsCost int `json:"periodontics_cost"`
	TotalCost        int `json:"total_cost"`
	
	// Terms
	EstimatedDuration string `json:"estimated_duration"` // e.g., "2-3 months"
	InstallmentMonths int    `json:"installment_months"`
	WarrantyDetails   string `json:"warranty_details"`
	Notes             string `json:"notes"`
	
	// Relationships
	Clinic Clinic `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
}

// Appointment between patient and clinic
type Appointment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	PatientID       uint      `gorm:"not null;index" json:"patient_id"`
	ClinicID        uint      `gorm:"not null;index" json:"clinic_id"`
	TreatmentPlanID uint      `json:"treatment_plan_id"`
	ClinicOfferID   uint      `json:"clinic_offer_id"`
	
	AppointmentDate time.Time `json:"appointment_date"`
	Specialization  string    `json:"specialization"`
	Status          string    `gorm:"default:'scheduled'" json:"status"`
	Notes           string    `json:"notes"`
	
	// Relationships
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Clinic  Clinic  `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
}

// Review patient feedback
type Review struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	PatientID uint   `gorm:"not null;index" json:"patient_id"`
	ClinicID  uint   `gorm:"not null;index" json:"clinic_id"`
	Rating    int    `gorm:"not null" json:"rating"` // 1-5
	Comment   string `json:"comment"`
	IsPublic  bool   `gorm:"default:false" json:"is_public"`
	
	// Relationships
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Clinic  Clinic  `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
}

// Complaint tracking
type Complaint struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	PatientID   uint   `gorm:"not null;index" json:"patient_id"`
	ClinicID    uint   `gorm:"not null;index" json:"clinic_id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      string `gorm:"default:'open'" json:"status"` // open, in_progress, resolved, closed
	Resolution  string `json:"resolution"`
	
	// Relationships
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Clinic  Clinic  `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
}

// Statistics for regulator dashboard (time-series data)
type Statistics struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Date     time.Time `gorm:"index" json:"date"`
	ClinicID *uint     `gorm:"index" json:"clinic_id"` // null for regional aggregates
	
	// Metrics
	TreatmentPlansGenerated int `json:"treatment_plans_generated"`
	AppointmentsScheduled   int `json:"appointments_scheduled"`
	AppointmentsCompleted   int `json:"appointments_completed"`
	TotalRevenue            int `json:"total_revenue"`
	PatientCount            int `json:"patient_count"`
	
	// Disease statistics
	CariesCount       int `json:"caries_count"`
	PulpitisCount     int `json:"pulpitis_count"`
	PeriodontitisCount int `json:"periodontitis_count"`
	GingivitisCount   int `json:"gingivitis_count"`
	ParodontitisCount int `json:"parodontitis_count"`
	
	// Average metrics
	AverageWaitDays   float64 `json:"average_wait_days"`
	AverageTreatmentCost int  `json:"average_treatment_cost"`
}
