package repositories

import (
	"context"

	"github.com/simplcommerce-go/pkg/repository"
	"github.com/simplcommerce-go/services/orders/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Order, error)
	FindByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]models.Order, int64, error)
	FindAllPaginated(ctx context.Context, page, pageSize int) ([]models.Order, int64, error)
	FindWithItems(ctx context.Context, id uint) (*models.Order, error)
	UpdateStatus(ctx context.Context, id uint, status string, updatedByID uint) error
}

type orderRepository struct {
	repo repository.Repository[models.Order]
	db   *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		repo: repository.NewGormRepository[models.Order](db),
		db:   db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.repo.Create(ctx, order)
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.repo.Update(ctx, order)
}

func (r *orderRepository) Delete(ctx context.Context, id uint) error {
	return r.repo.Delete(ctx, id)
}

func (r *orderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *orderRepository) FindByCustomerID(ctx context.Context, customerID uint, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("customer_id = ?", customerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Preload("Items").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindAllPaginated(ctx context.Context, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Preload("Items").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindWithItems(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Preload("Items").Preload("Histories").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status string, updatedByID uint) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"order_status": status,
		"updated_by_id": updatedByID,
	}).Error
}
