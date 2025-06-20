package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// OpenAPI3: summary: Apply for a job
// OpenAPI3: description: Applicant applies for a job with resume upload (requires Bearer token)
// OpenAPI3: tags: [Applications]
// OpenAPI3: security: BearerAuth
// OpenAPI3: requestBody: multipart/form-data (job_id, cover_letter, resume)
// OpenAPI3: responses: 200=BaseResponse, 400=BaseResponse
// @Security BearerAuth
// Requires Bearer token
func (h *ApplicationHandler) Apply(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.JSON(401, gin.H{"success": false, "message": "Missing or invalid Bearer token in Authorization header. Please provide: Authorization: Bearer <token>"})
		return
	}
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
// @Description Applicant views their applications. Requires Bearer token in Authorization header.
// @Tags Applications
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} map[string]interface{} "List of applications"
// @Security BearerAuth
// @Router /applicant/applications [get]
// Requires Bearer token
func (h *ApplicationHandler) TrackApplications(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.JSON(401, gin.H{"success": false, "message": "Missing or invalid Bearer token in Authorization header. Please provide: Authorization: Bearer <token>"})
		return
	}
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
// @Description Company views applications for their job. Requires Bearer token in Authorization header.
// @Tags Applications
// @Accept json
// @Produce json
// @Param job_id query string true "Job ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} map[string]interface{} "List of applications"
// @Failure 403 {object} map[string]interface{} "Unauthorized or not job owner"
// @Security BearerAuth
// @Router /company/applications/job [get]
// Requires Bearer token
func (h *ApplicationHandler) GetApplicationsForJob(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.JSON(401, gin.H{"success": false, "message": "Missing or invalid Bearer token in Authorization header. Please provide: Authorization: Bearer <token>"})
		return
	}
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
// @Description Company updates application status. Requires Bearer token in Authorization header.
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param status body updateStatusRequest true "New status"
// @Success 200 {object} map[string]interface{} "Status updated"
// @Failure 400 {object} map[string]interface{} "Validation or update error"
// @Failure 403 {object} map[string]interface{} "Unauthorized or not job owner"
// @Security BearerAuth
// @Router /company/applications/{id}/status [put]
// Requires Bearer token
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.JSON(401, gin.H{"success": false, "message": "Missing or invalid Bearer token in Authorization header. Please provide: Authorization: Bearer <token>"})
		return
	}
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
