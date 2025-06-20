package domain

import (
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleApplicant UserRole = "applicant"
	RoleCompany   UserRole = "company"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     UserRole  `json:"role"`
}
