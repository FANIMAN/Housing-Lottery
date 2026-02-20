package http

import (
	"strconv"

	"github.com/FANIMAN/housing-lottery/internal/domain"
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

// func (h *UploadHandler) ListApplicants(c *fiber.Ctx) error {
// 	adminID := c.Locals("admin_id")
// 	if adminID == nil {
// 		return fiber.ErrUnauthorized
// 	}

// 	subcityIDStr := c.Query("subcityId")
// 	search := c.Query("search")

// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	limit, _ := strconv.Atoi(c.Query("limit", "10"))

// 	applicants, err := h.usecase.GetApplicants(
// 		c.Context(),
// 		subcityIDStr,
// 		search,
// 		page,
// 		limit,
// 	)

// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch applicants")
// 	}

// 	return c.JSON(fiber.Map{
// 		"data": applicants,
// 	})
// }


func (h *UploadHandler) ListApplicants(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	if adminID == nil {
		return fiber.ErrUnauthorized
	}

	subcityIDStr := c.Query("subcityId")
	search := c.Query("search")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	applicants, total, err := h.usecase.GetApplicants(
		c.Context(),
		subcityIDStr,
		search,
		page,
		limit,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch applicants")
	}

	// IMPORTANT: never return null
	if applicants == nil {
		applicants = []*domain.Applicant{}
	}

	return c.JSON(fiber.Map{
		"data":  applicants,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}