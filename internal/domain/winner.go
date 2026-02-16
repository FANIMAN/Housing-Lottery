package domain

import "time"

type LotteryWinner struct {
	ID            string
	LotteryID     string
	ApplicantID   string
	PositionOrder int
	AnnouncedAt   *time.Time
}
