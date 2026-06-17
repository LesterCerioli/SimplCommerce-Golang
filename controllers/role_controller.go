package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

type RoleController struct {
	roleService *implementations.RoleService
}

func NewRoleController(roleService *implementations.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (ctrl *RoleController) List(c fiber.Ctx) error {
	roles, err := ctrl.roleService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": roles})
}

func (ctrl *RoleController) Create(c fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	role, err := ctrl.roleService.Create(c.Context(), req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

func (ctrl *RoleController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	if err := ctrl.roleService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role deleted successfully"})
}
