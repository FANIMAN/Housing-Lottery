package domain

import (
	"time"

	"github.com/google/uuid"
)

type LotteryStatus string

const (
	LotteryPending   LotteryStatus = "pending"
	LotteryCompleted LotteryStatus = "completed"
	LotteryCancelled LotteryStatus = "cancelled"
)

type Lottery struct {
	ID              string
	SubcityID       uuid.UUID
	TotalApplicants int
	WinnersCount    int
	SeedValue       int64
	Status          LotteryStatus
	CreatedAt       time.Time
}
