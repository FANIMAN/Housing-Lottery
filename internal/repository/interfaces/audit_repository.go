package interfaces

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type AuditRepository interface {
	Log(
		ctx context.Context,
		adminID string,
		action string,
		entityType string,
		entityID string,
		httpStatus int,
		ipAddress string,
		userAgent string,
		errorMessage string,
	) error

	List(
		ctx context.Context,
		adminID, action, entityType string,
		fromDate, toDate *time.Time,
		limit, offset int,
	) ([]*domain.AuditLogResponse, int, error)
}