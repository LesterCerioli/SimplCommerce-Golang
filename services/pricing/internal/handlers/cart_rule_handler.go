package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/pricing/internal/models"
	"github.com/simplcommerce-go/services/pricing/internal/repositories"
)

type CartRuleHandler struct {
	cartRuleRepo repositories.CartRuleRepository
}

func NewCartRuleHandler(cartRuleRepo repositories.CartRuleRepository) *CartRuleHandler {
	return &CartRuleHandler{cartRuleRepo: cartRuleRepo}
}

type createCartRuleRequest struct {
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	IsActive             *bool      `json:"isActive"`
	StartOn              *time.Time `json:"startOn"`
	EndOn                *time.Time `json:"endOn"`
	IsCouponRequired     *bool      `json:"isCouponRequired"`
	RuleToApply          string     `json:"ruleToApply"`
	DiscountAmount       float64    `json:"discountAmount"`
	MaxDiscountAmount    *float64   `json:"maxDiscountAmount"`
	DiscountStep         *int       `json:"discountStep"`
	UsageLimitPerCoupon  *int       `json:"usageLimitPerCoupon"`
	UsageLimitPerCustomer *int      `json:"usageLimitPerCustomer"`
	CategoryIDs          []uint     `json:"categoryIds"`
	ProductIDs           []uint     `json:"productIds"`
	CustomerGroupIDs     []uint     `json:"customerGroupIds"`
}

type updateCartRuleRequest struct {
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	IsActive             *bool      `json:"isActive"`
	StartOn              *time.Time `json:"startOn"`
	EndOn                *time.Time `json:"endOn"`
	IsCouponRequired     *bool      `json:"isCouponRequired"`
	RuleToApply          string     `json:"ruleToApply"`
	DiscountAmount       float64    `json:"discountAmount"`
	MaxDiscountAmount    *float64   `json:"maxDiscountAmount"`
	DiscountStep         *int       `json:"discountStep"`
	UsageLimitPerCoupon  *int       `json:"usageLimitPerCoupon"`
	UsageLimitPerCustomer *int      `json:"usageLimitPerCustomer"`
	CategoryIDs          []uint     `json:"categoryIds"`
	ProductIDs           []uint     `json:"productIds"`
	CustomerGroupIDs     []uint     `json:"customerGroupIds"`
}

func (h *CartRuleHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	rules, total, err := h.cartRuleRepo.Paginate(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(rules, page, pageSize, total))
}

func (h *CartRuleHandler) Create(c fiber.Ctx) error {
	var req createCartRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isCouponRequired := false
	if req.IsCouponRequired != nil {
		isCouponRequired = *req.IsCouponRequired
	}

	ruleToApply := req.RuleToApply
	if ruleToApply == "" {
		ruleToApply = "by_fixed"
	}

	rule := &models.CartRule{
		Name:                 req.Name,
		Description:          req.Description,
		IsActive:             isActive,
		StartOn:              req.StartOn,
		EndOn:                req.EndOn,
		IsCouponRequired:     isCouponRequired,
		RuleToApply:          ruleToApply,
		DiscountAmount:       req.DiscountAmount,
		MaxDiscountAmount:    req.MaxDiscountAmount,
		DiscountStep:         req.DiscountStep,
		UsageLimitPerCoupon:  req.UsageLimitPerCoupon,
		UsageLimitPerCustomer: req.UsageLimitPerCustomer,
	}

	for _, catID := range req.CategoryIDs {
		rule.Categories = append(rule.Categories, models.CartRuleCategory{
			CartRuleID: 0,
			CategoryID: catID,
		})
	}

	for _, prodID := range req.ProductIDs {
		rule.Products = append(rule.Products, models.CartRuleProduct{
			CartRuleID: 0,
			ProductID:  prodID,
		})
	}

	for _, cgID := range req.CustomerGroupIDs {
		rule.CustomerGroups = append(rule.CustomerGroups, models.CartRuleCustomerGroup{
			CartRuleID:      0,
			CustomerGroupID: cgID,
		})
	}

	if err := h.cartRuleRepo.Create(c.Context(), rule); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(rule))
}

func (h *CartRuleHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid cart rule id")
	}

	var req updateCartRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	rule, err := h.cartRuleRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "cart rule not found")
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	rule.Description = req.Description
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	rule.StartOn = req.StartOn
	rule.EndOn = req.EndOn
	if req.IsCouponRequired != nil {
		rule.IsCouponRequired = *req.IsCouponRequired
	}
	if req.RuleToApply != "" {
		rule.RuleToApply = req.RuleToApply
	}
	rule.DiscountAmount = req.DiscountAmount
	rule.MaxDiscountAmount = req.MaxDiscountAmount
	rule.DiscountStep = req.DiscountStep
	rule.UsageLimitPerCoupon = req.UsageLimitPerCoupon
	rule.UsageLimitPerCustomer = req.UsageLimitPerCustomer

	rule.Categories = nil
	for _, catID := range req.CategoryIDs {
		rule.Categories = append(rule.Categories, models.CartRuleCategory{
			CartRuleID: rule.ID,
			CategoryID: catID,
		})
	}

	rule.Products = nil
	for _, prodID := range req.ProductIDs {
		rule.Products = append(rule.Products, models.CartRuleProduct{
			CartRuleID: rule.ID,
			ProductID:  prodID,
		})
	}

	rule.CustomerGroups = nil
	for _, cgID := range req.CustomerGroupIDs {
		rule.CustomerGroups = append(rule.CustomerGroups, models.CartRuleCustomerGroup{
			CartRuleID:      rule.ID,
			CustomerGroupID: cgID,
		})
	}

	if err := h.cartRuleRepo.Update(c.Context(), rule); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(rule))
}

func (h *CartRuleHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid cart rule id")
	}

	if err := h.cartRuleRepo.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
