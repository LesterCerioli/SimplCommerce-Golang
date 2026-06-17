package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/internal/response"
	"simplcommerce/services/implementations"
)

type ProductOptionController struct {
	service *implementations.ProductOptionService
}

type ProductTemplateController struct {
	service *implementations.ProductTemplateService
}

func NewProductOptionController(service *implementations.ProductOptionService) *ProductOptionController {
	return &ProductOptionController{service: service}
}

func NewProductTemplateController(service *implementations.ProductTemplateService) *ProductTemplateController {
	return &ProductTemplateController{service: service}
}

func (h *ProductOptionController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	options, total, err := h.service.GetOptions(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(options, page, pageSize, total))
}

func (h *ProductOptionController) Create(c fiber.Ctx) error {
	var req implementations.CreateProductOptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	option, err := h.service.CreateOption(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(option))
}

func (h *ProductOptionController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid option id")
	}

	var req implementations.UpdateProductOptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	option, err := h.service.UpdateOption(c.Context(), id, &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(option))
}

func (h *ProductOptionController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid option id")
	}

	if err := h.service.DeleteOption(c.Context(), id); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *ProductTemplateController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	templates, total, err := h.service.GetTemplates(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(templates, page, pageSize, total))
}

func (h *ProductTemplateController) Create(c fiber.Ctx) error {
	var req implementations.CreateProductTemplateRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	template, err := h.service.CreateTemplate(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(template))
}
