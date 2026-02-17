package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	repo      interfaces.AdminRepository
	auditRepo interfaces.AuditRepository
	jwtSecret string
}

func NewAdminUsecase(
	repo interfaces.AdminRepository,
	auditRepo interfaces.AuditRepository,
	jwtSecret string,
) *AdminUsecase {
	return &AdminUsecase{
		repo:      repo,
		auditRepo: auditRepo,
		jwtSecret: jwtSecret,
	}
}


func (u *AdminUsecase) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		_ = u.auditRepo.Log(ctx, "", "admin_register_failed", "admin", "", 0, "", "", err.Error())
		return err
	}

	admin := &domain.Admin{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := u.repo.Create(ctx, admin); err != nil {
		_ = u.auditRepo.Log(ctx, "", "admin_register_failed", "admin", admin.ID, 0, "", "", err.Error())
		return err
	}

	_ = u.auditRepo.Log(ctx, admin.ID, "admin_register", "admin", admin.ID, 0, "", "", "")
	return nil
}

// Login authenticates an admin and returns a JWT token
func (u *AdminUsecase) Login(ctx context.Context, email, password string) (string, error) {
	admin, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Create JWT token
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"email":    admin.Email,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		_ = u.auditRepo.Log(ctx, "", "admin_login_failed", "admin", admin.ID, 0, "", "", err.Error())
		return "", err
	}
	_ = u.auditRepo.Log(ctx, admin.ID, "admin_login", "admin", admin.ID, 0, "", "", "")
	return signedToken, nil
}
