package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

type AdminHandler struct {
	usecase *usecase.AdminUsecase
}

func NewAdminHandler(u *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{usecase: u}
}

func (h *AdminHandler) Register(c *fiber.Ctx) error {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	err := h.usecase.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"message": "admin created"})
}
