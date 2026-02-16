package usecase

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	repo interfaces.AdminRepository
}

func NewAdminUsecase(r interfaces.AdminRepository) *AdminUsecase {
	return &AdminUsecase{repo: r}
}

func (u *AdminUsecase) Register(ctx context.Context, email, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &domain.Admin{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	return u.repo.Create(ctx, admin)
}
