package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/shipping/internal/models"
	"github.com/simplcommerce-go/services/shipping/internal/repositories"
)

type ShipmentService interface {
	CreateShipment(ctx context.Context, input CreateShipmentInput, createdByID uint) (*models.Shipment, error)
	GetShipment(ctx context.Context, id uint) (*models.Shipment, error)
	GetAllShipments(ctx context.Context, page, pageSize int) ([]models.Shipment, int64, error)
}

type CreateShipmentInput struct {
	OrderID        uint                 `json:"orderId"`
	TrackingNumber string               `json:"trackingNumber"`
	WarehouseID    uint                 `json:"warehouseId"`
	VendorID       *uint                `json:"vendorId"`
	Items          []ShipmentItemInput  `json:"items"`
}

type ShipmentItemInput struct {
	OrderItemID uint `json:"orderItemId"`
	ProductID   uint `json:"productId"`
	Quantity    int  `json:"quantity"`
}

type shipmentService struct {
	shipmentRepo     repositories.ShipmentRepository
	shipmentItemRepo repositories.ShipmentItemRepository
	db               *gorm.DB
}

func NewShipmentService(shipmentRepo repositories.ShipmentRepository, shipmentItemRepo repositories.ShipmentItemRepository, db *gorm.DB) ShipmentService {
	return &shipmentService{
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		db:               db,
	}
}

func (s *shipmentService) CreateShipment(ctx context.Context, input CreateShipmentInput, createdByID uint) (*models.Shipment, error) {
	if input.OrderID == 0 {
		return nil, errors.New("order id is required")
	}
	if input.WarehouseID == 0 {
		return nil, errors.New("warehouse id is required")
	}
	if len(input.Items) == 0 {
		return nil, errors.New("at least one shipment item is required")
	}

	shipment := &models.Shipment{
		OrderID:        input.OrderID,
		TrackingNumber: input.TrackingNumber,
		WarehouseID:    input.WarehouseID,
		VendorID:       input.VendorID,
		CreatedByID:    createdByID,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(shipment).Error; err != nil {
			return err
		}

		for _, item := range input.Items {
			shipmentItem := models.ShipmentItem{
				ShipmentID:  shipment.ID,
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
			}
			if err := tx.Create(&shipmentItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.shipmentRepo.FindByID(ctx, shipment.ID)
}

func (s *shipmentService) GetShipment(ctx context.Context, id uint) (*models.Shipment, error) {
	shipment, err := s.shipmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shipment not found")
		}
		return nil, err
	}
	return shipment, nil
}

func (s *shipmentService) GetAllShipments(ctx context.Context, page, pageSize int) ([]models.Shipment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.shipmentRepo.Paginate(ctx, page, pageSize)
}

type ShippingProviderService interface {
	GetAllProviders(ctx context.Context) ([]models.ShippingProvider, error)
	GetEnabledProviders(ctx context.Context) ([]models.ShippingProvider, error)
	GetProviderByID(ctx context.Context, id string) (*models.ShippingProvider, error)
	UpdateProvider(ctx context.Context, id string, input UpdateShippingProviderInput) (*models.ShippingProvider, error)
}

type UpdateShippingProviderInput struct {
	Name                                  *string `json:"name"`
	IsEnabled                             *bool   `json:"isEnabled"`
	ConfigureURL                          *string `json:"configureUrl"`
	ToAllShippingEnabledCountries         *bool   `json:"toAllShippingEnabledCountries"`
	OnlyCountryIDsString                  *string `json:"onlyCountryIdsString"`
	ToAllShippingEnabledStatesOrProvinces *bool   `json:"toAllShippingEnabledStatesOrProvinces"`
	OnlyStateOrProvinceIDsString          *string `json:"onlyStateOrProvinceIdsString"`
	AdditionalSettings                    *string `json:"additionalSettings"`
	ShippingPriceServiceTypeName          *string `json:"shippingPriceServiceTypeName"`
}

type shippingProviderService struct {
	providerRepo repositories.ShippingProviderRepository
}

func NewShippingProviderService(providerRepo repositories.ShippingProviderRepository) ShippingProviderService {
	return &shippingProviderService{providerRepo: providerRepo}
}

func (s *shippingProviderService) GetAllProviders(ctx context.Context) ([]models.ShippingProvider, error) {
	return s.providerRepo.FindAll(ctx)
}

func (s *shippingProviderService) GetEnabledProviders(ctx context.Context) ([]models.ShippingProvider, error) {
	return s.providerRepo.FindEnabled(ctx)
}

func (s *shippingProviderService) GetProviderByID(ctx context.Context, id string) (*models.ShippingProvider, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shipping provider not found")
		}
		return nil, err
	}
	return provider, nil
}

