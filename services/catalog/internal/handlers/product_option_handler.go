package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

type ProductOptionHandler struct {
	service services.ProductOptionService
}

func NewProductOptionHandler(service services.ProductOptionService) *ProductOptionHandler {
	return &ProductOptionHandler{service: service}
}

type createOptionRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *ProductOptionHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	options, total, err := h.service.GetOptions(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(options, page, pageSize, total))
}

func (h *ProductOptionHandler) Create(c fiber.Ctx) error {
	var req createOptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	option := &catalogmodels.ProductOption{
		Name: req.Name,
	}

	if err := h.service.CreateOption(c.Context(), option); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(option))
}

type ProductTemplateHandler struct {
	service services.ProductTemplateService
}

func NewProductTemplateHandler(service services.ProductTemplateService) *ProductTemplateHandler {
	return &ProductTemplateHandler{service: service}
}

type createTemplateRequest struct {
	Name       string `json:"name" validate:"required"`
	AttributeIDs []uint `json:"attributeIds,omitempty"`
}

func (h *ProductTemplateHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	templates, total, err := h.service.GetTemplates(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(templates, page, pageSize, total))
}

func (h *ProductTemplateHandler) Create(c fiber.Ctx) error {
	var req createTemplateRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	template := &catalogmodels.ProductTemplate{
		Name: req.Name,
	}

	for _, attrID := range req.AttributeIDs {
		template.Attributes = append(template.Attributes, catalogmodels.ProductTemplateProductAttribute{
			ProductTemplateID:  0,
			ProductAttributeID: attrID,
		})
	}

	if err := h.service.CreateTemplate(c.Context(), template); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(template))
}
