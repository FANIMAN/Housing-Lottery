package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type applicantRepository struct {
	db *pgxpool.Pool
}

func NewApplicantRepository(db *pgxpool.Pool) *applicantRepository {
	return &applicantRepository{db: db}
}

// Bulk insert applicants
func (r *applicantRepository) CreateBulk(ctx context.Context, applicants []*domain.Applicant) error {
	if len(applicants) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	stmt := `INSERT INTO applicants 
		(id, full_name, condominium_registration_id, subcity_id, upload_batch_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`

	for _, a := range applicants {
		_, err := tx.Exec(ctx, stmt,
			a.ID,
			a.FullName,
			a.CondominiumRegistrationID,
			a.SubcityID,
			a.UploadBatchID,
			a.CreatedAt,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

// Fetch existing registration IDs in the subcity to skip duplicates
func (r *applicantRepository) GetBySubcityRegistrationIDs(ctx context.Context, subcityID uuid.UUID, registrationIDs []string) (map[string]bool, error) {
	if len(registrationIDs) == 0 {
		return map[string]bool{}, nil
	}

	query := `SELECT condominium_registration_id FROM applicants WHERE subcity_id = $1 AND condominium_registration_id = ANY($2)`
	rows, err := r.db.Query(ctx, query, subcityID, registrationIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existing[id] = true
	}

	return existing, nil
}

// Fetch all applicant registration IDs for a given subcity
func (r *applicantRepository) GetAllBySubcityID(ctx context.Context, subcityID uuid.UUID) ([]*domain.Applicant, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, full_name, condominium_registration_id, subcity_id, upload_batch_id, created_at
		FROM applicants
		WHERE subcity_id=$1
	`, subcityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applicants []*domain.Applicant
	for rows.Next() {
		var a domain.Applicant
		if err := rows.Scan(&a.ID, &a.FullName, &a.CondominiumRegistrationID, &a.SubcityID, &a.UploadBatchID, &a.CreatedAt); err != nil {
			return nil, err
		}
		applicants = append(applicants, &a)
	}
	return applicants, nil
}

