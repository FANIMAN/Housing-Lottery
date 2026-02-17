package http

import (
	"github.com/FANIMAN/housing-lottery/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LotteryHandler struct {
	service *usecase.LotteryService
}

func NewLotteryHandler(service *usecase.LotteryService) *LotteryHandler {
	return &LotteryHandler{service: service}
}

// Spin a lottery winner
func (h *LotteryHandler) Spin(c *fiber.Ctx) error {
	lotteryID := c.Params("id")
	adminID := c.Locals("admin_id").(string)

	winner, err := h.service.SpinLottery(c.Context(), lotteryID, adminID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"winner_id":      winner.ApplicantID,
		"lottery_id":     winner.LotteryID,
		"position_order": winner.PositionOrder,
		"announced_at":   winner.AnnouncedAt,
	})
}

// Start lottery for subcity
func (h *LotteryHandler) Start(c *fiber.Ctx) error {
	subcityIDStr := c.Params("id")
	subcityID, err := uuid.Parse(subcityIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid subcity ID"})
	}

	adminID := c.Locals("admin_id").(string)

	lottery, err := h.service.StartLottery(c.Context(), subcityID, adminID)
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
