package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/cms/internal/models"
	"github.com/simplcommerce-go/services/cms/internal/services"
)

type PageHandler struct {
	pageService services.PageService
}

func NewPageHandler(pageService services.PageService) *PageHandler {
	return &PageHandler{pageService: pageService}
}

type createPageRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Body        string `json:"body"`
	IsPublished bool   `json:"isPublished"`
}

func (h *PageHandler) List(c fiber.Ctx) error {
	pages, err := h.pageService.FindAll(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response.Success(pages))
}

func (h *PageHandler) GetBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	page, err := h.pageService.FindBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, "page not found")
	}
	return c.JSON(response.Success(page))
}

func (h *PageHandler) Create(c fiber.Ctx) error {
	var req createPageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	page := &models.Page{
		Name:        req.Name,
		Slug:        req.Slug,
		Body:        req.Body,
		IsPublished: req.IsPublished,
	}

	if err := h.pageService.Create(c.Context(), page); err != nil {
		if err == services.ErrSlugExists {
			return response.Error(fiber.StatusConflict, err.Error())
		}
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(page))
}

func (h *PageHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	var req createPageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	page, err := h.pageService.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "page not found")
	}

	page.Name = req.Name
	page.Slug = req.Slug
	page.Body = req.Body
	page.IsPublished = req.IsPublished

	if err := h.pageService.Update(c.Context(), page); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(page))
}

func (h *PageHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.pageService.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
