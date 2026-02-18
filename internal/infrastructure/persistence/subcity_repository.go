package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubcityRepo struct {
	db *pgxpool.Pool
}

func NewSubcityRepository(db *pgxpool.Pool) *SubcityRepo {
	return &SubcityRepo{db: db}
}

func (r *SubcityRepo) Create(ctx context.Context, subcity *domain.Subcity) error {
	query := `
		INSERT INTO subcities (id, name, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, subcity.ID, subcity.Name, subcity.CreatedAt)
	return err
}

func (r *SubcityRepo) GetByID(ctx context.Context, id string) (*domain.Subcity, error) {
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

func (r *SubcityRepo) GetAll(ctx context.Context) ([]*domain.Subcity, error) {
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

func (r *SubcityRepo) Update(ctx context.Context, subcity *domain.Subcity) error {
	query := `UPDATE subcities SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, subcity.Name, subcity.ID)
	return err
}

func (r *SubcityRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subcities WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}


// ListAll returns all subcities
func (r *SubcityRepo) ListAll(ctx context.Context) ([]*domain.Subcity, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, created_at FROM subcities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcities []*domain.Subcity
	for rows.Next() {
		var s domain.Subcity
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		subcities = append(subcities, &s)
	}

	return subcities, nil
}
