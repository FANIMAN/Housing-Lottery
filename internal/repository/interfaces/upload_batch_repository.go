package interfaces

import (
	"context"
	// "github.com/google/uuid"
	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type UploadBatchRepository interface {
	Create(ctx context.Context, batch *domain.UploadBatch) error
	// UpdateTotalRecords(ctx context.Context, batchID uuid.UUID, total int) error
}


