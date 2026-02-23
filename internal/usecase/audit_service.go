package usecase

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
)

type AuditService struct {
	repo interfaces.AuditRepository
}

func NewAuditService(repo interfaces.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) List(
	ctx context.Context,
	adminID, action, entityType string,
	fromDate, toDate *time.Time,
	page, pageSize int,
) ([]*domain.AuditLogResponse, int, error) {

	offset := (page - 1) * pageSize

	return s.repo.List(
		ctx,
		adminID,
		action,
		entityType,
		fromDate,
		toDate,
		pageSize,
		offset,
	)
}