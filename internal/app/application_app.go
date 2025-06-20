package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yesetoda/Sera_Ale/internal/domain"
	"github.com/yesetoda/Sera_Ale/internal/repository"
	"github.com/yesetoda/Sera_Ale/internal/service"
)

type ApplicationApp interface {
	Apply(ctx context.Context, applicantID, jobID, coverLetter string, resumeFile interface{}) (*domain.Application, []string)
	TrackApplications(ctx context.Context, applicantID string, page, size int) ([]domain.Application, int64, error)
	GetApplicationsForJob(ctx context.Context, jobID, companyID string, page, size int) ([]domain.Application, int64, error)
	UpdateStatus(ctx context.Context, applicationID, companyID, status string) (*domain.Application, error)
}

type applicationApp struct {
	repo    repository.ApplicationRepository
	jobRepo repository.JobRepository
	cloud   service.CloudinaryService
}

func NewApplicationApp(repo repository.ApplicationRepository, jobRepo repository.JobRepository, cloud service.CloudinaryService) ApplicationApp {
	return &applicationApp{repo: repo, jobRepo: jobRepo, cloud: cloud}
}

func (a *applicationApp) Apply(ctx context.Context, applicantID, jobID, coverLetter string, resumeFile interface{}) (*domain.Application, []string) {
	if len(coverLetter) > 200 {
		return nil, []string{"Cover letter must be under 200 characters"}
	}
	if _, err := a.repo.FindByApplicantAndJob(ctx, applicantID, jobID); err == nil {
		return nil, []string{"You have already applied to this job"}
	}
	publicID := uuid.New().String()
	resumeURL, err := a.cloud.UploadPDF(ctx, resumeFile, publicID)
	if err != nil {
		return nil, []string{"Failed to upload resume"}
	}
	app := &domain.Application{
		ID:          uuid.New(),
		ApplicantID: uuid.MustParse(applicantID),
		JobID:       uuid.MustParse(jobID),
		ResumeLink:  resumeURL,
		CoverLetter: coverLetter,
		Status:      domain.StatusApplied,
	}
	if err := a.repo.Create(ctx, app); err != nil {
		return nil, []string{"Failed to create application"}
	}
	return app, nil
}

func (a *applicationApp) TrackApplications(ctx context.Context, applicantID string, page, size int) ([]domain.Application, int64, error) {
	return a.repo.FindByApplicant(ctx, applicantID, page, size)
}

func (a *applicationApp) GetApplicationsForJob(ctx context.Context, jobID, companyID string, page, size int) ([]domain.Application, int64, error) {
	job, err := a.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, 0, errors.New("Job not found")
	}
	if job.CreatedBy.String() != companyID {
		return nil, 0, errors.New("Unauthorized access")
	}
	return a.repo.FindByJob(ctx, jobID, page, size)
}

func (a *applicationApp) UpdateStatus(ctx context.Context, applicationID, companyID, status string) (*domain.Application, error) {
	app, err := a.repo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, errors.New("Application not found")
	}
	job, err := a.jobRepo.FindByID(ctx, app.JobID.String())
	if err != nil {
		return nil, errors.New("Job not found")
	}
	if job.CreatedBy.String() != companyID {
		return nil, errors.New("Unauthorized")
	}
	if err := a.repo.UpdateStatus(ctx, applicationID, domain.ApplicationStatus(status)); err != nil {
		return nil, errors.New("Failed to update status")
	}
	return a.repo.FindByID(ctx, applicationID)
}
