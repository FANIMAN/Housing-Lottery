package domain

import (
	"time"

	"github.com/google/uuid"
)

type Applicant struct {
	ID                        uuid.UUID `json:"id"`
	FullName                  string    `json:"full_name"`
	CondominiumRegistrationID string    `json:"condominium_registration_id"`
	SubcityID                 uuid.UUID `json:"subcity_id"`
	UploadBatchID             uuid.UUID `json:"upload_batch_id"`
	CreatedAt                 time.Time `json:"created_at"`
}