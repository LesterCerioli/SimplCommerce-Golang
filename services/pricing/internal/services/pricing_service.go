package services

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/pricing/internal/models"
	"github.com/simplcommerce-go/services/pricing/internal/repositories"
)

type CouponValidationResult struct {
	Valid     bool        `json:"valid"`
	CartRule  *models.CartRule `json:"cartRule,omitempty"`
	Error     string      `json:"error,omitempty"`
}

type CartRuleDiscount struct {
	CartRuleID     uint    `json:"cartRuleId"`
	RuleName       string  `json:"ruleName"`
	RuleToApply    string  `json:"ruleToApply"`
	DiscountAmount float64 `json:"discountAmount"`
}

type PricingService interface {
	ValidateCoupon(ctx context.Context, code string, userID uint) (*CouponValidationResult, error)
	ApplyCartRules(ctx context.Context, cartItems []CartItem, couponCode string, userID uint) ([]CartRuleDiscount, error)
	GetCatalogRuleDiscounts(ctx context.Context, productIDs []uint, customerGroupIDs []uint) (map[uint]float64, error)
}

type CartItem struct {
	ProductID  uint
	Quantity   int
	UnitPrice  float64
	CategoryID uint
}

type pricingService struct {
	cartRuleRepo    repositories.CartRuleRepository
	couponRepo      repositories.CouponRepository
	catalogRuleRepo repositories.CatalogRuleRepository
	usageRepo       repositories.CartRuleUsageRepository
	db              *gorm.DB
}

func NewPricingService(
	cartRuleRepo repositories.CartRuleRepository,
	couponRepo repositories.CouponRepository,
	catalogRuleRepo repositories.CatalogRuleRepository,
	usageRepo repositories.CartRuleUsageRepository,
	db *gorm.DB,
) PricingService {
	return &pricingService{
		cartRuleRepo:    cartRuleRepo,
		couponRepo:      couponRepo,
		catalogRuleRepo: catalogRuleRepo,
		usageRepo:       usageRepo,
		db:              db,
	}
}

func (s *pricingService) ValidateCoupon(ctx context.Context, code string, userID uint) (*CouponValidationResult, error) {
	if code == "" {
		return &CouponValidationResult{Valid: false, Error: "coupon code is required"}, nil
	}

	coupon, err := s.couponRepo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CouponValidationResult{Valid: false, Error: "coupon not found"}, nil
		}
		return nil, err
	}

	cartRule, err := s.cartRuleRepo.FindByID(ctx, coupon.CartRuleID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if !cartRule.IsActive {
		return &CouponValidationResult{Valid: false, Error: "coupon is not active"}, nil
	}

	if cartRule.StartOn != nil && now.Before(*cartRule.StartOn) {
		return &CouponValidationResult{Valid: false, Error: "coupon is not yet valid"}, nil
	}

	if cartRule.EndOn != nil && now.After(*cartRule.EndOn) {
		return &CouponValidationResult{Valid: false, Error: "coupon has expired"}, nil
	}

	if cartRule.UsageLimitPerCoupon != nil {
		usageCount, err := s.usageRepo.CountByCouponID(ctx, coupon.ID)
		if err != nil {
			return nil, err
		}
		if int(usageCount) >= *cartRule.UsageLimitPerCoupon {
			return &CouponValidationResult{Valid: false, Error: "coupon usage limit reached"}, nil
		}
	}

	if cartRule.UsageLimitPerCustomer != nil && userID != 0 {
		usageCount, err := s.usageRepo.CountByUserAndCartRule(ctx, userID, cartRule.ID)
		if err != nil {
			return nil, err
		}
		if int(usageCount) >= *cartRule.UsageLimitPerCustomer {
			return &CouponValidationResult{Valid: false, Error: "customer usage limit reached"}, nil
		}
	}

	return &CouponValidationResult{
		Valid:    true,
		CartRule: cartRule,
	}, nil
}

