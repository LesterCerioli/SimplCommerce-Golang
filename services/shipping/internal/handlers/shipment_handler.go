package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/shipping/internal/services"
)

type ShipmentHandler struct {
	shipmentService services.ShipmentService
}

func NewShipmentHandler(shipmentService services.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{shipmentService: shipmentService}
}

func (h *ShipmentHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	shipments, total, err := h.shipmentService.GetAllShipments(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(shipments, page, pageSize, total))
}

func (h *ShipmentHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid shipment id")
	}

	shipment, err := h.shipmentService.GetShipment(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(shipment))
}

func (h *ShipmentHandler) Create(c fiber.Ctx) error {
	var input struct {
		OrderID        uint                     `json:"orderId"`
		TrackingNumber string                   `json:"trackingNumber"`
		WarehouseID    uint                     `json:"warehouseId"`
		VendorID       *uint                    `json:"vendorId"`
		Items          []services.ShipmentItemInput `json:"items"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	userID, _ := c.Locals("userID").(uint)

	svcInput := services.CreateShipmentInput{
		OrderID:        input.OrderID,
		TrackingNumber: input.TrackingNumber,
		WarehouseID:    input.WarehouseID,
		VendorID:       input.VendorID,
		Items:          input.Items,
	}

	shipment, err := h.shipmentService.CreateShipment(c.Context(), svcInput, userID)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(shipment))
}

type ShippingProviderHandler struct {
	providerService services.ShippingProviderService
}

func NewShippingProviderHandler(providerService services.ShippingProviderService) *ShippingProviderHandler {
	return &ShippingProviderHandler{providerService: providerService}
}

func (h *ShippingProviderHandler) GetAll(c fiber.Ctx) error {
	providers, err := h.providerService.GetAllProviders(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(providers))
}

func (h *ShippingProviderHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "provider id is required")
	}

	var input services.UpdateShippingProviderInput
	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	provider, err := h.providerService.UpdateProvider(c.Context(), id, input)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(provider))
}

type ShippingRateHandler struct {
	rateService services.ShippingRateService
}

func NewShippingRateHandler(rateService services.ShippingRateService) *ShippingRateHandler {
	return &ShippingRateHandler{rateService: rateService}
}

func (h *ShippingRateHandler) Calculate(c fiber.Ctx) error {
	var input services.ShippingRateInput
	if err := c.Bind().Query(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid query parameters")
	}

	if input.CountryID == "" {
		input.CountryID = c.Query("countryId")
	}

	result, err := h.rateService.CalculateRate(c.Context(), input)
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(result))
}

type PriceDestinationHandler struct {
	priceDestinationService services.PriceDestinationService
}

func NewPriceDestinationHandler(priceDestinationService services.PriceDestinationService) *PriceDestinationHandler {
	return &PriceDestinationHandler{priceDestinationService: priceDestinationService}
}

func (h *PriceDestinationHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	prices, total, err := h.priceDestinationService.GetAll(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(prices, page, pageSize, total))
}

func (h *PriceDestinationHandler) Create(c fiber.Ctx) error {
	var input services.CreatePriceDestinationInput
	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	price, err := h.priceDestinationService.Create(c.Context(), input)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(price))
}
