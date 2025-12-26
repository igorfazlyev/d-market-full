package handlers

import (
	"dental-marketplace/backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

// DateRangeQuery represents date range parameters
type DateRangeQuery struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type CommonHandler struct {
	constantsRepo *repository.ConstantsRepository
}

func NewCommonHandler(constantsRepo *repository.ConstantsRepository) *CommonHandler {
	return &CommonHandler{
		constantsRepo: constantsRepo,
	}
}

// GetConstants returns all system constants from database
func (h *CommonHandler) GetConstants(c *gin.Context) {
	// Fetch all constants from database
	roles, _ := h.constantsRepo.GetRoles()
	specializations, _ := h.constantsRepo.GetSpecializations()
	treatmentStatuses, _ := h.constantsRepo.GetTreatmentStatuses()
	offerStatuses, _ := h.constantsRepo.GetOfferStatuses()
	appointmentStatuses, _ := h.constantsRepo.GetAppointmentStatuses()
	scanStatuses, _ := h.constantsRepo.GetScanStatuses()
	urgencyLevels, _ := h.constantsRepo.GetUrgencyLevels()
	genders, _ := h.constantsRepo.GetGenders()
	priceSegments, _ := h.constantsRepo.GetPriceSegments()
	cities, _ := h.constantsRepo.GetCities()
	districts, _ := h.constantsRepo.GetDistricts()

	// Convert to maps for easier frontend consumption
	rolesMap := make(map[string]string)
	for _, r := range roles {
		rolesMap[r.Code] = r.Name
	}

	specsMap := make(map[string]string)
	for _, s := range specializations {
		specsMap[s.Code] = s.Name
	}

	treatmentStatusesMap := make(map[string]string)
	for _, ts := range treatmentStatuses {
		treatmentStatusesMap[ts.Code] = ts.Name
	}

	offerStatusesMap := make(map[string]string)
	for _, os := range offerStatuses {
		offerStatusesMap[os.Code] = os.Name
	}

	appointmentStatusesMap := make(map[string]string)
	for _, as := range appointmentStatuses {
		appointmentStatusesMap[as.Code] = as.Name
	}

	scanStatusesMap := make(map[string]string)
	for _, ss := range scanStatuses {
		scanStatusesMap[ss.Code] = ss.Name
	}

	urgencyLevelsMap := make(map[string]string)
	for _, ul := range urgencyLevels {
		urgencyLevelsMap[ul.Code] = ul.Name
	}

	gendersMap := make(map[string]string)
	for _, g := range genders {
		gendersMap[g.Code] = g.Name
	}

	priceSegmentsList := make([]string, len(priceSegments))
	for i, ps := range priceSegments {
		priceSegmentsList[i] = ps.Name
	}

	citiesList := make([]string, len(cities))
	for i, c := range cities {
		citiesList[i] = c.Name
	}

	// Group districts by city
	districtsByCity := make(map[string][]string)
	for _, d := range districts {
		cityName := d.City.Name
		districtsByCity[cityName] = append(districtsByCity[cityName], d.Name)
	}

	constants := gin.H{
		"roles":                rolesMap,
		"specializations":      specsMap,
		"treatment_statuses":   treatmentStatusesMap,
		"offer_statuses":       offerStatusesMap,
		"appointment_statuses": appointmentStatusesMap,
		"scan_statuses":        scanStatusesMap,
		"urgency_levels":       urgencyLevelsMap,
		"genders":              gendersMap,
		"price_segments":       priceSegmentsList,
		"cities":               citiesList,
		"districts_by_city":    districtsByCity,
	}

	c.JSON(http.StatusOK, constants)
}
