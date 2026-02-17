package interfaces

import "context"

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
}
