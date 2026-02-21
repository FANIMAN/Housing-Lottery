package interfaces

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/jackc/pgx/v5"
)

type LotteryRepository interface {
	Create(ctx context.Context, lottery *domain.Lottery) error
	InsertWinners(ctx context.Context, winners []domain.LotteryWinner) error
	GetByID(ctx context.Context, id string) (*domain.Lottery, error)
	ListAll(ctx context.Context) ([]*domain.Lottery, error)
	IncrementWinnersCount(ctx context.Context, tx pgx.Tx, lotteryID string) error 
	UpdateStatus(ctx context.Context, tx pgx.Tx, lotteryID string, status domain.LotteryStatus) error 

}