func (s *shippingProviderService) UpdateProvider(ctx context.Context, id string, input UpdateShippingProviderInput) (*models.ShippingProvider, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shipping provider not found")
		}
		return nil, err
	}

	if input.Name != nil {
		provider.Name = *input.Name
	}
	if input.IsEnabled != nil {
		provider.IsEnabled = *input.IsEnabled
	}
	if input.ConfigureURL != nil {
		provider.ConfigureURL = *input.ConfigureURL
	}
	if input.ToAllShippingEnabledCountries != nil {
		provider.ToAllShippingEnabledCountries = *input.ToAllShippingEnabledCountries
	}
	if input.OnlyCountryIDsString != nil {
		provider.OnlyCountryIDsString = *input.OnlyCountryIDsString
	}
	if input.ToAllShippingEnabledStatesOrProvinces != nil {
		provider.ToAllShippingEnabledStatesOrProvinces = *input.ToAllShippingEnabledStatesOrProvinces
	}
	if input.OnlyStateOrProvinceIDsString != nil {
		provider.OnlyStateOrProvinceIDsString = *input.OnlyStateOrProvinceIDsString
	}
	if input.AdditionalSettings != nil {
		provider.AdditionalSettings = *input.AdditionalSettings
	}
	if input.ShippingPriceServiceTypeName != nil {
		provider.ShippingPriceServiceTypeName = *input.ShippingPriceServiceTypeName
	}

	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

type ShippingRateService interface {
	CalculateRate(ctx context.Context, input ShippingRateInput) (*ShippingRateResult, error)
}

type ShippingRateInput struct {
	CountryID         string `json:"countryId"`
	StateOrProvinceID *uint  `json:"stateOrProvinceId"`
	DistrictID        *uint  `json:"districtId"`
	ZipCode           string `json:"zipCode"`
	OrderSubtotal     float64 `json:"orderSubtotal"`
}

type ShippingRateResult struct {
	Price      float64 `json:"price"`
	Note       string  `json:"note"`
	ProviderID string  `json:"providerId"`
}

type shippingRateService struct {
	priceRepo repositories.PriceAndDestinationRepository
}

func NewShippingRateService(priceRepo repositories.PriceAndDestinationRepository) ShippingRateService {
	return &shippingRateService{priceRepo: priceRepo}
}

func (s *shippingRateService) CalculateRate(ctx context.Context, input ShippingRateInput) (*ShippingRateResult, error) {
	if input.CountryID == "" {
		return nil, errors.New("country id is required")
	}

	price, err := s.priceRepo.FindByDestination(ctx, input.CountryID, input.StateOrProvinceID, input.DistrictID, input.ZipCode, input.OrderSubtotal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no shipping rate available for the given destination")
		}
		return nil, err
	}

	return &ShippingRateResult{
		Price: price.ShippingPrice,
		Note:  price.Note,
	}, nil
}

type PriceDestinationService interface {
	GetAll(ctx context.Context, page, pageSize int) ([]models.PriceAndDestination, int64, error)
	Create(ctx context.Context, input CreatePriceDestinationInput) (*models.PriceAndDestination, error)
}

type CreatePriceDestinationInput struct {
	CountryID         string  `json:"countryId"`
	StateOrProvinceID *uint   `json:"stateOrProvinceId"`
	DistrictID        *uint   `json:"districtId"`
	ZipCode           string  `json:"zipCode"`
	Note              string  `json:"note"`
	MinOrderSubtotal  float64 `json:"minOrderSubtotal"`
	ShippingPrice     float64 `json:"shippingPrice"`
}

type priceDestinationService struct {
	priceRepo repositories.PriceAndDestinationRepository
}

func NewPriceDestinationService(priceRepo repositories.PriceAndDestinationRepository) PriceDestinationService {
	return &priceDestinationService{priceRepo: priceRepo}
}

func (s *priceDestinationService) GetAll(ctx context.Context, page, pageSize int) ([]models.PriceAndDestination, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.priceRepo.Paginate(ctx, page, pageSize)
}

func (s *priceDestinationService) Create(ctx context.Context, input CreatePriceDestinationInput) (*models.PriceAndDestination, error) {
	if input.CountryID == "" {
		return nil, errors.New("country id is required")
	}
	if input.ShippingPrice < 0 {
		return nil, errors.New("shipping price must be non-negative")
	}

	entity := &models.PriceAndDestination{
		CountryID:         input.CountryID,
		StateOrProvinceID: input.StateOrProvinceID,
		DistrictID:        input.DistrictID,
		ZipCode:           input.ZipCode,
		Note:              input.Note,
		MinOrderSubtotal:  input.MinOrderSubtotal,
		ShippingPrice:     input.ShippingPrice,
	}

	if err := s.priceRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
