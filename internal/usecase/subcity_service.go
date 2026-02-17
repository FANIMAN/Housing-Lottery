package usecase

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/google/uuid"
)


type SubcityUsecase struct {
	repo      interfaces.SubcityRepository
	auditRepo interfaces.AuditRepository
}

func NewSubcityUsecase(r interfaces.SubcityRepository, a interfaces.AuditRepository) *SubcityUsecase {
	return &SubcityUsecase{repo: r, auditRepo: a}
}

func (u *SubcityUsecase) Create(ctx context.Context, name string) error {
	subcity := &domain.Subcity{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, subcity); err != nil {
		_ = u.auditRepo.Log(ctx, "", "subcity_create_failed", "subcity", "", 0, "", "", err.Error())
		return err
	}

	_ = u.auditRepo.Log(ctx, "", "subcity_create", "subcity", subcity.ID, 0, "", "", "")
	return nil
}

func (u *SubcityUsecase) Update(ctx context.Context, id, name string) error {
	subcity := &domain.Subcity{ID: id, Name: name}
	if err := u.repo.Update(ctx, subcity); err != nil {
		_ = u.auditRepo.Log(ctx, "", "subcity_update_failed", "subcity", id, 0, "", "", err.Error())
		return err
	}
	_ = u.auditRepo.Log(ctx, "", "subcity_update", "subcity", id, 0, "", "", "")
	return nil
}

func (u *SubcityUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		_ = u.auditRepo.Log(ctx, "", "subcity_delete_failed", "subcity", id, 0, "", "", err.Error())
		return err
	}
	_ = u.auditRepo.Log(ctx, "", "subcity_delete", "subcity", id, 0, "", "", "")
	return nil
}


func (u *SubcityUsecase) GetByID(ctx context.Context, id string) (*domain.Subcity, error) {
	_ = u.auditRepo.Log(ctx, "", "subcity_get_by_id", "subcity", id, 0, "", "", "")
	return u.repo.GetByID(ctx, id)
}

func (u *SubcityUsecase) GetAll(ctx context.Context) ([]*domain.Subcity, error) {
	_ = u.auditRepo.Log(ctx, "", "subcity_get_all", "subcity", "", 0, "", "", "")
	return u.repo.GetAll(ctx)
}


