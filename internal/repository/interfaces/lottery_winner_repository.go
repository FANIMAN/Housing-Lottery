package interfaces

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5"
)

type LotteryWinnerRepository interface {
	Create(ctx context.Context, winner *domain.LotteryWinner) error
	CreateTx(ctx context.Context, tx pgx.Tx, winner *domain.LotteryWinner) error
	GetWinnersByLottery(ctx context.Context, lotteryID string) ([]*domain.LotteryWinner, error)

	ListWinners(
		ctx context.Context,
		subcityName, fullName, lotteryName string,
		fromDate, toDate *time.Time,
		limit, offset int,
	) ([]*domain.LotteryWinnerResponse, int, error)
}
