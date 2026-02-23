package http

import (
	"time"

	"github.com/FANIMAN/housing-lottery/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LotteryHandler struct {
	service *usecase.LotteryService
}

type StartLotteryRequest struct {
	Name string `json:"name"`
}

func NewLotteryHandler(service *usecase.LotteryService) *LotteryHandler {
	return &LotteryHandler{service: service}
}

// Spin a lottery winner
func (h *LotteryHandler) Spin(c *fiber.Ctx) error {
	lotteryID := c.Params("id")
	adminID := c.Locals("admin_id").(string)

	result, err := h.service.SpinLottery(c.Context(), lotteryID, adminID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"lottery_id":     result.Winner.LotteryID,
		"position_order": result.Winner.PositionOrder,
		"announced_at":   result.Winner.AnnouncedAt,

		"applicant": fiber.Map{
			"id":                          result.Applicant.ID,
			"full_name":                   result.Applicant.FullName,
			"condominium_registration_id": result.Applicant.CondominiumRegistrationID,
			"subcity_id":                  result.Applicant.SubcityID,
			"created_at":                  result.Applicant.CreatedAt,
		},
	})
}

// Start lottery for subcity
// func (h *LotteryHandler) Start(c *fiber.Ctx) error {
// 	subcityIDStr := c.Params("id")
// 	subcityID, err := uuid.Parse(subcityIDStr)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "invalid subcity ID"})
// 	}

// 	adminID := c.Locals("admin_id").(string)

// 	lottery, err := h.service.StartLottery(c.Context(), subcityID, adminID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(lottery)
// }

func (h *LotteryHandler) Start(c *fiber.Ctx) error {
	subcityIDStr := c.Params("id")
	subcityID, err := uuid.Parse(subcityIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid subcity ID"})
	}

	var req StartLotteryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "lottery name is required"})
	}

	adminID := c.Locals("admin_id").(string)

	lottery, err := h.service.StartLottery(c.Context(), subcityID, req.Name, adminID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lottery)
}

// Close lottery
func (h *LotteryHandler) Close(c *fiber.Ctx) error {
	lotteryID := c.Params("id")
	adminID := c.Locals("admin_id").(string)
	if err := h.service.CloseLottery(c.Context(), lotteryID, adminID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "lottery closed"})
}



func (h *LotteryHandler) ListWinners(c *fiber.Ctx) error {
	subcity := c.Query("subcity")
	fullName := c.Query("full_name")
	lotteryName := c.Query("lottery_name")
	from := c.Query("from")
	to := c.Query("to")

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	var fromDate, toDate *time.Time
	if from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			fromDate = &t
		}
	}
	if to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			toDate = &t
		}
	}

	winners,total, err := h.service.ListWinners(c.Context(), subcity, fullName, lotteryName, fromDate, toDate, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"data":      winners,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}