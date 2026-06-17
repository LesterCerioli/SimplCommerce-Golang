package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) Index(c fiber.Ctx) error {
	roles, err := h.roleService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(roles))
}

func (h *RoleHandler) Create(c fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "name is required"))
	}

	role, err := h.roleService.Create(c.Context(), req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(role))
}
