package http

import (
	"github.com/FANIMAN/housing-lottery/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	usecase *usecase.UploadService
}

func NewUploadHandler(u *usecase.UploadService) *UploadHandler {
	return &UploadHandler{usecase: u}
}

func (h *UploadHandler) UploadApplicants(c *fiber.Ctx) error {
	// Log token and admin_id
	adminID := c.Locals("admin_id")
	if adminID == nil {
		return fiber.ErrUnauthorized
	}

	// Log subcityID
	subcityID := c.Params("id")

	// Check if file is present in form-data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Try opening the file
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	defer file.Close()

	// Proceed with processing the file
	inserted, skipped, err := h.usecase.ProcessExcel(c.Context(), subcityID, adminID.(string), file, fileHeader.Filename)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"message":            "Upload completed",
		"inserted":           inserted,
		"skipped_duplicates": skipped,
	})
}
