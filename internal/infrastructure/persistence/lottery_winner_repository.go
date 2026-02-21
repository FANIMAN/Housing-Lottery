package persistence

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LotteryWinnerRepo struct {
	db *pgxpool.Pool
}

func NewLotteryWinnerRepository(db *pgxpool.Pool) *LotteryWinnerRepo {
	return &LotteryWinnerRepo{db: db}
}

func (r *LotteryWinnerRepo) Create(ctx context.Context, winner *domain.LotteryWinner) error {
	if winner.ID == "" {
		winner.ID = uuid.New().String()
	}
	if winner.AnnouncedAt == nil {
		now := time.Now()
		winner.AnnouncedAt = &now
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO lottery_winners (id, lottery_id, applicant_id, position_order, announced_at)
		VALUES ($1, $2, $3, $4, $5)
	`, winner.ID, winner.LotteryID, winner.ApplicantID, winner.PositionOrder, winner.AnnouncedAt)
	return err
}

func (r *LotteryWinnerRepo) CreateTx(ctx context.Context, tx pgx.Tx, winner *domain.LotteryWinner) error {
	if winner.ID == "" {
		winner.ID = uuid.New().String()
	}
	if winner.AnnouncedAt == nil {
		now := time.Now()
		winner.AnnouncedAt = &now
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO lottery_winners (id, lottery_id, applicant_id, position_order, announced_at)
		VALUES ($1, $2, $3, $4, $5)
	`, winner.ID, winner.LotteryID, winner.ApplicantID, winner.PositionOrder, winner.AnnouncedAt)
	return err
}

func (r *LotteryWinnerRepo) GetWinnersByLottery(ctx context.Context, lotteryID string) ([]*domain.LotteryWinner, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, lottery_id, applicant_id, position_order, announced_at
		FROM lottery_winners
		WHERE lottery_id=$1
		ORDER BY position_order ASC
	`, lotteryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var winners []*domain.LotteryWinner
	for rows.Next() {
		var w domain.LotteryWinner
		if err := rows.Scan(&w.ID, &w.LotteryID, &w.ApplicantID, &w.PositionOrder, &w.AnnouncedAt); err != nil {
			return nil, err
		}
		winners = append(winners, &w)
	}
	return winners, nil
}
