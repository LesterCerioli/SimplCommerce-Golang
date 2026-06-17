package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type StockController struct {
	svc *initializers.Services
}

func NewStockController(svc *initializers.Services) *StockController {
	return &StockController{svc: svc}
}

func (h *StockController) GetStock(c fiber.Ctx) error {
	productID := parseUint(c.Params("productId"))

	stock, err := h.svc.StockService.GetStock(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": stock})
}

func (h *StockController) UpdateStock(c fiber.Ctx) error {
	productID := parseUint(c.Params("productId"))
	createdByID := c.Locals("userID").(uint)

	var req struct {
		WarehouseID uint   `json:"warehouseId"`
		Quantity    int    `json:"quantity"`
		Note        string `json:"note"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	stock, err := h.svc.StockService.UpdateStock(c.Context(), productID, req.WarehouseID, req.Quantity, createdByID, req.Note)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": stock})
}

func (h *StockController) GetStockHistory(c fiber.Ctx) error {
	productID := parseUint(c.Params("productId"))
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	history, total, err := h.svc.StockService.GetStockHistory(c.Context(), productID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": history, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}
