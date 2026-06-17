package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/payment/internal/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, entity *models.Payment) error
	Update(ctx context.Context, entity *models.Payment) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Payment, error)
	FindAll(ctx context.Context) ([]models.Payment, error)
	FindByOrderID(ctx context.Context, orderID uint) ([]models.Payment, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Payment, int64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, entity *models.Payment) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *paymentRepository) Update(ctx context.Context, entity *models.Payment) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Payment{}, id).Error
}

func (r *paymentRepository) FindByID(ctx context.Context, id uint) (*models.Payment, error) {
	var entity models.Payment
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *paymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	var entities []models.Payment
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.Payment, error) {
	var entities []models.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&entities).Error
	return entities, err
}

func (r *paymentRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Payment, int64, error) {
	var entities []models.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Payment{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

type PaymentProviderRepository interface {
	Create(ctx context.Context, entity *models.PaymentProvider) error
	Update(ctx context.Context, entity *models.PaymentProvider) error
	FindByID(ctx context.Context, id string) (*models.PaymentProvider, error)
	FindAll(ctx context.Context) ([]models.PaymentProvider, error)
	FindEnabled(ctx context.Context) ([]models.PaymentProvider, error)
}

type paymentProviderRepository struct {
	db *gorm.DB
}

func NewPaymentProviderRepository(db *gorm.DB) PaymentProviderRepository {
	return &paymentProviderRepository{db: db}
}

func (r *paymentProviderRepository) Create(ctx context.Context, entity *models.PaymentProvider) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *paymentProviderRepository) Update(ctx context.Context, entity *models.PaymentProvider) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *paymentProviderRepository) FindByID(ctx context.Context, id string) (*models.PaymentProvider, error) {
	var entity models.PaymentProvider
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *paymentProviderRepository) FindAll(ctx context.Context) ([]models.PaymentProvider, error) {
	var entities []models.PaymentProvider
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *paymentProviderRepository) FindEnabled(ctx context.Context) ([]models.PaymentProvider, error) {
	var entities []models.PaymentProvider
	err := r.db.WithContext(ctx).Where("is_enabled = ?", true).Find(&entities).Error
	return entities, err
}
