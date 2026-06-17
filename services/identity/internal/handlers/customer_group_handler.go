package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type CustomerGroupHandler struct {
	groupService services.CustomerGroupService
}

func NewCustomerGroupHandler(groupService services.CustomerGroupService) *CustomerGroupHandler {
	return &CustomerGroupHandler{groupService: groupService}
}

func (h *CustomerGroupHandler) Index(c fiber.Ctx) error {
	groups, err := h.groupService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(groups))
}

func (h *CustomerGroupHandler) Create(c fiber.Ctx) error {
	var group models.CustomerGroup
	if err := c.Bind().JSON(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	result, err := h.groupService.Create(c.Context(), &group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(result))
}
