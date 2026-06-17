package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
	"simplcommerce/services/implementations"
)

type PaymentController struct {
	svc *initializers.Services
}

func NewPaymentController(svc *initializers.Services) *PaymentController {
	return &PaymentController{svc: svc}
}

func (h *PaymentController) ListPayments(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "10"), 10)

	payments, total, err := h.svc.PaymentService.GetAllPayments(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": payments, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}

func (h *PaymentController) GetPayment(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid payment id"})
	}

	payment, err := h.svc.PaymentService.GetPayment(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": payment})
}

func (h *PaymentController) CreatePayment(c fiber.Ctx) error {
	var req struct {
		OrderID              string  `json:"orderId"`
		Amount               float64 `json:"amount"`
		PaymentFee           float64 `json:"paymentFee"`
		PaymentMethod        string  `json:"paymentMethod"`
		GatewayTransactionID string  `json:"gatewayTransactionId"`
		Status               string  `json:"status"`
		FailureMessage       string  `json:"failureMessage"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	payment := &implementations.PaymentResponse{
		OrderID:              req.OrderID,
		Amount:               req.Amount,
		PaymentFee:           req.PaymentFee,
		PaymentMethod:        req.PaymentMethod,
		GatewayTransactionID: req.GatewayTransactionID,
		Status:               "Succeeded",
		FailureMessage:       req.FailureMessage,
	}
	if req.Status != "" {
		payment.Status = req.Status
	}

	if err := h.svc.PaymentService.CreatePayment(c.Context(), payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": payment})
}

func (h *PaymentController) ListProviders(c fiber.Ctx) error {
	providers, err := h.svc.PaymentService.ListProviders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": providers})
}
