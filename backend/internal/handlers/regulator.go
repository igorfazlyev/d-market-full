package handlers

import (
	"dental-marketplace/backend/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RegulatorHandler struct {
	repo *repository.Repository
}

func NewRegulatorHandler(repo *repository.Repository) *RegulatorHandler {
	return &RegulatorHandler{repo: repo}
}

// GetDashboard retrieves regional overview dashboard
// @Summary Get regulator dashboard
// @Description Get regional overview with key metrics
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/regulator/dashboard [get]
func (h *RegulatorHandler) GetDashboard(c *gin.Context) {
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

	// Get regional statistics (clinic_id IS NULL)
	stats, err := h.repo.GetStatistics(startDate, endDate, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve dashboard data",
		})
		return
	}

	// Calculate totals
	var totalPlans, totalAppointments, totalRevenue, totalPatients int
	var totalCaries, totalPulpitis, totalPeriodontitis, totalGingivitis, totalParodontitis int
	var avgWaitDays, avgTreatmentCost float64

	for _, stat := range stats {
		totalPlans += stat.TreatmentPlansGenerated
		totalAppointments += stat.AppointmentsCompleted
		totalRevenue += stat.TotalRevenue
		totalPatients += stat.PatientCount
		totalCaries += stat.CariesCount
		totalPulpitis += stat.PulpitisCount
		totalPeriodontitis += stat.PeriodontitisCount
		totalGingivitis += stat.GingivitisCount
		totalParodontitis += stat.ParodontitisCount
		avgWaitDays += stat.AverageWaitDays
		avgTreatmentCost += float64(stat.AverageTreatmentCost)
	}

	count := len(stats)
	if count > 0 {
		avgWaitDays /= float64(count)
		avgTreatmentCost /= float64(count)
	}

	// Get clinic count
	clinics, _ := h.repo.GetClinics("", "", "")
	clinicCount := len(clinics)

	c.JSON(http.StatusOK, gin.H{
		"period": period,
		"summary": gin.H{
			"total_clinics":             clinicCount,
			"total_treatment_plans":     totalPlans,
			"total_appointments":        totalAppointments,
			"total_revenue":             totalRevenue,
			"total_patients":            totalPatients,
			"average_wait_days":         avgWaitDays,
			"average_treatment_cost":    int(avgTreatmentCost),
		},
		"disease_statistics": gin.H{
			"caries":        totalCaries,
			"pulpitis":      totalPulpitis,
			"periodontitis": totalPeriodontitis,
			"gingivitis":    totalGingivitis,
			"parodontitis":  totalParodontitis,
		},
		"time_series": stats,
	})
}

// GetStatistics retrieves filtered statistics
// @Summary Get statistics
// @Description Get statistics with filters
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Param clinic_id query int false "Filter by clinic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/regulator/statistics [get]
func (h *RegulatorHandler) GetStatistics(c *gin.Context) {
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

	// Parse clinic ID filter if provided
	var clinicID *uint
	if clinicIDStr := c.Query("clinic_id"); clinicIDStr != "" {
		id, err := strconv.ParseUint(clinicIDStr, 10, 32)
		if err == nil {
			cID := uint(id)
			clinicID = &cID
		}
	}

	// Get statistics
	stats, err := h.repo.GetStatistics(startDate, endDate, clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve statistics",
		})
		return
	}

	// If clinic-specific, get clinic info
	var clinicInfo interface{}
	if clinicID != nil {
		clinic, err := h.repo.GetClinicByID(*clinicID)
		if err == nil {
			clinicInfo = clinic
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"period":     period,
		"clinic":     clinicInfo,
		"statistics": stats,
	})
}

