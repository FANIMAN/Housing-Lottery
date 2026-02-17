package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) interfaces.AuditRepository {
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
