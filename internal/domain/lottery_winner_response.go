package domain

import "time"

type LotteryWinnerResult struct {
	Winner    *LotteryWinner
	Applicant *Applicant
}


type LotteryWinnerResponse struct {
	WinnerID    string    `json:"winner_id"`
	FullName    string    `json:"full_name"`
	Subcity     string    `json:"subcity"`
	LotteryName string    `json:"lottery_name"`
	Position    int       `json:"position"`
	AnnouncedAt time.Time `json:"announced_at"`
}