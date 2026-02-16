package usecase

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/google/uuid"
)

type SubcityUsecase struct {
	repo interfaces.SubcityRepository
}

func NewSubcityUsecase(r interfaces.SubcityRepository) *SubcityUsecase {
	return &SubcityUsecase{repo: r}
}

func (u *SubcityUsecase) Create(ctx context.Context, name string) error {
	subcity := &domain.Subcity{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}
	return u.repo.Create(ctx, subcity)
}

func (u *SubcityUsecase) GetByID(ctx context.Context, id string) (*domain.Subcity, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *SubcityUsecase) GetAll(ctx context.Context) ([]*domain.Subcity, error) {
	return u.repo.GetAll(ctx)
}

func (u *SubcityUsecase) Update(ctx context.Context, id, name string) error {
	subcity := &domain.Subcity{
		ID:   id,
		Name: name,
	}
	return u.repo.Update(ctx, subcity)
}

func (u *SubcityUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
