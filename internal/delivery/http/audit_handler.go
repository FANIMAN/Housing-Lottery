package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

type AuditHandler struct {
	service *usecase.AuditService
}

func NewAuditHandler(service *usecase.AuditService) *AuditHandler {
	return &AuditHandler{service}
}

func (h *AuditHandler) List(c *fiber.Ctx) error {

	adminID := c.Query("adminId")
	action := c.Query("action")
	entityType := c.Query("entityType")

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	var fromDate *time.Time
	var toDate *time.Time

	if f := c.Query("fromDate"); f != "" {
		t, err := time.Parse(time.RFC3339, f)
		if err == nil {
			fromDate = &t
		}
	}

	if t := c.Query("toDate"); t != "" {
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			toDate = &parsed
		}
	}

	data, total, err := h.service.List(
		c.Context(),
		adminID,
		action,
		entityType,
		fromDate,
		toDate,
		page,
		pageSize,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":     data,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}