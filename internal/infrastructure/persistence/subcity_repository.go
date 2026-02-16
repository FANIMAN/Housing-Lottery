package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subcityRepository struct {
	db *pgxpool.Pool
}

func NewSubcityRepository(db *pgxpool.Pool) *subcityRepository {
	return &subcityRepository{db: db}
}

func (r *subcityRepository) Create(ctx context.Context, subcity *domain.Subcity) error {
	query := `
		INSERT INTO subcities (id, name, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, subcity.ID, subcity.Name, subcity.CreatedAt)
	return err
}

func (r *subcityRepository) GetByID(ctx context.Context, id string) (*domain.Subcity, error) {
	query := `
		SELECT id, name, created_at
		FROM subcities
		WHERE id = $1
	`
	var subcity domain.Subcity
	err := r.db.QueryRow(ctx, query, id).Scan(&subcity.ID, &subcity.Name, &subcity.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &subcity, nil
}

func (r *subcityRepository) GetAll(ctx context.Context) ([]*domain.Subcity, error) {
	query := `SELECT id, name, created_at FROM subcities ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subcities := []*domain.Subcity{}
	for rows.Next() {
		var s domain.Subcity
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		subcities = append(subcities, &s)
	}

	return subcities, nil
}

func (r *subcityRepository) Update(ctx context.Context, subcity *domain.Subcity) error {
	query := `UPDATE subcities SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, subcity.Name, subcity.ID)
	return err
}

func (r *subcityRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subcities WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
