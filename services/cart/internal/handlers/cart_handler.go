package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/cart/internal/services"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

type addItemRequest struct {
	ProductID uint `json:"productId" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type updateQuantityRequest struct {
	Quantity int `json:"quantity" validate:"required"`
}

func (h *CartHandler) AddItem(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	var req addItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	if err := h.cartService.AddItem(c.Context(), customerID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil))
}

func (h *CartHandler) UpdateQuantity(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid item id",
		})
	}

	var req updateQuantityRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := h.cartService.UpdateQuantity(c.Context(), customerID, uint(id), req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(nil))
}

func (h *CartHandler) RemoveItem(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid item id",
		})
	}

	if err := h.cartService.RemoveItem(c.Context(), customerID, uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(nil))
}

func (h *CartHandler) GetCart(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	items, err := h.cartService.GetCart(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(items))
}

func (h *CartHandler) ClearCart(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	if err := h.cartService.ClearCart(c.Context(), customerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(nil))
}
