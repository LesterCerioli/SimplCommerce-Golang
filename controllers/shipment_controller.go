package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
	"simplcommerce/services/implementations"
)

type ShipmentController struct {
	svc *initializers.Services
}

func NewShipmentController(svc *initializers.Services) *ShipmentController {
	return &ShipmentController{svc: svc}
}

func (h *ShipmentController) ListShipments(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "10"), 10)

	shipments, total, err := h.svc.ShipmentService.ListShipments(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": shipments, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}

func (h *ShipmentController) GetShipment(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid shipment id"})
	}

	shipment, err := h.svc.ShipmentService.GetShipment(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": shipment})
}

func (h *ShipmentController) CreateShipment(c fiber.Ctx) error {
	var req implementations.CreateShipmentInput
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	createdByID := c.Locals("userID").(string)

	shipment, err := h.svc.ShipmentService.CreateShipment(c.Context(), req, createdByID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": shipment})
}
