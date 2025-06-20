package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yesetoda/Sera_Ale/internal/app"
	"github.com/yesetoda/Sera_Ale/internal/domain"
)

type JobHandler struct {
	App app.JobApp
}

func NewJobHandler(app app.JobApp) *JobHandler {
	return &JobHandler{App: app}
}

type jobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// CreateJob godoc
// @Summary Create job
// @Description Company posts a new job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param job body jobRequest true "Job info"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /jobs [post]
func (h *JobHandler) CreateJob(c *gin.Context) {
	var req jobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Invalid input", Errors: []string{"Invalid JSON"}})
		return
	}
	if len(req.Title) < 1 || len(req.Title) > 100 {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Title must be 1-100 characters"})
		return
	}
	if len(req.Description) < 20 || len(req.Description) > 2000 {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Description must be 20-2000 characters"})
		return
	}
	userID := c.GetString("user_id")
	job := &domain.Job{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		CreatedBy:   uuid.MustParse(userID),
	}
	if err := h.App.CreateJob(c.Request.Context(), job); err != nil {
		c.JSON(http.StatusInternalServerError, domain.BaseResponse{Success: false, Message: "Failed to create job"})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Job created", Object: job})
}

// UpdateJob godoc
// @Summary Update job
// @Description Company updates their job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param job body jobRequest true "Job info"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /jobs/{id} [put]
func (h *JobHandler) UpdateJob(c *gin.Context) {
	id := c.Param("id")
	var req jobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Invalid input"})
		return
	}
	userID := c.GetString("user_id")
	job, err := h.App.GetJobByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.BaseResponse{Success: false, Message: "Job not found"})
		return
	}
	if job.CreatedBy.String() != userID {
		c.JSON(http.StatusForbidden, domain.BaseResponse{Success: false, Message: "Unauthorized access"})
		return
	}
	job.Title = req.Title
	job.Description = req.Description
	job.Location = req.Location
	if err := h.App.UpdateJob(c.Request.Context(), job); err != nil {
		c.JSON(http.StatusInternalServerError, domain.BaseResponse{Success: false, Message: "Failed to update job"})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Job updated", Object: job})
}

// DeleteJob godoc
// @Summary Delete job
// @Description Company deletes their job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.App.DeleteJob(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusForbidden, domain.BaseResponse{Success: false, Message: "Unauthorized access"})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Job deleted"})
}

// GetJob godoc
// @Summary Get job details
// @Description Get job details by ID
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} domain.BaseResponse
// @Failure 404 {object} domain.BaseResponse
// @Router /jobs/{id} [get]
func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.App.GetJobByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.BaseResponse{Success: false, Message: "Job not found"})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Job found", Object: job})
}

// SearchJobs godoc
// @Summary Search jobs
// @Description Applicant searches jobs with filters and pagination
// @Tags Jobs
// @Accept json
// @Produce json
// @Param title query string false "Title filter"
// @Param location query string false "Location filter"
// @Param company_name query string false "Company name filter"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} domain.PaginatedResponse
// @Router /jobs [get]
func (h *JobHandler) SearchJobs(c *gin.Context) {
	filters := map[string]interface{}{}
	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if company := c.Query("company_name"); company != "" {
		filters["company_name"] = company
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	jobs, total, err := h.App.SearchJobs(c.Request.Context(), filters, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.PaginatedResponse{Success: false, Message: "Failed to search jobs"})
		return
	}
	c.JSON(http.StatusOK, domain.PaginatedResponse{
		Success:    true,
		Message:    "Jobs found",
		Object:     jobs,
		PageNumber: page,
		PageSize:   size,
		TotalSize:  int(total),
	})
}
