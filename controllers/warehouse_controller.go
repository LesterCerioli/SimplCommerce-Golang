package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type WarehouseController struct {
	svc *initializers.Services
}

func NewWarehouseController(svc *initializers.Services) *WarehouseController {
	return &WarehouseController{svc: svc}
}

func (h *WarehouseController) List(c fiber.Ctx) error {
	warehouses, err := h.svc.WarehouseService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": warehouses})
}

func (h *WarehouseController) Create(c fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		AddressID uint   `json:"addressId"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "name is required"})
	}

	warehouse, err := h.svc.WarehouseService.Create(c.Context(), req.Name, req.AddressID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": warehouse})
}

func (h *WarehouseController) Update(c fiber.Ctx) error {
	warehouseID := parseUint(c.Params("id"))

	var req struct {
		Name      string `json:"name"`
		AddressID uint   `json:"addressId"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	warehouse, err := h.svc.WarehouseService.Update(c.Context(), warehouseID, req.Name, req.AddressID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": warehouse})
}

func (h *WarehouseController) Delete(c fiber.Ctx) error {
	warehouseID := parseUint(c.Params("id"))

	if err := h.svc.WarehouseService.Delete(c.Context(), warehouseID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
