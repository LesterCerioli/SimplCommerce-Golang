package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/inventory/internal/models"
)

type StockRepository interface {
	Create(ctx context.Context, stock *models.Stock) error
	Update(ctx context.Context, stock *models.Stock) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Stock, error)
	FindByProductID(ctx context.Context, productID uint) ([]models.Stock, error)
	FindByProductAndWarehouse(ctx context.Context, productID, warehouseID uint) (*models.Stock, error)
	FindAll(ctx context.Context) ([]models.Stock, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Stock, int64, error)
}

type GormStockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &GormStockRepository{db: db}
}

func (r *GormStockRepository) Create(ctx context.Context, stock *models.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

func (r *GormStockRepository) Update(ctx context.Context, stock *models.Stock) error {
	return r.db.WithContext(ctx).Save(stock).Error
}

func (r *GormStockRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Stock{}, id).Error
}

func (r *GormStockRepository) FindByID(ctx context.Context, id uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.WithContext(ctx).Preload("Warehouse").First(&stock, id).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *GormStockRepository) FindByProductID(ctx context.Context, productID uint) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.WithContext(ctx).Preload("Warehouse").Where("product_id = ?", productID).Find(&stocks).Error
	return stocks, err
}

func (r *GormStockRepository) FindByProductAndWarehouse(ctx context.Context, productID, warehouseID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.WithContext(ctx).Preload("Warehouse").Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *GormStockRepository) FindAll(ctx context.Context) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.WithContext(ctx).Preload("Warehouse").Find(&stocks).Error
	return stocks, err
}

func (r *GormStockRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Stock{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Warehouse").Offset(offset).Limit(pageSize).Find(&stocks).Error
	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}
