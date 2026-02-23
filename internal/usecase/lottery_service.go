package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/domain"
	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LotteryService struct {
	db                *pgxpool.Pool
	lotteryRepo       interfaces.LotteryRepository
	applicantRepo     interfaces.ApplicantRepository
	lotteryWinnerRepo interfaces.LotteryWinnerRepository
	auditRepo         interfaces.AuditRepository
}

func NewLotteryService(
	db *pgxpool.Pool,
	lotteryRepo interfaces.LotteryRepository,
	applicantRepo interfaces.ApplicantRepository,
	lotteryWinnerRepo interfaces.LotteryWinnerRepository,
	auditRepo interfaces.AuditRepository,
) *LotteryService {
	return &LotteryService{
		db:                db,
		lotteryRepo:       lotteryRepo,
		applicantRepo:     applicantRepo,
		lotteryWinnerRepo: lotteryWinnerRepo,
		auditRepo:         auditRepo,
	}
}

func (s *LotteryService) StartLottery(
	ctx context.Context,
	subcityID uuid.UUID,
	name string,
	adminID string,
) (*domain.Lottery, error) {

	// func (s *LotteryService) StartLottery(ctx context.Context, subcityID uuid.UUID, name string, adminID string) (*domain.Lottery, error) {
	applicants, err := s.applicantRepo.GetAllBySubcityID(ctx, subcityID)
	if err != nil {
		_ = s.auditRepo.Log(ctx, adminID, "start_lottery_failed", "lottery", "", 0, "", "", err.Error())
		return nil, err
	}

	lottery := &domain.Lottery{
		ID:              uuid.NewString(),
		Name:            name,
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

func (s *LotteryService) SpinLottery(
	ctx context.Context,
	lotteryID string,
	adminID string,
) (*domain.LotteryWinnerResult, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Get Lottery
	lottery, err := s.lotteryRepo.GetByID(ctx, lotteryID)
	if err != nil {
		return nil, err
	}

	if lottery.Status != domain.LotteryPending {
		return nil, errors.New("lottery is not active")
	}

	// 2️⃣ Get Applicants
	applicants, err := s.applicantRepo.GetAllBySubcityID(ctx, lottery.SubcityID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Get Existing Winners
	winners, err := s.lotteryWinnerRepo.GetWinnersByLottery(ctx, lotteryID)
	if err != nil {
		return nil, err
	}

	if len(winners) >= len(applicants) {
		return nil, errors.New("all applicants have already won")
	}

	// 4️⃣ Build winner lookup map
	wonIDs := make(map[string]bool)
	for _, w := range winners {
		wonIDs[w.ApplicantID] = true
	}

	// 5️⃣ Build eligible list
	var eligible []string
	for _, a := range applicants {
		aID := a.ID.String()
		if !wonIDs[aID] {
			eligible = append(eligible, aID)
		}
	}

	// 6️⃣ If no eligible left → complete lottery
	if len(eligible) == 0 {
		if err := s.lotteryRepo.UpdateStatus(ctx, tx, lotteryID, domain.LotteryCompleted); err != nil {
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return nil, errors.New("no eligible applicants remaining")
	}

	// 7️⃣ Pick random winner
	rand.Seed(time.Now().UnixNano())
	winnerID := eligible[rand.Intn(len(eligible))]

	winner := &domain.LotteryWinner{
		LotteryID:     lottery.ID,
		ApplicantID:   winnerID,
		PositionOrder: len(winners) + 1,
	}

	// 8️⃣ Save winner
	if err := s.lotteryWinnerRepo.CreateTx(ctx, tx, winner); err != nil {
		return nil, err
	}

	// 9️⃣ Increment winners count
	if err := s.lotteryRepo.IncrementWinnersCount(ctx, tx, lotteryID); err != nil {
		return nil, err
	}

	// 🔟 Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// 1️⃣1️⃣ Fetch full applicant info
	applicantUUID, err := uuid.Parse(winnerID)
	if err != nil {
		return nil, err
	}

	applicant, err := s.applicantRepo.GetByID(ctx, applicantUUID)
	if err != nil {
		return nil, err
	}

	// 1️⃣2️⃣ Return combined result
	return &domain.LotteryWinnerResult{
		Winner:    winner,
		Applicant: applicant,
	}, nil
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


func (s *LotteryService) ListWinners(
	ctx context.Context,
	subcityName, fullName, lotteryName string,
	fromDate, toDate *time.Time,
	page, pageSize int,
) ([]*domain.LotteryWinnerResponse, int, error) {

	offset := (page - 1) * pageSize

	return s.lotteryWinnerRepo.ListWinners(
		ctx,
		subcityName,
		fullName,
		lotteryName,
		fromDate,
		toDate,
		pageSize,
		offset,
	)
}