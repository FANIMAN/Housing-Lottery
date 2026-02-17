package interfaces

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type LotteryWinnerRepository interface {
	Create(ctx context.Context, winner *domain.LotteryWinner) error
	GetWinnersByLottery(ctx context.Context, lotteryID string) ([]*domain.LotteryWinner, error)
}
