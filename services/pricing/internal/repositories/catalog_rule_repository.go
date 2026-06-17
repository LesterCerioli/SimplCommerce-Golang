package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/pricing/internal/models"
)

type CatalogRuleRepository interface {
	Create(ctx context.Context, rule *models.CatalogRule) error
	Update(ctx context.Context, rule *models.CatalogRule) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.CatalogRule, error)
	FindAll(ctx context.Context) ([]models.CatalogRule, error)
	FindActive(ctx context.Context) ([]models.CatalogRule, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.CatalogRule, int64, error)
}

type GormCatalogRuleRepository struct {
	db *gorm.DB
}

func NewCatalogRuleRepository(db *gorm.DB) CatalogRuleRepository {
	return &GormCatalogRuleRepository{db: db}
}

func (r *GormCatalogRuleRepository) Create(ctx context.Context, rule *models.CatalogRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *GormCatalogRuleRepository) Update(ctx context.Context, rule *models.CatalogRule) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(rule).Error
}

func (r *GormCatalogRuleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Select("CustomerGroups").Delete(&models.CatalogRule{}, id).Error
}

func (r *GormCatalogRuleRepository) FindByID(ctx context.Context, id uint) (*models.CatalogRule, error) {
	var rule models.CatalogRule
	err := r.db.WithContext(ctx).
		Preload("CustomerGroups").
		First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *GormCatalogRuleRepository) FindAll(ctx context.Context) ([]models.CatalogRule, error) {
	var rules []models.CatalogRule
	err := r.db.WithContext(ctx).
		Preload("CustomerGroups").
		Find(&rules).Error
	return rules, err
}

func (r *GormCatalogRuleRepository) FindActive(ctx context.Context) ([]models.CatalogRule, error) {
	now := time.Now()
	var rules []models.CatalogRule
	err := r.db.WithContext(ctx).
		Preload("CustomerGroups").
		Where("is_active = ?", true).
		Where("(start_on IS NULL OR start_on <= ?)", now).
		Where("(end_on IS NULL OR end_on >= ?)", now).
		Find(&rules).Error
	return rules, err
}

func (r *GormCatalogRuleRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.CatalogRule, int64, error) {
	var rules []models.CatalogRule
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CatalogRule{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("CustomerGroups").
		Offset(offset).Limit(pageSize).Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}
