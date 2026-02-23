package interfaces

import (
	"context"
	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *domain.Admin) error
	GetByEmail(ctx context.Context, email string) (*domain.Admin, error)
	List(ctx context.Context, email, id string, limit, offset int) ([]*domain.Admin, int, error)

}
