package domain

import (
	"github.com/google/uuid"
)

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"unique;not null" json:"name"`
}

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	RoleID   uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	Role     Role      `gorm:"foreignKey:RoleID" json:"role"`
}
