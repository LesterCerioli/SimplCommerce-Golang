package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/internal/response"
	"simplcommerce/services/implementations"
)

type ProductAttributeController struct {
	service *implementations.ProductAttributeService
}

type ProductAttributeGroupController struct {
	service *implementations.ProductAttributeGroupService
}

func NewProductAttributeController(service *implementations.ProductAttributeService) *ProductAttributeController {
	return &ProductAttributeController{service: service}
}

func NewProductAttributeGroupController(service *implementations.ProductAttributeGroupService) *ProductAttributeGroupController {
	return &ProductAttributeGroupController{service: service}
}

func (h *ProductAttributeController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	attrs, total, err := h.service.GetAttributes(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(attrs, page, pageSize, total))
}

func (h *ProductAttributeController) Create(c fiber.Ctx) error {
	var req implementations.CreateProductAttributeRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	attr, err := h.service.CreateAttribute(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(attr))
}

func (h *ProductAttributeController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid attribute id")
	}

	if err := h.service.DeleteAttribute(c.Context(), id); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *ProductAttributeGroupController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	groups, total, err := h.service.GetGroups(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(groups, page, pageSize, total))
}

func (h *ProductAttributeGroupController) Create(c fiber.Ctx) error {
	var req implementations.CreateProductAttributeGroupRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	group, err := h.service.CreateGroup(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(group))
}

func (h *ProductAttributeGroupController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid group id")
	}

	var req implementations.UpdateProductAttributeGroupRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	group, err := h.service.UpdateGroup(c.Context(), id, &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(group))
}
