package repositories

import (
	"context"

	"github.com/simplcommerce-go/pkg/repository"
	"github.com/simplcommerce-go/services/cart/internal/models"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	Create(ctx context.Context, item *models.CartItem) error
	Update(ctx context.Context, item *models.CartItem) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.CartItem, error)
	FindByCustomer(ctx context.Context, customerID uint) ([]models.CartItem, error)
	FindByCustomerAndProduct(ctx context.Context, customerID, productID uint) (*models.CartItem, error)
	ClearCustomerCart(ctx context.Context, customerID uint) error
}

type cartItemRepository struct {
	repo repository.Repository[models.CartItem]
	db   *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{
		repo: repository.NewGormRepository[models.CartItem](db),
		db:   db,
	}
}

func (r *cartItemRepository) Create(ctx context.Context, item *models.CartItem) error {
	return r.repo.Create(ctx, item)
}

func (r *cartItemRepository) Update(ctx context.Context, item *models.CartItem) error {
	return r.repo.Update(ctx, item)
}

func (r *cartItemRepository) Delete(ctx context.Context, id uint) error {
	return r.repo.Delete(ctx, id)
}

func (r *cartItemRepository) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *cartItemRepository) FindByCustomer(ctx context.Context, customerID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&items).Error
	return items, err
}

func (r *cartItemRepository) FindByCustomerAndProduct(ctx context.Context, customerID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).Where("customer_id = ? AND product_id = ?", customerID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepository) ClearCustomerCart(ctx context.Context, customerID uint) error {
	return r.db.WithContext(ctx).Where("customer_id = ?", customerID).Delete(&models.CartItem{}).Error
}
