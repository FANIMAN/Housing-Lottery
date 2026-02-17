package domain

import "time"
import "github.com/google/uuid"

type UploadBatch struct {
	ID           uuid.UUID `json:"id"`
	SubcityID    uuid.UUID `json:"subcity_id"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	FileName     string    `json:"file_name"`
	TotalRecords int       `json:"total_records"`
	CreatedAt    time.Time `json:"created_at"`
}