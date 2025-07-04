package app

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/yesetoda/Sera_Ale/internal/domain"
	"github.com/yesetoda/Sera_Ale/internal/repository"
	"github.com/yesetoda/Sera_Ale/internal/service"
	"gorm.io/gorm"
)

type UserApp interface {
	Signup(ctx context.Context, name, email, password, role string) (*domain.User, []string)
	Login(ctx context.Context, email, password string) (*domain.User, string, []string)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type userApp struct {
	repo     repository.UserRepository
	jwt      service.JWTService
	password service.PasswordService
	orm      *gorm.DB
}

func NewUserApp(repo repository.UserRepository, jwt service.JWTService, password service.PasswordService) UserApp {
	return &userApp{repo: repo, jwt: jwt, password: password, orm: repo.GetDB()}
}

func (a *userApp) Signup(ctx context.Context, name, email, password, roleName string) (*domain.User, []string) {
	errMsgs := validateSignupInput(name, email, password, roleName)
	if len(errMsgs) > 0 {
		return nil, errMsgs
	}
	_, err := a.repo.FindByEmail(ctx, email)
	if err == nil {
		return nil, []string{"Email already exists"}
	}
	var role domain.Role
	if err := a.orm.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(roleName)).First(&role).Error; err != nil {
		return nil, []string{"Role does not exist"}
	}
	hash, err := a.password.HashPassword(password)
	if err != nil {
		return nil, []string{"Failed to hash password"}
	}
	user := &domain.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: hash,
		RoleID:   role.ID,
		Role:     role,
	}
	if err := a.repo.Create(ctx, user); err != nil {
		return nil, []string{"Failed to create user"}
	}
	return user, nil
}

func (a *userApp) Login(ctx context.Context, email, password string) (*domain.User, string, []string) {
	user, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", []string{"User not found"}
	}
	if err := a.password.ComparePassword(user.Password, password); err != nil {
		return nil, "", []string{"Incorrect password"}
	}
	token, err := a.jwt.GenerateToken(user.ID.String(), user.Role.Name)
	if err != nil {
		return nil, "", []string{"Failed to generate token"}
	}
	return user, token, nil
}

func (a *userApp) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return a.repo.FindByID(ctx, id)
}

func validateSignupInput(name, email, password, role string) []string {
	errs := []string{}
	if name == "" || !regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(name) {
		errs = append(errs, "Name must contain only alphabets and spaces")
	}
	if email == "" || !regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`).MatchString(email) {
		errs = append(errs, "Invalid email address")
	}
	if len(password) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(password) ||
		!regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password) {
		errs = append(errs, "Password must be at least 8 characters, include upper, lower, number, and special char")
	}
	role = strings.ToLower(role)
	if role != "company" && role != "applicant" {
		errs = append(errs, "Role must be either 'company' or 'applicant'")
	}
	return errs
}
