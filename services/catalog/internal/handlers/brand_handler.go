package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

type BrandHandler struct {
	service services.BrandService
}

func NewBrandHandler(service services.BrandService) *BrandHandler {
	return &BrandHandler{service: service}
}

type createBrandRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublished bool   `json:"isPublished,omitempty"`
}

type updateBrandRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublished bool   `json:"isPublished,omitempty"`
}

func (h *BrandHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	brands, total, err := h.service.GetBrands(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(brands, page, pageSize, total))
}

func (h *BrandHandler) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	brand, err := h.service.GetBrandBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, "brand not found")
	}

	return c.JSON(response.Success(brand))
}

func (h *BrandHandler) Create(c fiber.Ctx) error {
	var req createBrandRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	brand := &catalogmodels.Brand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsPublished: req.IsPublished,
	}

	if err := h.service.CreateBrand(c.Context(), brand); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(brand))
}

func (h *BrandHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid brand id")
	}

	var req updateBrandRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	brand, err := h.service.GetBrandByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "brand not found")
	}

	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Slug != "" {
		brand.Slug = req.Slug
	}
	brand.Description = req.Description
	brand.IsPublished = req.IsPublished

	if err := h.service.UpdateBrand(c.Context(), brand); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(brand))
}

func (h *BrandHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid brand id")
	}

	if err := h.service.DeleteBrand(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
