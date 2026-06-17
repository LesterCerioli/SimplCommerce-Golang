package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type PageController struct {
	svc *initializers.Services
}

func NewPageController(svc *initializers.Services) *PageController {
	return &PageController{svc: svc}
}

func (h *PageController) ListPages(c fiber.Ctx) error {
	pages, err := h.svc.PageService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": pages})
}

func (h *PageController) GetPage(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid slug"})
	}

	page, err := h.svc.PageService.GetBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": page})
}

func (h *PageController) CreatePage(c fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Body        string `json:"body"`
		IsPublished bool   `json:"isPublished"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "name is required"})
	}

	page, err := h.svc.PageService.Create(c.Context(), req.Name, req.Slug, req.Body, req.IsPublished)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": page})
}

func (h *PageController) UpdatePage(c fiber.Ctx) error {
	pageID := parseUint(c.Params("id"))

	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Body        string `json:"body"`
		IsPublished bool   `json:"isPublished"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	page, err := h.svc.PageService.Update(c.Context(), pageID, req.Name, req.Slug, req.Body, req.IsPublished)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": page})
}

func (h *PageController) DeletePage(c fiber.Ctx) error {
	pageID := parseUint(c.Params("id"))

	if err := h.svc.PageService.Delete(c.Context(), pageID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
