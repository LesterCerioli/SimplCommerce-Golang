package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
	"simplcommerce/services/implementations"
)

type CheckoutController struct {
	svc *initializers.Services
}

func NewCheckoutController(svc *initializers.Services) *CheckoutController {
	return &CheckoutController{svc: svc}
}

func (h *CheckoutController) GetCheckout(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	info, err := h.svc.CheckoutService.GetCheckout(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": info})
}

func (h *CheckoutController) ProcessCheckout(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	var req implementations.CheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	order, err := h.svc.CheckoutService.ProcessCheckout(c.Context(), customerID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": order})
}
