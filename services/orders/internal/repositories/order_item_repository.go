package repositories

import (
	"context"

	"github.com/simplcommerce-go/services/orders/internal/models"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	Create(ctx context.Context, item *models.OrderItem) error
	CreateBatch(ctx context.Context, items []models.OrderItem) error
	FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error)
	DeleteByOrderID(ctx context.Context, orderID uint) error
}

type orderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) Create(ctx context.Context, item *models.OrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *orderItemRepository) CreateBatch(ctx context.Context, items []models.OrderItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *orderItemRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (r *orderItemRepository) DeleteByOrderID(ctx context.Context, orderID uint) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error
}
