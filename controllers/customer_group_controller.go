package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

type CustomerGroupController struct {
	groupService *implementations.CustomerGroupService
}

func NewCustomerGroupController(groupService *implementations.CustomerGroupService) *CustomerGroupController {
	return &CustomerGroupController{groupService: groupService}
}

func (ctrl *CustomerGroupController) List(c fiber.Ctx) error {
	groups, err := ctrl.groupService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": groups})
}

func (ctrl *CustomerGroupController) Create(c fiber.Ctx) error {
	var req implementations.CreateCustomerGroupRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	group, err := ctrl.groupService.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func (ctrl *CustomerGroupController) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer group ID"})
	}

	var req implementations.UpdateCustomerGroupRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	group, err := ctrl.groupService.Update(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(group)
}

func (ctrl *CustomerGroupController) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer group ID"})
	}

	if err := ctrl.groupService.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Customer group deleted successfully"})
}
