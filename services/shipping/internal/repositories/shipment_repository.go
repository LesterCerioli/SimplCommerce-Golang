package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/shipping/internal/models"
)

type ShipmentRepository interface {
	Create(ctx context.Context, entity *models.Shipment) error
	Update(ctx context.Context, entity *models.Shipment) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Shipment, error)
	FindAll(ctx context.Context) ([]models.Shipment, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Shipment, int64, error)
}

type shipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(ctx context.Context, entity *models.Shipment) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *shipmentRepository) Update(ctx context.Context, entity *models.Shipment) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *shipmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Shipment{}, id).Error
}

func (r *shipmentRepository) FindByID(ctx context.Context, id uint) (*models.Shipment, error) {
	var entity models.Shipment
	err := r.db.WithContext(ctx).Preload("Items").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *shipmentRepository) FindAll(ctx context.Context) ([]models.Shipment, error) {
	var entities []models.Shipment
	err := r.db.WithContext(ctx).Preload("Items").Find(&entities).Error
	return entities, err
}

func (r *shipmentRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Shipment, int64, error) {
	var entities []models.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Shipment{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Items").Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

type ShipmentItemRepository interface {
	Create(ctx context.Context, entity *models.ShipmentItem) error
	CreateBatch(ctx context.Context, entities []models.ShipmentItem) error
	FindByShipmentID(ctx context.Context, shipmentID uint) ([]models.ShipmentItem, error)
	DeleteByShipmentID(ctx context.Context, shipmentID uint) error
}

type shipmentItemRepository struct {
	db *gorm.DB
}

func NewShipmentItemRepository(db *gorm.DB) ShipmentItemRepository {
	return &shipmentItemRepository{db: db}
}

func (r *shipmentItemRepository) Create(ctx context.Context, entity *models.ShipmentItem) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *shipmentItemRepository) CreateBatch(ctx context.Context, entities []models.ShipmentItem) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

func (r *shipmentItemRepository) FindByShipmentID(ctx context.Context, shipmentID uint) ([]models.ShipmentItem, error) {
	var entities []models.ShipmentItem
	err := r.db.WithContext(ctx).Where("shipment_id = ?", shipmentID).Find(&entities).Error
	return entities, err
}

func (r *shipmentItemRepository) DeleteByShipmentID(ctx context.Context, shipmentID uint) error {
	return r.db.WithContext(ctx).Where("shipment_id = ?", shipmentID).Delete(&models.ShipmentItem{}).Error
}

type ShippingProviderRepository interface {
	Create(ctx context.Context, entity *models.ShippingProvider) error
	Update(ctx context.Context, entity *models.ShippingProvider) error
	FindByID(ctx context.Context, id string) (*models.ShippingProvider, error)
	FindAll(ctx context.Context) ([]models.ShippingProvider, error)
	FindEnabled(ctx context.Context) ([]models.ShippingProvider, error)
}

type shippingProviderRepository struct {
	db *gorm.DB
}

func NewShippingProviderRepository(db *gorm.DB) ShippingProviderRepository {
	return &shippingProviderRepository{db: db}
}

func (r *shippingProviderRepository) Create(ctx context.Context, entity *models.ShippingProvider) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *shippingProviderRepository) Update(ctx context.Context, entity *models.ShippingProvider) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *shippingProviderRepository) FindByID(ctx context.Context, id string) (*models.ShippingProvider, error) {
	var entity models.ShippingProvider
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *shippingProviderRepository) FindAll(ctx context.Context) ([]models.ShippingProvider, error) {
	var entities []models.ShippingProvider
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *shippingProviderRepository) FindEnabled(ctx context.Context) ([]models.ShippingProvider, error) {
	var entities []models.ShippingProvider
	err := r.db.WithContext(ctx).Where("is_enabled = ?", true).Find(&entities).Error
	return entities, err
}

type PriceAndDestinationRepository interface {
	Create(ctx context.Context, entity *models.PriceAndDestination) error
	Update(ctx context.Context, entity *models.PriceAndDestination) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.PriceAndDestination, error)
	FindAll(ctx context.Context) ([]models.PriceAndDestination, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.PriceAndDestination, int64, error)
	FindByDestination(ctx context.Context, countryID string, stateOrProvinceID, districtID *uint, zipCode string, orderSubtotal float64) (*models.PriceAndDestination, error)
}

type priceAndDestinationRepository struct {
	db *gorm.DB
}

func NewPriceAndDestinationRepository(db *gorm.DB) PriceAndDestinationRepository {
	return &priceAndDestinationRepository{db: db}
}

func (r *priceAndDestinationRepository) Create(ctx context.Context, entity *models.PriceAndDestination) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *priceAndDestinationRepository) Update(ctx context.Context, entity *models.PriceAndDestination) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *priceAndDestinationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PriceAndDestination{}, id).Error
}

func (r *priceAndDestinationRepository) FindByID(ctx context.Context, id uint) (*models.PriceAndDestination, error) {
	var entity models.PriceAndDestination
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *priceAndDestinationRepository) FindAll(ctx context.Context) ([]models.PriceAndDestination, error) {
	var entities []models.PriceAndDestination
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *priceAndDestinationRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.PriceAndDestination, int64, error) {
	var entities []models.PriceAndDestination
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PriceAndDestination{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *priceAndDestinationRepository) FindByDestination(ctx context.Context, countryID string, stateOrProvinceID, districtID *uint, zipCode string, orderSubtotal float64) (*models.PriceAndDestination, error) {
	var entity models.PriceAndDestination

	query := r.db.WithContext(ctx).Where("country_id = ?", countryID).
		Where("min_order_subtotal <= ?", orderSubtotal).
		Order("min_order_subtotal DESC, shipping_price ASC")

	if stateOrProvinceID != nil {
		query = query.Where("(state_or_province_id IS NULL OR state_or_province_id = ?)", *stateOrProvinceID)
	} else {
		query = query.Where("state_or_province_id IS NULL")
	}

	if districtID != nil {
		query = query.Where("(district_id IS NULL OR district_id = ?)", *districtID)
	} else {
		query = query.Where("district_id IS NULL")
	}

	if zipCode != "" {
		query = query.Where("(zip_code = '' OR zip_code = ?)", zipCode)
	} else {
		query = query.Where("zip_code = ''")
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
