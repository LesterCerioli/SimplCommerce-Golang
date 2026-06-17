package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/pricing/internal/models"
)

type CouponRepository interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	Update(ctx context.Context, coupon *models.Coupon) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Coupon, error)
	FindByCode(ctx context.Context, code string) (*models.Coupon, error)
	FindByCartRuleID(ctx context.Context, cartRuleID uint) ([]models.Coupon, error)
	FindAll(ctx context.Context) ([]models.Coupon, error)
	PaginateByCartRuleID(ctx context.Context, cartRuleID uint, page, pageSize int) ([]models.Coupon, int64, error)
}

type GormCouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) CouponRepository {
	return &GormCouponRepository{db: db}
}

func (r *GormCouponRepository) Create(ctx context.Context, coupon *models.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

func (r *GormCouponRepository) Update(ctx context.Context, coupon *models.Coupon) error {
	return r.db.WithContext(ctx).Save(coupon).Error
}

func (r *GormCouponRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Coupon{}, id).Error
}

func (r *GormCouponRepository) FindByID(ctx context.Context, id uint) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.WithContext(ctx).Preload("CartRule").First(&coupon, id).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *GormCouponRepository) FindByCode(ctx context.Context, code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.WithContext(ctx).Preload("CartRule").Where("code = ?", code).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *GormCouponRepository) FindByCartRuleID(ctx context.Context, cartRuleID uint) ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := r.db.WithContext(ctx).Where("cart_rule_id = ?", cartRuleID).Find(&coupons).Error
	return coupons, err
}

func (r *GormCouponRepository) FindAll(ctx context.Context) ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := r.db.WithContext(ctx).Preload("CartRule").Find(&coupons).Error
	return coupons, err
}

func (r *GormCouponRepository) PaginateByCartRuleID(ctx context.Context, cartRuleID uint, page, pageSize int) ([]models.Coupon, int64, error) {
	var coupons []models.Coupon
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Coupon{}).Where("cart_rule_id = ?", cartRuleID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&coupons).Error
	if err != nil {
		return nil, 0, err
	}

	return coupons, total, nil
}
