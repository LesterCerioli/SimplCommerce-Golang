package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/inventory/internal/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.ProductBackInStockSubscription) error
	FindByProductID(ctx context.Context, productID uint) ([]models.ProductBackInStockSubscription, error)
	FindByProductAndEmail(ctx context.Context, productID uint, email string) (*models.ProductBackInStockSubscription, error)
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]models.ProductBackInStockSubscription, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.ProductBackInStockSubscription, int64, error)
}

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) Create(ctx context.Context, sub *models.ProductBackInStockSubscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *GormSubscriptionRepository) FindByProductID(ctx context.Context, productID uint) ([]models.ProductBackInStockSubscription, error) {
	var subs []models.ProductBackInStockSubscription
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&subs).Error
	return subs, err
}

func (r *GormSubscriptionRepository) FindByProductAndEmail(ctx context.Context, productID uint, email string) (*models.ProductBackInStockSubscription, error) {
	var sub models.ProductBackInStockSubscription
	err := r.db.WithContext(ctx).Where("product_id = ? AND customer_email = ?", productID, email).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *GormSubscriptionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductBackInStockSubscription{}, id).Error
}

func (r *GormSubscriptionRepository) FindAll(ctx context.Context) ([]models.ProductBackInStockSubscription, error) {
	var subs []models.ProductBackInStockSubscription
	err := r.db.WithContext(ctx).Find(&subs).Error
	return subs, err
}

func (r *GormSubscriptionRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.ProductBackInStockSubscription, int64, error) {
	var subs []models.ProductBackInStockSubscription
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ProductBackInStockSubscription{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&subs).Error
	if err != nil {
		return nil, 0, err
	}

	return subs, total, nil
}
