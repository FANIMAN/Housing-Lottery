package usecase

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/FANIMAN/housing-lottery/internal/utils"
	"github.com/google/uuid"
)

// ✅ Make struct exported
type UploadService struct {
	applicantRepo interfaces.ApplicantRepository
	batchRepo     interfaces.UploadBatchRepository
}

// ✅ Make constructor exported
func NewUploadService(applicantRepo interfaces.ApplicantRepository, batchRepo interfaces.UploadBatchRepository) *UploadService {
	return &UploadService{
		applicantRepo: applicantRepo,
		batchRepo:     batchRepo,
	}
}

// ProcessExcel method
func (s *UploadService) ProcessExcel(ctx context.Context, subcityIDStr, adminIDStr string, file io.Reader, fileName string) (inserted int, skipped int, err error) {
	log.Println("Start processing Excel file:", fileName)

    subcityID, err := uuid.Parse(subcityIDStr)
    if err != nil {
        log.Println("Error parsing subcity ID:", err)
        return 0, 0, err
    }
    adminID, err := uuid.Parse(adminIDStr)
    if err != nil {
        log.Println("Error parsing admin ID:", err)
        return 0, 0, err
    }

    log.Println("Subcity ID:", subcityID, "Admin ID:", adminID)

    // Parse the Excel file
    rows, err := utils.ParseApplicantExcel(file)
    if err != nil {
        log.Println("Error parsing Excel file:", err)
        return 0, 0, err
    }
    log.Println("Parsed rows:", len(rows))

	// Check duplicates in DB
	regIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		regIDs = append(regIDs, r.CondominiumRegistrationID)
	}

	existingMap, err := s.applicantRepo.GetBySubcityRegistrationIDs(ctx, subcityID, regIDs)
	if err != nil {
		return 0, 0, err
	}

	// Filter out duplicates
	var applicants []*domain.Applicant
	for _, r := range rows {
		if _, ok := existingMap[r.CondominiumRegistrationID]; ok {
			skipped++
			continue
		}

		applicants = append(applicants, &domain.Applicant{
			ID:                        uuid.New(), // type uuid.UUID
			FullName:                  r.FullName,
			CondominiumRegistrationID: r.CondominiumRegistrationID,
			SubcityID:                 subcityID,
			CreatedAt:                 time.Now(),
		})
	}

	// Create batch record
	batch := &domain.UploadBatch{
		ID:           uuid.New(),
		SubcityID:    subcityID,
		UploadedBy:   adminID,
		FileName:     fileName,
		TotalRecords: len(applicants),
		CreatedAt:    time.Now(),
	}

	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return 0, 0, err
	}

	// Assign batch ID to applicants
	for _, a := range applicants {
		a.UploadBatchID = batch.ID
	}

	// Bulk insert applicants
	if err := s.applicantRepo.CreateBulk(ctx, applicants); err != nil {
		return 0, 0, err
	}

	return len(applicants), skipped, nil
}
