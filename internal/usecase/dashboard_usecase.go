package usecase

import (
	"context"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/infrastructure/persistence"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
)

type DashboardUsecase struct {
	repo        *persistence.DashboardRepo
	subcityRepo *persistence.SubcityRepo // exported type
	lotteryRepo *persistence.LotteryRepo
}

// Constructor
func NewDashboardUsecase(
	dashboardRepo *persistence.DashboardRepo,
	subcityRepo *persistence.SubcityRepo,
	lotteryRepo *persistence.LotteryRepo,
) *DashboardUsecase {
	return &DashboardUsecase{
		repo:        dashboardRepo,
		subcityRepo: subcityRepo,
		lotteryRepo: lotteryRepo,
	}
}

// Get summary with optional filters
func (u *DashboardUsecase) GetSummary(subcityId, status, startDate, endDate string) (*interfaces.DashboardSummary, error) {
	return u.repo.GetSummary(context.Background(), subcityId, status, startDate, endDate)
}

// ListSubcities returns all subcities
func (u *DashboardUsecase) ListSubcities() ([]*domain.Subcity, error) {
	return u.subcityRepo.ListAll(context.Background())
}

// ListLotteries returns all lotteries
func (u *DashboardUsecase) ListLotteries() ([]*domain.Lottery, error) {
	return u.lotteryRepo.ListAll(context.Background())
}
