package app

import (
	"context"

	"github.com/yesetoda/Sera_Ale/internal/domain"
	"github.com/yesetoda/Sera_Ale/internal/repository"
)

type JobApp interface {
	CreateJob(ctx context.Context, job *domain.Job) error
	UpdateJob(ctx context.Context, job *domain.Job) error
	DeleteJob(ctx context.Context, jobID string, userID string) error
	GetJobByID(ctx context.Context, jobID string) (*domain.Job, error)
	GetJobsByCompany(ctx context.Context, companyID string, page, size int) ([]domain.Job, int64, error)
	SearchJobs(ctx context.Context, filters map[string]interface{}, page, size int) ([]domain.Job, int64, error)
}

type jobApp struct {
	repo repository.JobRepository
}

func NewJobApp(repo repository.JobRepository) JobApp {
	return &jobApp{repo: repo}
}

func (a *jobApp) CreateJob(ctx context.Context, job *domain.Job) error {
	return a.repo.Create(ctx, job)
}

func (a *jobApp) UpdateJob(ctx context.Context, job *domain.Job) error {
	return a.repo.Update(ctx, job)
}

func (a *jobApp) DeleteJob(ctx context.Context, jobID string, userID string) error {
	job, err := a.repo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}
	if job.CreatedBy.String() != userID {
		return domain.ErrUnauthorized
	}
	return a.repo.Delete(ctx, jobID)
}

func (a *jobApp) GetJobByID(ctx context.Context, jobID string) (*domain.Job, error) {
	return a.repo.FindByID(ctx, jobID)
}

func (a *jobApp) GetJobsByCompany(ctx context.Context, companyID string, page, size int) ([]domain.Job, int64, error) {
	return a.repo.FindByCompany(ctx, companyID, page, size)
}

func (a *jobApp) SearchJobs(ctx context.Context, filters map[string]interface{}, page, size int) ([]domain.Job, int64, error) {
	return a.repo.Search(ctx, filters, page, size)
}
