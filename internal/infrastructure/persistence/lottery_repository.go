package persistence

import (
	"context"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LotteryRepo struct {
	db *pgxpool.Pool
}

func NewLotteryRepository(db *pgxpool.Pool) *LotteryRepo {
	return &LotteryRepo{db: db}
}

// Create or update lottery
func (r *LotteryRepo) Create(ctx context.Context, lottery *domain.Lottery) error {
	if lottery.ID == "" {
		lottery.ID = uuid.New().String()
	}
	if lottery.CreatedAt.IsZero() {
		lottery.CreatedAt = time.Now()
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO lotteries (id, subcity_id, total_applicants, winners_count, seed_value, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE
		SET status = EXCLUDED.status,
		    total_applicants = EXCLUDED.total_applicants,
		    winners_count = EXCLUDED.winners_count
	`, lottery.ID, lottery.SubcityID, lottery.TotalApplicants, lottery.WinnersCount, lottery.SeedValue, lottery.Status, lottery.CreatedAt)

	return err
}

// Get a lottery by ID
func (r *LotteryRepo) GetByID(ctx context.Context, id string) (*domain.Lottery, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, subcity_id, total_applicants, winners_count, seed_value, status, created_at
		FROM lotteries
		WHERE id = $1
	`, id)

	var l domain.Lottery
	if err := row.Scan(&l.ID, &l.SubcityID, &l.TotalApplicants, &l.WinnersCount, &l.SeedValue, &l.Status, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

// Insert multiple winners for a lottery
func (r *LotteryRepo) InsertWinners(ctx context.Context, winners []domain.LotteryWinner) error {
	batch := &pgx.Batch{}
	for _, w := range winners {
		if w.ID == "" {
			w.ID = uuid.New().String()
		}
		if w.AnnouncedAt == nil {
			now := time.Now()
			w.AnnouncedAt = &now
		}
		batch.Queue(`
			INSERT INTO lottery_winners (id, lottery_id, applicant_id, position_order, announced_at)
			VALUES ($1,$2,$3,$4,$5)
		`, w.ID, w.LotteryID, w.ApplicantID, w.PositionOrder, w.AnnouncedAt)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for range winners {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}
