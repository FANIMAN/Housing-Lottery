package persistence

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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



func (r *applicantRepository) GetFiltered(ctx context.Context, subcityID *uuid.UUID, search string, limit int, offset int) ([]*domain.Applicant, int, error) {
	// Build base query
	query := `SELECT id, full_name, condominium_registration_id, subcity_id, upload_batch_id, created_at
	          FROM applicants WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if subcityID != nil {
		query += ` AND subcity_id=$` + strconv.Itoa(argIdx)
		args = append(args, *subcityID)
		argIdx++
	}

	if search != "" {
		query += ` AND (LOWER(full_name) LIKE $` + strconv.Itoa(argIdx) + ` OR LOWER(condominium_registration_id) LIKE $` + strconv.Itoa(argIdx+1) + `)`
		args = append(args, "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
		argIdx += 2
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS subquery"
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	applicants := []*domain.Applicant{}
	for rows.Next() {
		var a domain.Applicant
		if err := rows.Scan(&a.ID, &a.FullName, &a.CondominiumRegistrationID, &a.SubcityID, &a.UploadBatchID, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		applicants = append(applicants, &a)
	}

	return applicants, total, nil
}