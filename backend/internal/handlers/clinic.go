package handlers

import (
	"dental-marketplace/backend/internal/models"
	"dental-marketplace/backend/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	repo *repository.Repository
}

func NewClinicHandler(repo *repository.Repository) *ClinicHandler {
	return &ClinicHandler{repo: repo}
}

// GetDashboard retrieves dashboard metrics for clinic
// @Summary Get clinic dashboard
// @Description Get key metrics for clinic dashboard
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/dashboard [get]
func (h *ClinicHandler) GetDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	// Parse period
	period := c.DefaultQuery("period", "30d")
	var days int
	switch period {
	case "7d":
		days = 7
	case "90d":
		days = 90
	default:
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	metrics, err := h.repo.GetClinicDashboardMetrics(clinic.ID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve dashboard metrics",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetIncomingPlans retrieves treatment plans for clinic to review
// @Summary Get incoming treatment plans
// @Description Get treatment plans that clinic can bid on
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.TreatmentPlan
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/incoming-plans [get]
func (h *ClinicHandler) GetIncomingPlans(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	plans, err := h.repo.GetIncomingTreatmentPlans(clinic.ID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve treatment plans",
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// CreateOffer creates a clinic offer for a treatment plan
type CreateOfferRequest struct {
	TreatmentPlanID   uint   `json:"treatment_plan_id" binding:"required"`
	TherapyCost       int    `json:"therapy_cost"`
	OrthopedicsCost   int    `json:"orthopedics_cost"`
	SurgeryCost       int    `json:"surgery_cost"`
	HygieneCost       int    `json:"hygiene_cost"`
	PeriodonticsCost  int    `json:"periodontics_cost"`
	TotalCost         int    `json:"total_cost" binding:"required"`
	EstimatedDuration string `json:"estimated_duration"`
	InstallmentMonths int    `json:"installment_months"`
	WarrantyDetails   string `json:"warranty_details"`
	Notes             string `json:"notes"`
}

// @Summary Create clinic offer
// @Description Create an offer for a treatment plan
// @Tags clinic
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOfferRequest true "Offer details"
// @Success 201 {object} models.ClinicOffer
// @Failure 400 {object} ErrorResponse
// @Router /api/clinic/offers [post]
func (h *ClinicHandler) CreateOffer(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	offer := &models.ClinicOffer{
		TreatmentPlanID:   req.TreatmentPlanID,
		ClinicID:          clinic.ID,
		Status:            models.OfferStatusSent,
		TherapyCost:       req.TherapyCost,
		OrthopedicsCost:   req.OrthopedicsCost,
		SurgeryCost:       req.SurgeryCost,
		HygieneCost:       req.HygieneCost,
		PeriodonticsCost:  req.PeriodonticsCost,
		TotalCost:         req.TotalCost,
		EstimatedDuration: req.EstimatedDuration,
		InstallmentMonths: req.InstallmentMonths,
		WarrantyDetails:   req.WarrantyDetails,
		Notes:             req.Notes,
	}

	err = h.repo.CreateClinicOffer(offer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create offer",
		})
		return
	}

	// Update treatment plan status
	h.repo.UpdateTreatmentPlanStatus(req.TreatmentPlanID, models.PlanStatusOffersReceived)

	c.JSON(http.StatusCreated, offer)
}

// GetLeads retrieves accepted offers (leads) for clinic
// @Summary Get clinic leads
// @Description Get all accepted offers (leads) for the clinic
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ClinicOffer
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/leads [get]
func (h *ClinicHandler) GetLeads(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	leads, err := h.repo.GetClinicLeads(clinic.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve leads",
		})
		return
	}

	c.JSON(http.StatusOK, leads)
}

// GetAppointments retrieves appointments for clinic
// @Summary Get clinic appointments
// @Description Get all appointments for the clinic
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Appointment
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/appointments [get]
func (h *ClinicHandler) GetAppointments(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	status := c.Query("status")
	appointments, err := h.repo.GetClinicAppointments(clinic.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve appointments",
		})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// UpdateAppointment updates appointment status
type UpdateAppointmentRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// @Summary Update appointment
// @Description Update appointment status and notes
// @Tags clinic
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Appointment ID"
// @Param request body UpdateAppointmentRequest true "Update details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/clinic/appointments/{id} [put]
func (h *ClinicHandler) UpdateAppointment(c *gin.Context) {
	appointmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid appointment ID",
		})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err = h.repo.UpdateAppointmentStatus(uint(appointmentID), req.Status, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update appointment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment updated successfully",
	})
}

// GetPriceList retrieves clinic price list
// @Summary Get price list
// @Description Get clinic's price list by specialization
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Param specialization query string false "Filter by specialization"
// @Success 200 {array} models.PriceList
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/price-list [get]
func (h *ClinicHandler) GetPriceList(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	specialization := c.Query("specialization")
	priceList, err := h.repo.GetClinicPriceList(clinic.ID, specialization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve price list",
		})
		return
	}

	c.JSON(http.StatusOK, priceList)
}

// UpdatePriceList updates clinic price list
// @Summary Update price list
// @Description Update clinic's price list
// @Tags clinic
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []models.PriceList true "Price list items"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/clinic/price-list [put]
func (h *ClinicHandler) UpdatePriceList(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	var items []models.PriceList
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Set clinic ID for all items
	for i := range items {
		items[i].ClinicID = clinic.ID
	}

	err = h.repo.UpdatePriceList(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update price list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Price list updated successfully",
	})
}

// GetAnalytics retrieves analytics data for clinic
// @Summary Get clinic analytics
// @Description Get revenue and performance analytics
// @Tags clinic
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/clinic/analytics [get]
func (h *ClinicHandler) GetAnalytics(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	clinic, err := h.repo.GetClinicByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic profile not found",
		})
		return
	}

	// Parse period
	period := c.DefaultQuery("period", "30d")
	var days int
	switch period {
	case "7d":
		days = 7
	case "90d":
		days = 90
	default:
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get statistics
	stats, err := h.repo.GetStatistics(startDate, endDate, &clinic.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period":     period,
		"statistics": stats,
		"clinic":     clinic,
	})
}
