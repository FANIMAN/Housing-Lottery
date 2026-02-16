package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

type SubcityHandler struct {
	usecase *usecase.SubcityUsecase
}

func NewSubcityHandler(u *usecase.SubcityUsecase) *SubcityHandler {
	return &SubcityHandler{usecase: u}
}

// Create subcity
func (h *SubcityHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.usecase.Create(c.Context(), req.Name); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"message": "subcity created"})
}

// List all subcities
func (h *SubcityHandler) List(c *fiber.Ctx) error {
	subcities, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(subcities)
}

// Update subcity
func (h *SubcityHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.usecase.Update(c.Context(), id, req.Name); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"message": "subcity updated"})
}

// Delete subcity
func (h *SubcityHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.usecase.Delete(c.Context(), id); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"message": "subcity deleted"})
}
