package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

type DashboardHandler struct {
	usecase *usecase.DashboardUsecase
}

func NewDashboardHandler(u *usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{usecase: u}
}

// GET /api/dashboard/summary
func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	subcityId := c.Query("subcityId")
	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	summary, err := h.usecase.GetSummary(subcityId, status, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(summary)
}

func (h *DashboardHandler) ListSubcities(c *fiber.Ctx) error {
	subcities, err := h.usecase.ListSubcities()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subcities)
}

func (h *DashboardHandler) ListLotteries(c *fiber.Ctx) error {
	lotteries, err := h.usecase.ListLotteries()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(lotteries)
}