// GetClinics retrieves all clinics with statistics
// @Summary Get all clinics
// @Description Get list of all clinics with their statistics
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param city query string false "Filter by city"
// @Param district query string false "Filter by district"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/regulator/clinics [get]
func (h *RegulatorHandler) GetClinics(c *gin.Context) {
	city := c.Query("city")
	district := c.Query("district")

	clinics, err := h.repo.GetAllClinicsWithStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve clinics",
		})
		return
	}

	// Filter by city and district if provided
	if city != "" || district != "" {
		filtered := make([]map[string]interface{}, 0)
		for _, clinic := range clinics {
			match := true
			if city != "" && clinic["city"] != city {
				match = false
			}
			if district != "" && clinic["district"] != district {
				match = false
			}
			if match {
				filtered = append(filtered, clinic)
			}
		}
		clinics = filtered
	}

	c.JSON(http.StatusOK, clinics)
}

// GetComplaints retrieves all complaints
// @Summary Get complaints
// @Description Get all complaints with optional status filter
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Complaint
// @Failure 500 {object} ErrorResponse
// @Router /api/regulator/complaints [get]
func (h *RegulatorHandler) GetComplaints(c *gin.Context) {
	status := c.Query("status")

	complaints, err := h.repo.GetAllComplaints(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve complaints",
		})
		return
	}

	c.JSON(http.StatusOK, complaints)
}

// GetClinicDetails retrieves detailed information about a specific clinic
// @Summary Get clinic details
// @Description Get detailed information and statistics for a specific clinic
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param id path int true "Clinic ID"
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} ErrorResponse
// @Router /api/regulator/clinics/{id} [get]
func (h *RegulatorHandler) GetClinicDetails(c *gin.Context) {
	clinicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid clinic ID",
		})
		return
	}

	// Get clinic info
	clinic, err := h.repo.GetClinicByID(uint(clinicID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Clinic not found",
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

	// Get clinic statistics
	cID := clinic.ID
	stats, err := h.repo.GetStatistics(startDate, endDate, &cID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve clinic statistics",
		})
		return
	}

	// Get price list
	priceList, _ := h.repo.GetClinicPriceList(clinic.ID, "")

	// Get appointments count
	appointments, _ := h.repo.GetClinicAppointments(clinic.ID, "")

	c.JSON(http.StatusOK, gin.H{
		"clinic":            clinic,
		"period":            period,
		"statistics":        stats,
		"price_list_count":  len(priceList),
		"appointments_count": len(appointments),
	})
}

// GetDiseaseAnalytics retrieves disease prevalence analytics
// @Summary Get disease analytics
// @Description Get disease prevalence statistics for a period
// @Tags regulator
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (7d, 30d, 90d)" default(30d)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/regulator/disease-analytics [get]
func (h *RegulatorHandler) GetDiseaseAnalytics(c *gin.Context) {
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

	// Get regional statistics
	stats, err := h.repo.GetStatistics(startDate, endDate, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve disease analytics",
		})
		return
	}

	// Aggregate disease counts
	diseases := map[string]int{
		"Caries":        0,
		"Pulpitis":      0,
		"Periodontitis": 0,
		"Gingivitis":    0,
		"Parodontitis":  0,
	}

	for _, stat := range stats {
		diseases["Caries"] += stat.CariesCount
		diseases["Pulpitis"] += stat.PulpitisCount
		diseases["Periodontitis"] += stat.PeriodontitisCount
		diseases["Gingivitis"] += stat.GingivitisCount
		diseases["Parodontitis"] += stat.ParodontitisCount
	}

	// Calculate total cases
	totalCases := 0
	for _, count := range diseases {
		totalCases += count
	}

	// Create chart data
	chartData := make([]map[string]interface{}, 0)
	for disease, count := range diseases {
		percentage := 0.0
		if totalCases > 0 {
			percentage = float64(count) / float64(totalCases) * 100
		}
		chartData = append(chartData, map[string]interface{}{
			"disease":    disease,
			"count":      count,
			"percentage": percentage,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"period":      period,
		"total_cases": totalCases,
		"diseases":    chartData,
		"time_series": stats,
	})
}
