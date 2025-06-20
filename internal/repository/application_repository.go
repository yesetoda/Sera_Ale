package repository

import (
	"context"

	"github.com/yesetoda/Sera_Ale/internal/domain"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) error
	FindByApplicant(ctx context.Context, applicantID string, page, size int) ([]domain.Application, int64, error)
	FindByJob(ctx context.Context, jobID string, page, size int) ([]domain.Application, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Application, error)
	UpdateStatus(ctx context.Context, id string, status domain.ApplicationStatus) error
	FindByApplicantAndJob(ctx context.Context, applicantID, jobID string) (*domain.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) FindByApplicant(ctx context.Context, applicantID string, page, size int) ([]domain.Application, int64, error) {
	var apps []domain.Application
	var total int64
	db := r.db.WithContext(ctx).Where("applicant_id = ?", applicantID)
	db.Model(&domain.Application{}).Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) FindByJob(ctx context.Context, jobID string, page, size int) ([]domain.Application, int64, error) {
	var apps []domain.Application
	var total int64
	db := r.db.WithContext(ctx).Where("job_id = ?", jobID)
	db.Model(&domain.Application{}).Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) FindByID(ctx context.Context, id string) (*domain.Application, error) {
	var app domain.Application
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, id string, status domain.ApplicationStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Application{}).Where("id = ?", id).Update("status", status).Error
}

func (r *applicationRepository) FindByApplicantAndJob(ctx context.Context, applicantID, jobID string) (*domain.Application, error) {
	var app domain.Application
	err := r.db.WithContext(ctx).Where("applicant_id = ? AND job_id = ?", applicantID, jobID).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}
