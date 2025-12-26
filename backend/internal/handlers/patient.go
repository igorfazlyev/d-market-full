package handlers

import (
	"dental-marketplace/backend/internal/models"
	"dental-marketplace/backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	repo *repository.Repository
}

func NewPatientHandler(repo *repository.Repository) *PatientHandler {
	return &PatientHandler{repo: repo}
}

// GetScans retrieves all CT scans for the patient
// @Summary Get patient CT scans
// @Description Get all CT scans for authenticated patient
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CTScan
// @Failure 401 {object} ErrorResponse
// @Router /api/patient/scans [get]
func (h *PatientHandler) GetScans(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Get patient profile
	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	// Get scans
	scans, err := h.repo.GetPatientCTScans(patient.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve scans",
		})
		return
	}

	c.JSON(http.StatusOK, scans)
}

// GetScanByID retrieves a specific CT scan with treatment plan
// @Summary Get CT scan by ID
// @Description Get CT scan details with treatment plan
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Param id path int true "Scan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} ErrorResponse
// @Router /api/patient/scans/{id} [get]
func (h *PatientHandler) GetScanByID(c *gin.Context) {
	scanID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid scan ID",
		})
		return
	}

	scan, err := h.repo.GetCTScanByID(uint(scanID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Scan not found",
		})
		return
	}

	// Get treatment plan if exists
	var treatmentPlan *models.TreatmentPlan
	if scan.AIProcessed {
		treatmentPlan, _ = h.repo.GetTreatmentPlanByScanID(scan.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"scan":           scan,
		"treatment_plan": treatmentPlan,
	})
}

// GetTreatmentPlan retrieves treatment plan with items
// @Summary Get treatment plan
// @Description Get treatment plan by scan ID with all items
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Param id path int true "Scan ID"
// @Success 200 {object} models.TreatmentPlan
// @Failure 404 {object} ErrorResponse
// @Router /api/patient/scans/{id}/plan [get]
func (h *PatientHandler) GetTreatmentPlan(c *gin.Context) {
	scanID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid scan ID",
		})
		return
	}

	plan, err := h.repo.GetTreatmentPlanByScanID(uint(scanID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Treatment plan not found",
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// UpdateSearchCriteria updates patient's clinic search preferences
type SearchCriteriaRequest struct {
	City         string `json:"city"`
	District     string `json:"district"`
	PriceSegment string `json:"price_segment"`
}

// @Summary Update search criteria
// @Description Update patient's clinic search preferences
// @Tags patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SearchCriteriaRequest true "Search criteria"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/patient/search-criteria [post]
func (h *PatientHandler) UpdateSearchCriteria(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req SearchCriteriaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Get patient profile
	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	// Update criteria
	err = h.repo.UpdatePatientSearchCriteria(patient.ID, req.City, req.District, req.PriceSegment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update search criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search criteria updated successfully",
	})
}

// GetOffers retrieves clinic offers for a treatment plan
// @Summary Get clinic offers
// @Description Get all clinic offers for a treatment plan
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Param plan_id path int true "Treatment Plan ID"
// @Success 200 {array} models.ClinicOffer
// @Failure 404 {object} ErrorResponse
// @Router /api/patient/plans/{plan_id}/offers [get]
func (h *PatientHandler) GetOffers(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("plan_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	offers, err := h.repo.GetOffersForTreatmentPlan(uint(planID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve offers",
		})
		return
	}

	c.JSON(http.StatusOK, offers)
}

// SelectOffer accepts a clinic offer
type SelectOfferRequest struct {
	OfferID uint `json:"offer_id" binding:"required"`
}

// @Summary Select clinic offer
// @Description Accept a clinic offer and create appointment
// @Tags patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SelectOfferRequest true "Offer selection"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/patient/select-offer [post]
func (h *PatientHandler) SelectOffer(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req SelectOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Get patient profile
	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	// Accept offer and create appointment
	err = h.repo.AcceptClinicOffer(req.OfferID, patient.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to accept offer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Offer accepted successfully. Appointment created.",
	})
}

// GetAppointments retrieves all appointments for patient
// @Summary Get appointments
// @Description Get all appointments for authenticated patient
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Appointment
// @Failure 500 {object} ErrorResponse
// @Router /api/patient/appointments [get]
func (h *PatientHandler) GetAppointments(c *gin.Context) {
	userID, _ := c.Get("userID")

	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	appointments, err := h.repo.GetPatientAppointments(patient.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve appointments",
		})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// CreateReview creates a review for a clinic
type CreateReviewRequest struct {
	ClinicID uint   `json:"clinic_id" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comment  string `json:"comment"`
}

// @Summary Create review
// @Description Create a review for a clinic
// @Tags patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateReviewRequest true "Review details"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/patient/reviews [post]
func (h *PatientHandler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	review := &models.Review{
		PatientID: patient.ID,
		ClinicID:  req.ClinicID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		IsPublic:  false, // Reviews are not public by default per requirements
	}

	err = h.repo.CreateReview(review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create review",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review submitted successfully",
	})
}

// CreateComplaint creates a complaint
type CreateComplaintRequest struct {
	ClinicID    uint   `json:"clinic_id" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// @Summary Create complaint
// @Description Create a complaint about a clinic
// @Tags patient
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateComplaintRequest true "Complaint details"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/patient/complaints [post]
func (h *PatientHandler) CreateComplaint(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	complaint := &models.Complaint{
		PatientID:   patient.ID,
		ClinicID:    req.ClinicID,
		Subject:     req.Subject,
		Description: req.Description,
		Status:      "open",
	}

	err = h.repo.CreateComplaint(complaint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create complaint",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Complaint submitted successfully",
	})
}

// GetTreatmentPlans retrieves all treatment plans for patient
// @Summary Get treatment plans
// @Description Get all treatment plans for authenticated patient
// @Tags patient
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.TreatmentPlan
// @Failure 500 {object} ErrorResponse
// @Router /api/patient/plans [get]
func (h *PatientHandler) GetTreatmentPlans(c *gin.Context) {
	userID, _ := c.Get("userID")

	patient, err := h.repo.GetPatientByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient profile not found",
		})
		return
	}

	plans, err := h.repo.GetPatientTreatmentPlans(patient.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve treatment plans",
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}
