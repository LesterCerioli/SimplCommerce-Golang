package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/orders/internal/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type createOrderRequest struct {
	Items             []services.CreateOrderItemRequest `json:"items"`
	ShippingAddressID uint                              `json:"shippingAddressId"`
	BillingAddressID  uint                              `json:"billingAddressId"`
	CouponCode        string                            `json:"couponCode"`
	CouponRuleName    string                            `json:"couponRuleName"`
	DiscountAmount    float64                           `json:"discountAmount"`
	OrderNote         string                            `json:"orderNote"`
	ShippingMethod    string                            `json:"shippingMethod"`
	ShippingFeeAmount float64                           `json:"shippingFeeAmount"`
	TaxAmount         float64                           `json:"taxAmount"`
	PaymentMethod     string                            `json:"paymentMethod"`
	PaymentFeeAmount  float64                           `json:"paymentFeeAmount"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	var req createOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	orderReq := services.CreateOrderRequest{
		Items:             req.Items,
		ShippingAddressID: req.ShippingAddressID,
		BillingAddressID:  req.BillingAddressID,
		CouponCode:        req.CouponCode,
		CouponRuleName:    req.CouponRuleName,
		DiscountAmount:    req.DiscountAmount,
		OrderNote:         req.OrderNote,
		ShippingMethod:    req.ShippingMethod,
		ShippingFeeAmount: req.ShippingFeeAmount,
		TaxAmount:         req.TaxAmount,
		PaymentMethod:     req.PaymentMethod,
		PaymentFeeAmount:  req.PaymentFeeAmount,
	}

	order, err := h.orderService.CreateOrder(c.Context(), customerID, orderReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(order))
}

func (h *OrderHandler) GetOrderByID(c fiber.Ctx) error {
	customerID := c.Locals("userID").(uint)

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid order id",
		})
	}

	order, err := h.orderService.GetOrderByID(c.Context(), uint(id), customerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(order))
}

func (h *OrderHandler) ListOrders(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	roles, ok := c.Locals("roles").([]string)
	isAdmin := ok && hasRole(roles, "admin")

	if isAdmin {
		o, t, err := h.orderService.GetAllOrders(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(response.Paginated(o, page, pageSize, t))
	}

	customerID := c.Locals("userID").(uint)
	o, t, err := h.orderService.GetCustomerOrders(c.Context(), customerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Paginated(o, page, pageSize, t))
}

func (h *OrderHandler) UpdateStatus(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid order id",
		})
	}

	var req updateStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "status is required",
		})
	}

	updatedByID := c.Locals("userID").(uint)

	if err := h.orderService.UpdateStatus(c.Context(), uint(id), req.Status, updatedByID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(response.Success(nil))
}

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
