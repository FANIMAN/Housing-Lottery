package domain

import "time"

type LotteryStatus string

const (
	LotteryPending   LotteryStatus = "pending"
	LotteryCompleted LotteryStatus = "completed"
	LotteryCancelled LotteryStatus = "cancelled"
)

type Lottery struct {
	ID              string
	SubcityID       string
	TotalApplicants int
	WinnersCount    int
	SeedValue       int64
	Status          LotteryStatus
	CreatedAt       time.Time
}
