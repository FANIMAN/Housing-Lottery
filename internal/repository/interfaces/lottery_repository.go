package interfaces

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
)

type LotteryRepository interface {
	Create(ctx context.Context, lottery *domain.Lottery) error
	InsertWinners(ctx context.Context, winners []domain.LotteryWinner) error
	GetByID(ctx context.Context, id string) (*domain.Lottery, error)
	ListAll(ctx context.Context) ([]*domain.Lottery, error)

}
