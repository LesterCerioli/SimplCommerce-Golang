package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type createCategoryRequest struct {
	Name            string `json:"name" validate:"required"`
	Slug            string `json:"slug,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayOrder    int    `json:"displayOrder,omitempty"`
	IsPublished     bool   `json:"isPublished,omitempty"`
	IncludeInMenu   bool   `json:"includeInMenu,omitempty"`
	ParentID        *uint  `json:"parentId,omitempty"`
	ThumbnailImageID *uint `json:"thumbnailImageId,omitempty"`
}

type updateCategoryRequest struct {
	Name            string `json:"name"`
	Slug            string `json:"slug,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayOrder    int    `json:"displayOrder,omitempty"`
	IsPublished     bool   `json:"isPublished,omitempty"`
	IncludeInMenu   bool   `json:"includeInMenu,omitempty"`
	ParentID        *uint  `json:"parentId,omitempty"`
	ThumbnailImageID *uint `json:"thumbnailImageId,omitempty"`
}

func (h *CategoryHandler) Index(c fiber.Ctx) error {
	categories, err := h.service.GetCategoryTree(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(categories))
}

func (h *CategoryHandler) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.service.GetCategoryBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, "category not found")
	}

	return c.JSON(response.Success(category))
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	var req createCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	category := &catalogmodels.Category{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		DisplayOrder:     req.DisplayOrder,
		IsPublished:      req.IsPublished,
		IncludeInMenu:    req.IncludeInMenu,
		ParentID:         req.ParentID,
		ThumbnailImageID: req.ThumbnailImageID,
	}

	if err := h.service.CreateCategory(c.Context(), category); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(category))
}

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid category id")
	}

	var req updateCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	category, err := h.service.GetCategoryByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	category.Description = req.Description
	category.DisplayOrder = req.DisplayOrder
	category.IsPublished = req.IsPublished
	category.IncludeInMenu = req.IncludeInMenu
	category.ParentID = req.ParentID
	category.ThumbnailImageID = req.ThumbnailImageID

	if err := h.service.UpdateCategory(c.Context(), category); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(category))
}

func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid category id")
	}

	if err := h.service.DeleteCategory(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
