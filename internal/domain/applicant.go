package domain

import "time"

type Applicant struct {
	ID                         string
	FullName                   string
	CondominiumRegistrationID  string
	SubcityID                  string
	UploadBatchID              string
	CreatedAt                  time.Time
}
