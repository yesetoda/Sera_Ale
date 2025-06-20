package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	StatusApplied   ApplicationStatus = "Applied"
	StatusReviewed  ApplicationStatus = "Reviewed"
	StatusInterview ApplicationStatus = "Interview"
	StatusRejected  ApplicationStatus = "Rejected"
	StatusHired     ApplicationStatus = "Hired"
)

type Application struct {
	ID          uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ApplicantID uuid.UUID         `json:"applicant_id"`
	JobID       uuid.UUID         `json:"job_id"`
	ResumeLink  string            `json:"resume_link"`
	CoverLetter string            `json:"cover_letter"`
	Status      ApplicationStatus `json:"status"`
	AppliedAt   time.Time         `json:"applied_at"`
}
