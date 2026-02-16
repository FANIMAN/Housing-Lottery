package interfaces

import (
	"context"
	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type ApplicantRepository interface {
	BulkInsert(ctx context.Context, applicants []domain.Applicant) error
	GetBySubcity(ctx context.Context, subcityID string) ([]domain.Applicant, error)
}
