package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/payment/internal/models"
	"github.com/simplcommerce-go/services/payment/internal/services"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	payments, total, err := h.paymentService.GetAllPayments(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(payments, page, pageSize, total))
}

func (h *PaymentHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid payment id")
	}

	payment, err := h.paymentService.GetPayment(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(payment))
}

func (h *PaymentHandler) Create(c fiber.Ctx) error {
	var input struct {
		OrderID              uint    `json:"orderId"`
		Amount               float64 `json:"amount"`
		PaymentFee           float64 `json:"paymentFee"`
		PaymentMethod        string  `json:"paymentMethod"`
		GatewayTransactionID string  `json:"gatewayTransactionId"`
		Status               string  `json:"status"`
		FailureMessage       string  `json:"failureMessage"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	payment := &models.Payment{
		OrderID:              input.OrderID,
		Amount:               input.Amount,
		PaymentFee:           input.PaymentFee,
		PaymentMethod:        input.PaymentMethod,
		GatewayTransactionID: input.GatewayTransactionID,
		Status:               "Succeeded",
		FailureMessage:       input.FailureMessage,
	}

	if input.Status != "" {
		payment.Status = input.Status
	}

	if err := h.paymentService.ProcessPayment(c.Context(), payment); err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(payment))
}

type PaymentProviderHandler struct {
	providerService services.PaymentProviderService
}

func NewPaymentProviderHandler(providerService services.PaymentProviderService) *PaymentProviderHandler {
	return &PaymentProviderHandler{providerService: providerService}
}

func (h *PaymentProviderHandler) GetAll(c fiber.Ctx) error {
	providers, err := h.providerService.GetAllProviders(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(providers))
}

func (h *PaymentProviderHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "provider id is required")
	}

	var input services.UpdateProviderInput
	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	provider, err := h.providerService.UpdateProvider(c.Context(), id, input)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(provider))
}
