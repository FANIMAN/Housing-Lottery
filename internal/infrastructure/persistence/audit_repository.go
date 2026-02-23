package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Log(
	ctx context.Context,
	adminID string,
	action string,
	entityType string,
	entityID string,
	httpStatus int,
	ipAddress string,
	userAgent string,
	errorMessage string,
) error {

	query := `
		INSERT INTO audit_logs
		(admin_id, action, entity_type, entity_id, http_status, ip_address, user_agent, error_message)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		nullIfEmpty(adminID),
		action,
		entityType,
		nullIfEmpty(entityID),
		nullIfZero(httpStatus),
		nullIfEmpty(ipAddress),
		nullIfEmpty(userAgent),
		nullIfEmpty(errorMessage),
	)
	return err
}

func nullIfEmpty(val string) interface{} {
	if val == "" {
		return nil
	}
	return val
}

func nullIfZero(val int) interface{} {
	if val == 0 {
		return nil
	}
	return val
}

func (r *AuditRepo) List(
	ctx context.Context,
	adminID, action, entityType string,
	fromDate, toDate *time.Time,
	limit, offset int,
) ([]*domain.AuditLogResponse, int, error) {

	query := `
	SELECT 
		a.id,
		a.admin_id,
		ad.email,
		a.action,
		a.entity_type,
		a.entity_id,
		a.http_status,
		a.ip_address,
		a.user_agent,
		a.error_message,
		a.created_at
	FROM audit_logs a
	LEFT JOIN admins ad ON a.admin_id = ad.id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM audit_logs a
	WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if adminID != "" {
		query += fmt.Sprintf(" AND a.admin_id = $%d", argID)
		countQuery += fmt.Sprintf(" AND a.admin_id = $%d", argID)
		args = append(args, adminID)
		argID++
	}

	if action != "" {
		query += fmt.Sprintf(" AND a.action ILIKE $%d", argID)
		countQuery += fmt.Sprintf(" AND a.action ILIKE $%d", argID)
		args = append(args, "%"+action+"%")
		argID++
	}

	if entityType != "" {
		query += fmt.Sprintf(" AND a.entity_type = $%d", argID)
		countQuery += fmt.Sprintf(" AND a.entity_type = $%d", argID)
		args = append(args, entityType)
		argID++
	}

	if fromDate != nil {
		query += fmt.Sprintf(" AND a.created_at >= $%d", argID)
		countQuery += fmt.Sprintf(" AND a.created_at >= $%d", argID)
		args = append(args, *fromDate)
		argID++
	}

	if toDate != nil {
		query += fmt.Sprintf(" AND a.created_at <= $%d", argID)
		countQuery += fmt.Sprintf(" AND a.created_at <= $%d", argID)
		args = append(args, *toDate)
		argID++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var audits []*domain.AuditLogResponse

	for rows.Next() {
		var a domain.AuditLogResponse

		err := rows.Scan(
			&a.ID,
			&a.AdminID,
			&a.AdminEmail,
			&a.Action,
			&a.EntityType,
			&a.EntityID,
			&a.HTTPStatus,
			&a.IPAddress,
			&a.UserAgent,
			&a.ErrorMessage,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		audits = append(audits, &a)
	}

	return audits, total, nil
}
