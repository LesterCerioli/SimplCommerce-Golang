package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/inventory/internal/models"
)

type StockHistoryRepository interface {
	Create(ctx context.Context, history *models.StockHistory) error
	FindByProductID(ctx context.Context, productID uint) ([]models.StockHistory, error)
	FindAll(ctx context.Context) ([]models.StockHistory, error)
	PaginateByProductID(ctx context.Context, productID uint, page, pageSize int) ([]models.StockHistory, int64, error)
}

type GormStockHistoryRepository struct {
	db *gorm.DB
}

func NewStockHistoryRepository(db *gorm.DB) StockHistoryRepository {
	return &GormStockHistoryRepository{db: db}
}

func (r *GormStockHistoryRepository) Create(ctx context.Context, history *models.StockHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *GormStockHistoryRepository) FindByProductID(ctx context.Context, productID uint) ([]models.StockHistory, error) {
	var histories []models.StockHistory
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at DESC").Find(&histories).Error
	return histories, err
}

func (r *GormStockHistoryRepository) FindAll(ctx context.Context) ([]models.StockHistory, error) {
	var histories []models.StockHistory
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&histories).Error
	return histories, err
}

func (r *GormStockHistoryRepository) PaginateByProductID(ctx context.Context, productID uint, page, pageSize int) ([]models.StockHistory, int64, error) {
	var histories []models.StockHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StockHistory{}).Where("product_id = ?", productID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
