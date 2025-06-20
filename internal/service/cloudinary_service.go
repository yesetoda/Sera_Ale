package service

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
)

type CloudinaryService interface {
	UploadPDF(ctx context.Context, file interface{}, publicID string) (string, error)
}

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return nil, err
	}
	return &cloudinaryService{cld: cld}, nil
}

func (s *cloudinaryService) UploadPDF(ctx context.Context, file interface{}, publicID string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "resumes",
		ResourceType: "raw",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
} 