package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
	"simplcommerce/services/implementations"
)

type OrderController struct {
	svc *initializers.Services
}

func NewOrderController(svc *initializers.Services) *OrderController {
	return &OrderController{svc: svc}
}

func (h *OrderController) ListOrders(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	roles, ok := c.Locals("roles").([]string)
	isAdmin := ok && hasRole(roles, "admin")

	if isAdmin {
		orders, total, err := h.svc.OrderService.GetAllOrders(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true, "data": orders, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
	}

	customerID := c.Locals("userID").(string)
	orders, total, err := h.svc.OrderService.GetCustomerOrders(c.Context(), customerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": orders, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}

func (h *OrderController) GetOrder(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)
	id := c.Params("id")

	order, err := h.svc.OrderService.GetOrderByID(c.Context(), id, customerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": order})
}

func (h *OrderController) CreateOrder(c fiber.Ctx) error {
	customerID := c.Locals("userID").(string)

	var req struct {
		Items             []implementations.CreateOrderItemRequest `json:"items"`
		ShippingAddressID string                                 `json:"shippingAddressId"`
		BillingAddressID  string                                 `json:"billingAddressId"`
		CouponCode        string                                 `json:"couponCode"`
		CouponRuleName    string                                 `json:"couponRuleName"`
		DiscountAmount    float64                                `json:"discountAmount"`
		OrderNote         string                                 `json:"orderNote"`
		ShippingMethod    string                                 `json:"shippingMethod"`
		ShippingFeeAmount float64                                `json:"shippingFeeAmount"`
		TaxAmount         float64                                `json:"taxAmount"`
		PaymentMethod     string                                 `json:"paymentMethod"`
		PaymentFeeAmount  float64                                `json:"paymentFeeAmount"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	var subTotal float64
	for _, item := range req.Items {
		subTotal += item.ProductPrice * float64(item.Quantity)
	}
	subTotalWithDiscount := subTotal - req.DiscountAmount

	createReq := implementations.CreateOrderRequest{
		Items:                req.Items,
		ShippingAddressID:    req.ShippingAddressID,
		BillingAddressID:     req.BillingAddressID,
		CouponCode:           req.CouponCode,
		CouponRuleName:       req.CouponRuleName,
		DiscountAmount:       req.DiscountAmount,
		SubTotal:             subTotal,
		SubTotalWithDiscount: subTotalWithDiscount,
		OrderNote:            req.OrderNote,
		ShippingMethod:       req.ShippingMethod,
		ShippingFeeAmount:    req.ShippingFeeAmount,
		TaxAmount:            req.TaxAmount,
		PaymentMethod:        req.PaymentMethod,
		PaymentFeeAmount:     req.PaymentFeeAmount,
	}

	order, err := h.svc.OrderService.CreateOrder(c.Context(), customerID, createReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": order})
}

func (h *OrderController) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")
	updatedByID := c.Locals("userID").(string)

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "status is required"})
	}

	if err := h.svc.OrderService.UpdateStatus(c.Context(), id, req.Status, updatedByID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
