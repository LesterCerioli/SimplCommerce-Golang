package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/internal/response"
	"simplcommerce/services/implementations"
)

type BrandController struct {
	service *implementations.BrandService
}

func NewBrandController(service *implementations.BrandService) *BrandController {
	return &BrandController{service: service}
}

func (h *BrandController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	brands, total, err := h.service.GetBrands(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(brands, page, pageSize, total))
}

func (h *BrandController) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	brand, err := h.service.GetBrandBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(brand))
}

func (h *BrandController) Create(c fiber.Ctx) error {
	var req implementations.CreateBrandRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	brand, err := h.service.CreateBrand(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(brand))
}

func (h *BrandController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid brand id")
	}

	var req implementations.UpdateBrandRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	brand, err := h.service.UpdateBrand(c.Context(), id, &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(brand))
}

func (h *BrandController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid brand id")
	}

	if err := h.service.DeleteBrand(c.Context(), id); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
