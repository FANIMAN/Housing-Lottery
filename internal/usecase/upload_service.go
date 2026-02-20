package usecase

import (
	"context"
	"io"
	"strings"
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


func (s *UploadService) GetApplicantsBySubcity(ctx context.Context, subcityIDStr string, search string) ([]*domain.Applicant, error) {
    subcityID, err := uuid.Parse(subcityIDStr)
    if err != nil {
        return nil, err
    }
    applicants, err := s.applicantRepo.GetAllBySubcityID(ctx, subcityID)
    if err != nil {
        return nil, err
    }
    if search != "" {
        filtered := []*domain.Applicant{}
        for _, a := range applicants {
            if strings.Contains(strings.ToLower(a.FullName), strings.ToLower(search)) ||
               strings.Contains(strings.ToLower(a.CondominiumRegistrationID), strings.ToLower(search)) {
                filtered = append(filtered, a)
            }
        }
        return filtered, nil
    }
    return applicants, nil
}

// func (s *UploadService) GetAllApplicants(ctx context.Context, search string) ([]*domain.Applicant, error) {
//     // This could list all applicants across subcities
//     // You can implement later if needed
//     return []*domain.Applicant{}, nil
// }

func (s *UploadService) GetApplicants(
	ctx context.Context,
	subcityIDStr string,
	search string,
	page int,
	limit int,
) ([]*domain.Applicant, int, error) {

	var subcityID *uuid.UUID
	if subcityIDStr != "" {
		parsed, err := uuid.Parse(subcityIDStr)
		if err != nil {
			return nil, 0, err
		}
		subcityID = &parsed
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return s.applicantRepo.GetFiltered(ctx, subcityID, search, limit, offset)
}