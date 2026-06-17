package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/pricing/internal/models"
	"github.com/simplcommerce-go/services/pricing/internal/repositories"
	"github.com/simplcommerce-go/services/pricing/internal/services"
)

type CouponHandler struct {
	couponRepo     repositories.CouponRepository
	cartRuleRepo   repositories.CartRuleRepository
	pricingService services.PricingService
}

func NewCouponHandler(
	couponRepo repositories.CouponRepository,
	cartRuleRepo repositories.CartRuleRepository,
	pricingService services.PricingService,
) *CouponHandler {
	return &CouponHandler{
		couponRepo:     couponRepo,
		cartRuleRepo:   cartRuleRepo,
		pricingService: pricingService,
	}
}

type createCouponRequest struct {
	Code string `json:"code"`
}

type validateCouponRequest struct {
	Code   string `json:"code"`
	UserID uint   `json:"userId"`
}

func (h *CouponHandler) ListByCartRule(c fiber.Ctx) error {
	cartRuleIDStr := c.Params("id")
	cartRuleID, err := strconv.ParseUint(cartRuleIDStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid cart rule id")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	coupons, total, err := h.couponRepo.PaginateByCartRuleID(c.Context(), uint(cartRuleID), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(coupons, page, pageSize, total))
}

func (h *CouponHandler) Create(c fiber.Ctx) error {
	cartRuleIDStr := c.Params("id")
	cartRuleID, err := strconv.ParseUint(cartRuleIDStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid cart rule id")
	}

	_, err = h.cartRuleRepo.FindByID(c.Context(), uint(cartRuleID))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "cart rule not found")
	}

	var req createCouponRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Code == "" {
		return response.Error(fiber.StatusBadRequest, "coupon code is required")
	}

	coupon := &models.Coupon{
		CartRuleID: uint(cartRuleID),
		Code:       req.Code,
	}

	if err := h.couponRepo.Create(c.Context(), coupon); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(coupon))
}

func (h *CouponHandler) Validate(c fiber.Ctx) error {
	var req validateCouponRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := h.pricingService.ValidateCoupon(c.Context(), req.Code, req.UserID)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(result))
}
