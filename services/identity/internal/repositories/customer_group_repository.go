package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type CustomerGroupRepository interface {
	Create(ctx context.Context, group *models.CustomerGroup) error
	Update(ctx context.Context, group *models.CustomerGroup) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.CustomerGroup, error)
	FindAll(ctx context.Context) ([]models.CustomerGroup, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.CustomerGroup, int64, error)
}

type GormCustomerGroupRepository struct {
	db *gorm.DB
}

func NewCustomerGroupRepository(db *gorm.DB) CustomerGroupRepository {
	return &GormCustomerGroupRepository{db: db}
}

func (r *GormCustomerGroupRepository) Create(ctx context.Context, group *models.CustomerGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GormCustomerGroupRepository) Update(ctx context.Context, group *models.CustomerGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *GormCustomerGroupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CustomerGroup{}, id).Error
}

func (r *GormCustomerGroupRepository) FindByID(ctx context.Context, id uint) (*models.CustomerGroup, error) {
	var group models.CustomerGroup
	err := r.db.WithContext(ctx).Preload("Users").First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GormCustomerGroupRepository) FindAll(ctx context.Context) ([]models.CustomerGroup, error) {
	var groups []models.CustomerGroup
	err := r.db.WithContext(ctx).Find(&groups).Error
	return groups, err
}

func (r *GormCustomerGroupRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.CustomerGroup, int64, error) {
	var groups []models.CustomerGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CustomerGroup{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}
