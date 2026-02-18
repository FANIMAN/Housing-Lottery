package interfaces

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type SubcityRepository interface {
	Create(ctx context.Context, subcity *domain.Subcity) error
	GetByID(ctx context.Context, id string) (*domain.Subcity, error)
	GetAll(ctx context.Context) ([]*domain.Subcity, error)
	Update(ctx context.Context, subcity *domain.Subcity) error
	Delete(ctx context.Context, id string) error
	ListAll(ctx context.Context) ([]*domain.Subcity, error)

}
