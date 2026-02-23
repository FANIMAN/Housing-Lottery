package http

import (
	"os"

	"github.com/FANIMAN/housing-lottery/internal/usecase"
	"github.com/gofiber/fiber/v2"
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

func (h *AdminHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	token, err := h.usecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}


func (h *AdminHandler) VerifyPIN(c *fiber.Ctx) error {
	type request struct {
		PIN string `json:"pin"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if body.PIN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PIN is required"})
	}

	// Compare with env PIN
	adminPin := os.Getenv("ADMIN_PIN")
	if body.PIN != adminPin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin PIN"})
	}

	return c.JSON(fiber.Map{"success": true})
}


func (h *AdminHandler) List(c *fiber.Ctx) error {
	email := c.Query("email")
	id := c.Query("id")

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	admins, total, err := h.usecase.ListAdmins(c.Context(), email, id, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":     admins,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}