func (s *pricingService) ApplyCartRules(ctx context.Context, cartItems []CartItem, couponCode string, userID uint) ([]CartRuleDiscount, error) {
	var discounts []CartRuleDiscount

	rules, err := s.cartRuleRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	if couponCode != "" {
		validation, err := s.ValidateCoupon(ctx, couponCode, userID)
		if err != nil {
			return nil, err
		}
		if validation.Valid && validation.CartRule != nil {
			rule := validation.CartRule
			if s.isRuleApplicable(ctx, rule, cartItems) {
				discount := s.calculateDiscount(rule, cartItems)
				discounts = append(discounts, discount)
			}
		}
	}

	for _, rule := range rules {
		if rule.IsCouponRequired {
			continue
		}
		if s.isRuleApplicable(ctx, &rule, cartItems) {
			discount := s.calculateDiscount(&rule, cartItems)
			discounts = append(discounts, discount)
		}
	}

	return discounts, nil
}

func (s *pricingService) isRuleApplicable(ctx context.Context, rule *models.CartRule, cartItems []CartItem) bool {
	if len(rule.Products) > 0 {
		productIDs := make(map[uint]bool)
		for _, p := range rule.Products {
			productIDs[p.ProductID] = true
		}
		for _, item := range cartItems {
			if productIDs[item.ProductID] {
				return true
			}
		}
	}

	if len(rule.Categories) > 0 {
		categoryIDs := make(map[uint]bool)
		for _, c := range rule.Categories {
			categoryIDs[c.CategoryID] = true
		}
		for _, item := range cartItems {
			if categoryIDs[item.CategoryID] {
				return true
			}
		}
	}

	if len(rule.Products) == 0 && len(rule.Categories) == 0 {
		return true
	}

	return false
}

func (s *pricingService) calculateDiscount(rule *models.CartRule, cartItems []CartItem) CartRuleDiscount {
	switch rule.RuleToApply {
	case "by_percent":
		total := 0.0
		for _, item := range cartItems {
			total += float64(item.Quantity) * item.UnitPrice
		}
		discount := total * rule.DiscountAmount / 100.0
		if rule.MaxDiscountAmount != nil && discount > *rule.MaxDiscountAmount {
			discount = *rule.MaxDiscountAmount
		}
		return CartRuleDiscount{
			CartRuleID:     rule.ID,
			RuleName:       rule.Name,
			RuleToApply:    rule.RuleToApply,
			DiscountAmount: discount,
		}
	case "by_fixed":
		return CartRuleDiscount{
			CartRuleID:     rule.ID,
			RuleName:       rule.Name,
			RuleToApply:    rule.RuleToApply,
			DiscountAmount: rule.DiscountAmount,
		}
	default:
		return CartRuleDiscount{
			CartRuleID:     rule.ID,
			RuleName:       rule.Name,
			RuleToApply:    rule.RuleToApply,
			DiscountAmount: rule.DiscountAmount,
		}
	}
}

func (s *pricingService) GetCatalogRuleDiscounts(ctx context.Context, productIDs []uint, customerGroupIDs []uint) (map[uint]float64, error) {
	discounts := make(map[uint]float64)

	rules, err := s.catalogRuleRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return discounts, nil
	}

	customerGroupSet := make(map[uint]bool)
	for _, id := range customerGroupIDs {
		customerGroupSet[id] = true
	}

	for _, rule := range rules {
		if len(rule.CustomerGroups) > 0 {
			hasGroup := false
			for _, cg := range rule.CustomerGroups {
				if customerGroupSet[cg.CustomerGroupID] {
					hasGroup = true
					break
				}
			}
			if !hasGroup {
				continue
			}
		}

		for _, pid := range productIDs {
			existing, ok := discounts[pid]
			discountVal := rule.DiscountAmount
			if !ok || discountVal > existing {
				discounts[pid] = discountVal
			}
		}
	}

	return discounts, nil
}
