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
func (r *DashboardRepo) GetSummary(ctx context.Context, subcityId, status, startDate, endDate string) (*interfaces.DashboardSummary, error) {
	var args []interface{}
	var conditions []string
	argID := 1

	if subcityId != "" {
		conditions = append(conditions, fmt.Sprintf("subcity_id = $%d", argID))
		args = append(args, subcityId)
		argID++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, status)
		argID++
	}
	if startDate != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argID))
		args = append(args, startDate)
		argID++
	}
	if endDate != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argID))
		args = append(args, endDate)
		argID++
	}

	query := "SELECT COUNT(DISTINCT subcity_id) as subcities, COUNT(*) as lotteries, " +
		"COALESCE(SUM(total_applicants),0) as applicants, COALESCE(SUM(winners_count),0) as winners FROM lotteries"

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	row := r.db.QueryRow(ctx, query, args...)

	summary := &interfaces.DashboardSummary{}
	if err := row.Scan(&summary.Subcities, &summary.Lotteries, &summary.Applicants, &summary.Winners); err != nil {
		return nil, err
	}

	return summary, nil
}
