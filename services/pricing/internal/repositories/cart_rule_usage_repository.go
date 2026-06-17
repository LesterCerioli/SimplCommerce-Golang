package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/pricing/internal/models"
)

type CartRuleUsageRepository interface {
	Create(ctx context.Context, usage *models.CartRuleUsage) error
	FindByCartRuleID(ctx context.Context, cartRuleID uint) ([]models.CartRuleUsage, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.CartRuleUsage, error)
	FindByOrderID(ctx context.Context, orderID uint) ([]models.CartRuleUsage, error)
	CountByCartRuleID(ctx context.Context, cartRuleID uint) (int64, error)
	CountByCouponID(ctx context.Context, couponID uint) (int64, error)
	CountByUserAndCartRule(ctx context.Context, userID, cartRuleID uint) (int64, error)
}

type GormCartRuleUsageRepository struct {
	db *gorm.DB
}

func NewCartRuleUsageRepository(db *gorm.DB) CartRuleUsageRepository {
	return &GormCartRuleUsageRepository{db: db}
}

func (r *GormCartRuleUsageRepository) Create(ctx context.Context, usage *models.CartRuleUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *GormCartRuleUsageRepository) FindByCartRuleID(ctx context.Context, cartRuleID uint) ([]models.CartRuleUsage, error) {
	var usages []models.CartRuleUsage
	err := r.db.WithContext(ctx).Preload("Coupon").Where("cart_rule_id = ?", cartRuleID).Find(&usages).Error
	return usages, err
}

func (r *GormCartRuleUsageRepository) FindByUserID(ctx context.Context, userID uint) ([]models.CartRuleUsage, error) {
	var usages []models.CartRuleUsage
	err := r.db.WithContext(ctx).Preload("CartRule").Preload("Coupon").Where("user_id = ?", userID).Find(&usages).Error
	return usages, err
}

func (r *GormCartRuleUsageRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.CartRuleUsage, error) {
	var usages []models.CartRuleUsage
	err := r.db.WithContext(ctx).Preload("CartRule").Preload("Coupon").Where("order_id = ?", orderID).Find(&usages).Error
	return usages, err
}

func (r *GormCartRuleUsageRepository) CountByCartRuleID(ctx context.Context, cartRuleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CartRuleUsage{}).Where("cart_rule_id = ?", cartRuleID).Count(&count).Error
	return count, err
}

func (r *GormCartRuleUsageRepository) CountByCouponID(ctx context.Context, couponID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CartRuleUsage{}).Where("coupon_id = ?", couponID).Count(&count).Error
	return count, err
}

func (r *GormCartRuleUsageRepository) CountByUserAndCartRule(ctx context.Context, userID, cartRuleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CartRuleUsage{}).Where("user_id = ? AND cart_rule_id = ?", userID, cartRuleID).Count(&count).Error
	return count, err
}
