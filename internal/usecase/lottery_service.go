package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/google/uuid"
)

type LotteryService struct {
	lotteryRepo       interfaces.LotteryRepository
	applicantRepo     interfaces.ApplicantRepository
	lotteryWinnerRepo interfaces.LotteryWinnerRepository
	auditRepo         interfaces.AuditRepository
}

func NewLotteryService(
	lotteryRepo interfaces.LotteryRepository,
	applicantRepo interfaces.ApplicantRepository,
	lotteryWinnerRepo interfaces.LotteryWinnerRepository,
	auditRepo interfaces.AuditRepository,
) *LotteryService {
	return &LotteryService{
		lotteryRepo:       lotteryRepo,
		applicantRepo:     applicantRepo,
		lotteryWinnerRepo: lotteryWinnerRepo,
		auditRepo:         auditRepo,
	}
}

func (s *LotteryService) StartLottery(ctx context.Context, subcityID uuid.UUID, adminID string) (*domain.Lottery, error) {
	applicants, err := s.applicantRepo.GetAllBySubcityID(ctx, subcityID)
	if err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "start_lottery_failed", "lottery", "", 0, "", "", err.Error())
		return nil, err
	}

	lottery := &domain.Lottery{
		ID:              uuid.NewString(),
		SubcityID:       subcityID,
		Status:          domain.LotteryPending,
		SeedValue:       time.Now().UnixNano(),
		CreatedAt:       time.Now(),
		TotalApplicants: len(applicants),
	}

	if err := s.lotteryRepo.Create(ctx, lottery); err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "start_lottery_failed", "lottery", lottery.ID, 0, "", "", err.Error())
		return nil, err
	}

	_ = s.auditRepo.Log(ctx, adminID, "start_lottery", "lottery", lottery.ID, 0, "", "", "")
	return lottery, nil
}

func (s *LotteryService) SpinLottery(ctx context.Context, lotteryID, adminID string) (*domain.LotteryWinner, error) {
	lottery, err := s.lotteryRepo.GetByID(ctx, lotteryID)
	if err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "lottery_spin_failed", "lottery", lotteryID, 0, "", "", err.Error())
		return nil, err
	}

	if lottery.Status != domain.LotteryPending {
		return nil, errors.New("lottery is not active")
	}

	applicants, err := s.applicantRepo.GetAllBySubcityID(ctx, lottery.SubcityID)
	if err != nil {
		return nil, err
	}

	winners, err := s.lotteryWinnerRepo.GetWinnersByLottery(ctx, lotteryID)
	if err != nil {
		return nil, err
	}

	wonIDs := make(map[string]bool)
	for _, w := range winners {
		wonIDs[w.ApplicantID] = true
	}

	var eligible []string
	for _, a := range applicants {
		aID := a.ID.String()
		if !wonIDs[aID] {
			eligible = append(eligible, aID)
		}
	}

	if len(eligible) == 0 {
		lottery.Status = domain.LotteryCompleted
		_ = s.lotteryRepo.Create(ctx, lottery)
		_ = s.auditRepo.Log(ctx, adminID, "lottery_no_eligible_applicants", "lottery", lotteryID, 0, "", "", "")
		return nil, errors.New("no eligible applicants remaining")
	}

	rand.Seed(time.Now().UnixNano())
	winnerID := eligible[rand.Intn(len(eligible))]

	winner := &domain.LotteryWinner{
		LotteryID:     lottery.ID,
		ApplicantID:   winnerID,
		PositionOrder: len(winners) + 1,
	}

	if err := s.lotteryWinnerRepo.Create(ctx, winner); err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "lottery_spin_failed", "lottery_winner", winner.ID, 0, "", "", err.Error())
		return nil, err
	}

	_ = s.auditRepo.Log(ctx, adminID, "lottery_spin", "lottery_winner", winner.ID, 0, "", "", "")
	return winner, nil
}

func (s *LotteryService) CloseLottery(ctx context.Context, lotteryID, adminID string) error {
	lottery, err := s.lotteryRepo.GetByID(ctx, lotteryID)
	if err != nil {
		return err
	}

	lottery.Status = domain.LotteryCompleted
	if err := s.lotteryRepo.Create(ctx, lottery); err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "close_lottery_failed", "lottery", lottery.ID, 0, "", "", err.Error())
		return err
	}

	_ = s.auditRepo.Log(ctx, adminID, "close_lottery", "lottery", lottery.ID, 0, "", "", "")
	return nil
}
