package interfaces

import (
	"context"
	"time"
)

type DashboardSummary struct {
	Subcities  int64 `json:"subcities"`
	Lotteries  int64 `json:"lotteries"`
	Applicants int64 `json:"applicants"`
	Winners    int64 `json:"winners"`
}

type DashboardRepository interface {
	GetSummary(ctx context.Context, from, to *time.Time, subcityID, lotteryStatus *string) (*DashboardSummary, error)
}
