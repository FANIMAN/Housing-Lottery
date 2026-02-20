package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepo struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepo {
	return &DashboardRepo{db: db}
}

// Get summary with optional filters
func (r *DashboardRepo) GetSummary(
	ctx context.Context,
	subcityId, lotteryId, status, startDate, endDate string,
) (*interfaces.DashboardSummary, error) {

	var args []interface{}
	var conditions []string
	argID := 1

	baseQuery := `
	SELECT
		COUNT(DISTINCT s.id) AS subcities,
		COUNT(DISTINCT l.id) AS lotteries,
		COUNT(DISTINCT a.id) AS applicants,
		COUNT(DISTINCT w.id) AS winners
	FROM subcities s
	LEFT JOIN lotteries l ON l.subcity_id = s.id
	LEFT JOIN applicants a ON a.subcity_id = s.id
	LEFT JOIN lottery_winners w ON w.lottery_id = l.id
	`

	// Filters

	if subcityId != "" {
		conditions = append(conditions, fmt.Sprintf("s.id = $%d", argID))
		args = append(args, subcityId)
		argID++
	}

	if lotteryId != "" {
		conditions = append(conditions, fmt.Sprintf("l.id = $%d", argID))
		args = append(args, lotteryId)
		argID++
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("l.status = $%d", argID))
		args = append(args, status)
		argID++
	}

	if startDate != "" {
		conditions = append(conditions, fmt.Sprintf("l.created_at >= $%d", argID))
		args = append(args, startDate)
		argID++
	}

	if endDate != "" {
		conditions = append(conditions, fmt.Sprintf("l.created_at <= $%d", argID))
		args = append(args, endDate)
		argID++
	}

	query := baseQuery

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	row := r.db.QueryRow(ctx, query, args...)

	summary := &interfaces.DashboardSummary{}
	if err := row.Scan(
		&summary.Subcities,
		&summary.Lotteries,
		&summary.Applicants,
		&summary.Winners,
	); err != nil {
		return nil, err
	}

	return summary, nil
}
