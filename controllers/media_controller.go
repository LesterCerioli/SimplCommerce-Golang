package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

type MediaController struct {
	mediaService *implementations.MediaService
}

func NewMediaController(mediaService *implementations.MediaService) *MediaController {
	return &MediaController{mediaService: mediaService}
}

func (ctrl *MediaController) List(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "10"), 10)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	media, total, err := ctrl.mediaService.Paginate(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":     media,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *MediaController) Upload(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	caption := c.FormValue("caption", "")

	req := implementations.CreateMediaRequest{
		Caption:   caption,
		FileName:  file.Filename,
		FileSize:  file.Size,
		MediaType: 0,
	}

	media, err := ctrl.mediaService.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}

func (ctrl *MediaController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid media ID"})
	}

	if err := ctrl.mediaService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Media deleted successfully"})
}
