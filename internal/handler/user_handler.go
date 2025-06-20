package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yesetoda/Sera_Ale/internal/app"
	"github.com/yesetoda/Sera_Ale/internal/domain"
)

// User model
// @Description User entity
// @Description Role: "applicant" or "company"
// @Description Password is always hashed
// @Description Example: {"id": "uuid", "name": "John Doe", "email": "john@example.com", "role": "applicant"}
type UserSwagger struct {
	ID    string `json:"id" example:"uuid"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
	Role  string `json:"role" example:"applicant"`
}

type UserHandler struct {
	App app.UserApp
}

type AuthHandler struct {
	App app.UserApp
}

func NewUserHandler(app app.UserApp) *UserHandler {
	return &UserHandler{App: app}
}

func NewAuthHandler(app app.UserApp) *AuthHandler {
	return &AuthHandler{App: app}
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup godoc
// @Summary Register as a new user (company or applicant)
// @Description Register as a new user (company or applicant)
// @Tags Auth
// @Accept json
// @Produce json
// @Param signupRequest body signupRequest true "Signup request"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
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
	resp := UserSwagger{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.Name,
	}
	c.JSON(http.StatusOK, domain.BaseResponse{Success: true, Message: "Signup successful", Object: resp})
}

// Login godoc
// @Summary Login with email and password
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login request"
// @Success 200 {object} domain.BaseResponse
// @Failure 400 {object} domain.BaseResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
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

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Returns the authenticated user's profile (requires Bearer token)
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} UserSwagger
// @Failure 404 {object} domain.BaseResponse
// @Security BearerAuth
// @Router /user/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.JSON(401, gin.H{"success": false, "message": "Missing or invalid Bearer token in Authorization header. Please provide: Authorization: Bearer <token>"})
		return
	}
	userID := c.GetString("user_id")
	user, err := h.App.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
