package repositories

import (
	"context"

	"github.com/simplcommerce-go/services/orders/internal/models"
	"gorm.io/gorm"
)

type OrderHistoryRepository interface {
	Create(ctx context.Context, history *models.OrderHistory) error
	FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderHistory, error)
}

type orderHistoryRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) OrderHistoryRepository {
	return &orderHistoryRepository{db: db}
}

func (r *orderHistoryRepository) Create(ctx context.Context, history *models.OrderHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *orderHistoryRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderHistory, error) {
	var histories []models.OrderHistory
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at ASC").Find(&histories).Error
	return histories, err
}
