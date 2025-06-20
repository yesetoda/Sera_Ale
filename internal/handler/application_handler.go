package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yesetoda/Sera_Ale/internal/app"
)

type ApplicationHandler struct {
	App app.ApplicationApp
}

func NewApplicationHandler(app app.ApplicationApp) *ApplicationHandler {
	return &ApplicationHandler{App: app}
}

type applyRequest struct {
	JobID       string `form:"job_id" json:"job_id"`
	CoverLetter string `form:"cover_letter" json:"cover_letter"`
}

// Apply godoc
// @Summary Apply for a job
// @Description Applicant applies for a job with resume upload
// @Tags Applications
// @Accept multipart/form-data
// @Produce json
// @Param job_id formData string true "Job ID"
// @Param cover_letter formData string false "Cover letter"
// @Param resume formData file true "Resume PDF"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /applicant/applications [post]
func (h *ApplicationHandler) Apply(c *gin.Context) {
	applicantID := c.GetString("user_id")
	var req applyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}
	file, _, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Resume file required"})
		return
	}
	defer file.Close()
	application, errs := h.App.Apply(c.Request.Context(), applicantID, req.JobID, req.CoverLetter, file)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Application failed", "errors": errs})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Application submitted", "object": application})
}

// TrackApplications godoc
// @Summary Track my applications
// @Description Applicant views their applications
// @Tags Applications
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /applicant/applications [get]
func (h *ApplicationHandler) TrackApplications(c *gin.Context) {
	applicantID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	apps, total, err := h.App.TrackApplications(c.Request.Context(), applicantID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Applications found", "object": apps, "page_number": page, "page_size": size, "total_size": total})
}

// GetApplicationsForJob godoc
// @Summary View applications for a job
// @Description Company views applications for their job
// @Tags Applications
// @Accept json
// @Produce json
// @Param job_id query string true "Job ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /company/applications/job [get]
func (h *ApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	companyID := c.GetString("user_id")
	jobID := c.Query("job_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	apps, total, err := h.App.GetApplicationsForJob(c.Request.Context(), jobID, companyID, page, size)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Applications found", "object": apps, "page_number": page, "page_size": size, "total_size": total})
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

// UpdateStatus godoc
// @Summary Update application status
// @Description Company updates application status
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param status body updateStatusRequest true "New status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /company/applications/{id}/status [put]
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	companyID := c.GetString("user_id")
	id := c.Param("id")
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}
	app, err := h.App.UpdateStatus(c.Request.Context(), id, companyID, req.Status)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status updated", "object": app})
}

// ... existing code ...
