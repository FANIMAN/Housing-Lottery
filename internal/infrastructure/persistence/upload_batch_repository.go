package persistence

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type uploadBatchRepository struct {
	db *pgxpool.Pool
}

func NewUploadBatchRepository(db *pgxpool.Pool) *uploadBatchRepository {
	return &uploadBatchRepository{db: db}
}

func (r *uploadBatchRepository) Create(ctx context.Context, batch *domain.UploadBatch) error {
	query := `
		INSERT INTO upload_batches (id, subcity_id, uploaded_by, file_name, total_records, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		batch.ID,
		batch.SubcityID,
		batch.UploadedBy,
		batch.FileName,
		batch.TotalRecords,
		batch.CreatedAt,
	)
	return err
}



