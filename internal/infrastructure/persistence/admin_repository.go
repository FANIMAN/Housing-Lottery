package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *adminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(ctx context.Context, admin *domain.Admin) error {
	query := `
		INSERT INTO admins (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		admin.ID,
		admin.Email,
		admin.PasswordHash,
		admin.CreatedAt,
	)

	return err
}

func (r *adminRepository) GetByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM admins
		WHERE email = $1
	`

	var admin domain.Admin

	err := r.db.QueryRow(ctx, query, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

