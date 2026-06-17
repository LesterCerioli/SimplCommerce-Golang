package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/pricing/internal/models"
)

type CartRuleRepository interface {
	Create(ctx context.Context, rule *models.CartRule) error
	Update(ctx context.Context, rule *models.CartRule) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.CartRule, error)
	FindAll(ctx context.Context) ([]models.CartRule, error)
	FindActive(ctx context.Context) ([]models.CartRule, error)
	FindByCouponCode(ctx context.Context, code string) (*models.CartRule, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.CartRule, int64, error)
}

type GormCartRuleRepository struct {
	db *gorm.DB
}

func NewCartRuleRepository(db *gorm.DB) CartRuleRepository {
	return &GormCartRuleRepository{db: db}
}

func (r *GormCartRuleRepository) Create(ctx context.Context, rule *models.CartRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *GormCartRuleRepository) Update(ctx context.Context, rule *models.CartRule) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(rule).Error
}

func (r *GormCartRuleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Select("Coupons", "Categories", "Products", "CustomerGroups", "Usages").Delete(&models.CartRule{}, id).Error
}

func (r *GormCartRuleRepository) FindByID(ctx context.Context, id uint) (*models.CartRule, error) {
	var rule models.CartRule
	err := r.db.WithContext(ctx).
		Preload("Coupons").
		Preload("Categories").
		Preload("Products").
		Preload("CustomerGroups").
		Preload("Usages").
		First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *GormCartRuleRepository) FindAll(ctx context.Context) ([]models.CartRule, error) {
	var rules []models.CartRule
	err := r.db.WithContext(ctx).
		Preload("Coupons").
		Preload("Categories").
		Preload("Products").
		Preload("CustomerGroups").
		Find(&rules).Error
	return rules, err
}

func (r *GormCartRuleRepository) FindActive(ctx context.Context) ([]models.CartRule, error) {
	now := time.Now()
	var rules []models.CartRule
	err := r.db.WithContext(ctx).
		Preload("Coupons").
		Preload("Categories").
		Preload("Products").
		Preload("CustomerGroups").
		Where("is_active = ?", true).
		Where("(start_on IS NULL OR start_on <= ?)", now).
		Where("(end_on IS NULL OR end_on >= ?)", now).
		Find(&rules).Error
	return rules, err
}

func (r *GormCartRuleRepository) FindByCouponCode(ctx context.Context, code string) (*models.CartRule, error) {
	var rule models.CartRule
	err := r.db.WithContext(ctx).
		Preload("Coupons", "code = ?", code).
		Preload("Categories").
		Preload("Products").
		Preload("CustomerGroups").
		Joins("JOIN pricing_coupons ON pricing_coupons.cart_rule_id = pricing_cart_rules.id").
		Where("pricing_coupons.code = ?", code).
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *GormCartRuleRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.CartRule, int64, error) {
	var rules []models.CartRule
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CartRule{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Coupons").
		Preload("Categories").
		Preload("Products").
		Preload("CustomerGroups").
		Offset(offset).Limit(pageSize).Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}
