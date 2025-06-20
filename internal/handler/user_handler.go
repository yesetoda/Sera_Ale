package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yesetoda/Sera_Ale/internal/app"
	"github.com/yesetoda/Sera_Ale/internal/domain"
)

type UserHandler struct {
	App app.UserApp
}

func NewUserHandler(app app.UserApp) *UserHandler {
	return &UserHandler{App: app}
}

// Signup godoc
// @Summary Signup
// @Description Register as a new user (company or applicant)
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body signupRequest true "Signup info"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /signup [post]
// @Security None
// No authentication required
type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Invalid input", Errors: []string{"Invalid JSON"}})
		return
	}
	user, errs := h.App.Signup(c.Request.Context(), req.Name, req.Email, req.Password, req.Role)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Signup failed", Errors: errs})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Signup successful", Object: user})
}

// Login godoc
// @Summary Login
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body loginRequest true "Login info"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /login [post]
// @Security None
// No authentication required
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Invalid input", Errors: []string{"Invalid JSON"}})
		return
	}
	user, token, errs := h.App.Login(c.Request.Context(), req.Email, req.Password)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, domain.BaseResponse{Success: false, Message: "Login failed", Errors: errs})
		return
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Login successful", Object: gin.H{"user": user, "token": token}})
}
