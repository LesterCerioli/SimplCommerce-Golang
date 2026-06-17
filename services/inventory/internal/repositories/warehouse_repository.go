package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/inventory/internal/models"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *models.Warehouse) error
	Update(ctx context.Context, warehouse *models.Warehouse) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Warehouse, error)
	FindAll(ctx context.Context) ([]models.Warehouse, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Warehouse, int64, error)
}

type GormWarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &GormWarehouseRepository{db: db}
}

func (r *GormWarehouseRepository) Create(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *GormWarehouseRepository) Update(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *GormWarehouseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Warehouse{}, id).Error
}

func (r *GormWarehouseRepository) FindByID(ctx context.Context, id uint) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.WithContext(ctx).First(&warehouse, id).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *GormWarehouseRepository) FindAll(ctx context.Context) ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	err := r.db.WithContext(ctx).Find(&warehouses).Error
	return warehouses, err
}

func (r *GormWarehouseRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Warehouse, int64, error) {
	var warehouses []models.Warehouse
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Warehouse{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&warehouses).Error
	if err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}
