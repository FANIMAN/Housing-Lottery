package interfaces

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/google/uuid"
)

type ApplicantRepository interface {
	CreateBulk(ctx context.Context, applicants []*domain.Applicant) error
	GetBySubcityRegistrationIDs(ctx context.Context, subcityID uuid.UUID, registrationIDs []string) (map[string]bool, error)
    GetAllBySubcityID(ctx context.Context, subcityID uuid.UUID) ([]*domain.Applicant, error)

}
