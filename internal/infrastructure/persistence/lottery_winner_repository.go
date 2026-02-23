package persistence

import (
	"context"
	"fmt"
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


func (r *LotteryWinnerRepo) ListWinners(
	ctx context.Context,
	subcityName, fullName, lotteryName string,
	fromDate, toDate *time.Time,
	limit, offset int,
) ([]*domain.LotteryWinnerResponse, int, error) {

	baseQuery := `
		FROM lottery_winners w
		JOIN applicants a ON w.applicant_id = a.id
		JOIN lotteries l ON w.lottery_id = l.id
		JOIN subcities s ON a.subcity_id = s.id
		WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	// Filters
	if fullName != "" {
		baseQuery += fmt.Sprintf(" AND a.full_name ILIKE $%d", argID)
		args = append(args, "%"+fullName+"%")
		argID++
	}

	if lotteryName != "" {
		baseQuery += fmt.Sprintf(" AND l.name ILIKE $%d", argID)
		args = append(args, "%"+lotteryName+"%")
		argID++
	}

	if subcityName != "" {
		baseQuery += fmt.Sprintf(" AND s.name ILIKE $%d", argID)
		args = append(args, "%"+subcityName+"%")
		argID++
	}

	if fromDate != nil {
		baseQuery += fmt.Sprintf(" AND w.announced_at >= $%d", argID)
		args = append(args, *fromDate)
		argID++
	}

	if toDate != nil {
		baseQuery += fmt.Sprintf(" AND w.announced_at <= $%d", argID)
		args = append(args, *toDate)
		argID++
	}

	// ✅ 1. Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// ✅ 2. Get paginated data
	dataQuery := `
		SELECT 
			w.id,
			a.full_name,
			s.name AS subcity_name,
			l.name AS lottery_name,
			w.position_order,
			w.announced_at
	` + baseQuery +
		fmt.Sprintf(" ORDER BY w.announced_at DESC LIMIT $%d OFFSET $%d", argID, argID+1)

	argsWithPagination := append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, argsWithPagination...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	winners := []*domain.LotteryWinnerResponse{}

	for rows.Next() {
		var w domain.LotteryWinnerResponse
		if err := rows.Scan(
			&w.WinnerID,
			&w.FullName,
			&w.Subcity,
			&w.LotteryName,
			&w.Position,
			&w.AnnouncedAt,
		); err != nil {
			return nil, 0, err
		}
		winners = append(winners, &w)
	}

	return winners, total, nil
}