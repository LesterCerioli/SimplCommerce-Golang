package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/internal/response"
	"simplcommerce/services/implementations"
)

type CategoryController struct {
	service *implementations.CategoryService
}

func NewCategoryController(service *implementations.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (h *CategoryController) Index(c fiber.Ctx) error {
	catList, err := h.service.GetCategories(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	tree := implementations.BuildCategoryTree(catList)
	return c.JSON(response.Success(tree))
}

func (h *CategoryController) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.service.GetCategoryBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(category))
}

func (h *CategoryController) Create(c fiber.Ctx) error {
	var req implementations.CreateCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	category, err := h.service.CreateCategory(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(category))
}

func (h *CategoryController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid category id")
	}

	var req implementations.UpdateCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	category, err := h.service.UpdateCategory(c.Context(), id, &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(category))
}

func (h *CategoryController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid category id")
	}

	if err := h.service.DeleteCategory(c.Context(), id); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
