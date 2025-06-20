package repository

import (
	"context"

	"github.com/yesetoda/Sera_Ale/internal/domain"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	Update(ctx context.Context, job *domain.Job) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*domain.Job, error)
	FindByCompany(ctx context.Context, companyID string, page, size int) ([]domain.Job, int64, error)
	Search(ctx context.Context, filters map[string]interface{}, page, size int) ([]domain.Job, int64, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) Update(ctx context.Context, job *domain.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Job{}, "id = ?", id).Error
}

func (r *jobRepository) FindByID(ctx context.Context, id string) (*domain.Job, error) {
	var job domain.Job
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) FindByCompany(ctx context.Context, companyID string, page, size int) ([]domain.Job, int64, error) {
	var jobs []domain.Job
	var total int64
	db := r.db.WithContext(ctx).Where("created_by = ?", companyID)
	db.Model(&domain.Job{}).Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&jobs).Error
	return jobs, total, err
}

func (r *jobRepository) Search(ctx context.Context, filters map[string]interface{}, page, size int) ([]domain.Job, int64, error) {
	var jobs []domain.Job
	var total int64
	db := r.db.WithContext(ctx).Model(&domain.Job{})
	if title, ok := filters["title"]; ok {
		db = db.Where("LOWER(title) LIKE ?", "%"+title.(string)+"%")
	}
	if location, ok := filters["location"]; ok {
		db = db.Where("location LIKE ?", "%"+location.(string)+"%")
	}
	// companyName filter requires a join, implement in service layer if needed
	db.Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&jobs).Error
	return jobs, total, err
}
