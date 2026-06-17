package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/orders/internal/services"
)

type CheckoutHandler struct {
	checkoutService services.CheckoutService
}

func NewCheckoutHandler(checkoutService services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService: checkoutService}
}

func (h *CheckoutHandler) GetCheckout(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	info, err := h.checkoutService.GetCheckout(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(info))
}

func (h *CheckoutHandler) ProcessCheckout(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	var req services.CheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	order, err := h.checkoutService.ProcessCheckout(c.Context(), customerID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(order))
}
