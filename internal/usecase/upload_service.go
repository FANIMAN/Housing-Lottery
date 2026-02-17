package usecase

import (
	"context"
	"io"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/FANIMAN/housing-lottery/internal/utils"
	"github.com/google/uuid"
)

type UploadService struct {
	applicantRepo interfaces.ApplicantRepository
	batchRepo     interfaces.UploadBatchRepository
	auditRepo     interfaces.AuditRepository
}

func NewUploadService(applicantRepo interfaces.ApplicantRepository, batchRepo interfaces.UploadBatchRepository, auditRepo interfaces.AuditRepository) *UploadService {
	return &UploadService{
		applicantRepo: applicantRepo,
		batchRepo:     batchRepo,
		auditRepo:     auditRepo,
	}
}

func (s *UploadService) ProcessExcel(ctx context.Context, subcityIDStr, adminIDStr string, file io.Reader, fileName string) (inserted int, skipped int, err error) {
	subcityID, err := uuid.Parse(subcityIDStr)
	if err != nil {
		_ = s.auditRepo.Log(ctx, adminIDStr, "upload_failed", "upload_batch", "", 0, "", "", err.Error())
		return 0, 0, err
	}
	adminID, _ := uuid.Parse(adminIDStr)

	rows, err := utils.ParseApplicantExcel(file)
	if err != nil {
		_ = s.auditRepo.Log(ctx, adminIDStr, "upload_failed", "upload_batch", "", 0, "", "", err.Error())
		return 0, 0, err
	}

	// Filter duplicates
	var applicants []*domain.Applicant
	for _, r := range rows {
		existing, _ := s.applicantRepo.GetBySubcityRegistrationIDs(ctx, subcityID, []string{r.CondominiumRegistrationID})
		if len(existing) > 0 {
			skipped++
			continue
		}
		applicants = append(applicants, &domain.Applicant{
			ID:                        uuid.New(),
			FullName:                  r.FullName,
			CondominiumRegistrationID: r.CondominiumRegistrationID,
			SubcityID:                 subcityID,
			CreatedAt:                 time.Now(),
		})
	}

	batch := &domain.UploadBatch{
		ID:           uuid.New(),
		SubcityID:    subcityID,
		UploadedBy:   adminID,
		FileName:     fileName,
		TotalRecords: len(applicants),
		CreatedAt:    time.Now(),
	}

	if err := s.batchRepo.Create(ctx, batch); err != nil {
		_ = s.auditRepo.Log(ctx, adminIDStr, "upload_failed", "upload_batch", batch.ID.String(), 0, "", "", err.Error())
		return 0, 0, err
	}

	for _, a := range applicants {
		a.UploadBatchID = batch.ID
	}

	if err := s.applicantRepo.CreateBulk(ctx, applicants); err != nil {
		_ = s.auditRepo.Log(ctx, adminIDStr, "upload_failed", "upload_batch", batch.ID.String(), 0, "", "", err.Error())
		return 0, 0, err
	}

	_ = s.auditRepo.Log(ctx, adminIDStr, "upload_success", "upload_batch", batch.ID.String(), 0, "", "", "")
	return len(applicants), skipped, nil
}
