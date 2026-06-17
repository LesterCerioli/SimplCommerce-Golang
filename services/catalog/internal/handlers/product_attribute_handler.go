package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

type ProductAttributeHandler struct {
	service services.ProductAttributeService
}

func NewProductAttributeHandler(service services.ProductAttributeService) *ProductAttributeHandler {
	return &ProductAttributeHandler{service: service}
}

type createAttributeRequest struct {
	Name    string `json:"name" validate:"required"`
	GroupID uint   `json:"groupId"`
}

func (h *ProductAttributeHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	attributes, total, err := h.service.GetAttributes(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(attributes, page, pageSize, total))
}

func (h *ProductAttributeHandler) Create(c fiber.Ctx) error {
	var req createAttributeRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	attr := &catalogmodels.ProductAttribute{
		Name:    req.Name,
		GroupID: req.GroupID,
	}

	if err := h.service.CreateAttribute(c.Context(), attr); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(attr))
}

type ProductAttributeGroupHandler struct {
	service services.ProductAttributeGroupService
}

func NewProductAttributeGroupHandler(service services.ProductAttributeGroupService) *ProductAttributeGroupHandler {
	return &ProductAttributeGroupHandler{service: service}
}

type createAttributeGroupRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *ProductAttributeGroupHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	groups, total, err := h.service.GetGroups(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(groups, page, pageSize, total))
}

func (h *ProductAttributeGroupHandler) Create(c fiber.Ctx) error {
	var req createAttributeGroupRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	group := &catalogmodels.ProductAttributeGroup{
		Name: req.Name,
	}

	if err := h.service.CreateGroup(c.Context(), group); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(group))
}
