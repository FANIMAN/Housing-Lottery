package persistence

import (
	"context"
	"fmt"

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



func (r *adminRepository) List(ctx context.Context, email, id string, limit, offset int) ([]*domain.Admin, int, error) {
	admins := []*domain.Admin{}
	var args []interface{}
	query := "SELECT id, email, password_hash, created_at FROM admins WHERE 1=1"

	// Filtering
	if email != "" {
		args = append(args, "%"+email+"%")
		query += " AND email ILIKE $" + fmt.Sprint(len(args))
	}
	if id != "" {
		args = append(args, "%"+id+"%")
		query += " AND id ILIKE $" + fmt.Sprint(len(args))
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS sub"
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	args = append(args, limit, offset)
	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)-1) + " OFFSET $" + fmt.Sprint(len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var admin domain.Admin
		if err := rows.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt); err != nil {
			return nil, 0, err
		}
		admins = append(admins, &admin)
	}

	return admins, total, nil
}