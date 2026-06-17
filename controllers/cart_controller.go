package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type CartController struct {
	svc *initializers.Services
}

func NewCartController(svc *initializers.Services) *CartController {
	return &CartController{svc: svc}
}

func (h *CartController) GetCart(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	items, err := h.svc.CartService.GetCart(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

func (h *CartController) AddItem(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	var req struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	if err := h.svc.CartService.AddItem(c.Context(), customerID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": nil})
}

func (h *CartController) UpdateQuantity(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	id := c.Params("id")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	if err := h.svc.CartService.UpdateQuantity(c.Context(), customerID, id, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *CartController) RemoveItem(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	id := c.Params("id")
	if err := h.svc.CartService.RemoveItem(c.Context(), customerID, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *CartController) ClearCart(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	if err := h.svc.CartService.ClearCart(c.Context(), customerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
