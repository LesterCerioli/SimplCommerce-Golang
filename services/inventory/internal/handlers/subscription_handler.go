package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/inventory/internal/models"
	"github.com/simplcommerce-go/services/inventory/internal/repositories"
)

type SubscriptionHandler struct {
	subRepo repositories.SubscriptionRepository
}

func NewSubscriptionHandler(subRepo repositories.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{subRepo: subRepo}
}

type subscribeRequest struct {
	ProductID     uint   `json:"productId"`
	CustomerEmail string `json:"customerEmail"`
}

func (h *SubscriptionHandler) Subscribe(c fiber.Ctx) error {
	var req subscribeRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	if req.CustomerEmail == "" {
		return response.Error(fiber.StatusBadRequest, "customer email is required")
	}

	if req.ProductID == 0 {
		return response.Error(fiber.StatusBadRequest, "product id is required")
	}

	existing, err := h.subRepo.FindByProductAndEmail(c.Context(), req.ProductID, req.CustomerEmail)
	if err == nil && existing != nil {
		return c.JSON(response.Success(existing))
	}

	sub := &models.ProductBackInStockSubscription{
		ProductID:     req.ProductID,
		CustomerEmail: req.CustomerEmail,
	}

	if err := h.subRepo.Create(c.Context(), sub); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(sub))
}
