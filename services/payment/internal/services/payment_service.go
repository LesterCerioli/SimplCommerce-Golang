package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/payment/internal/models"
	"github.com/simplcommerce-go/services/payment/internal/repositories"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, id uint) (*models.Payment, error)
	GetAllPayments(ctx context.Context, page, pageSize int) ([]models.Payment, int64, error)
	GetPaymentsByOrderID(ctx context.Context, orderID uint) ([]models.Payment, error)
}

type paymentService struct {
	paymentRepo repositories.PaymentRepository
}

func NewPaymentService(paymentRepo repositories.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) ProcessPayment(ctx context.Context, payment *models.Payment) error {
	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}
	if payment.PaymentMethod == "" {
		return errors.New("payment method is required")
	}
	return s.paymentRepo.Create(ctx, payment)
}

func (s *paymentService) GetPayment(ctx context.Context, id uint) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) GetAllPayments(ctx context.Context, page, pageSize int) ([]models.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.paymentRepo.Paginate(ctx, page, pageSize)
}

func (s *paymentService) GetPaymentsByOrderID(ctx context.Context, orderID uint) ([]models.Payment, error) {
	return s.paymentRepo.FindByOrderID(ctx, orderID)
}

type PaymentProviderService interface {
	GetAllProviders(ctx context.Context) ([]models.PaymentProvider, error)
	GetEnabledProviders(ctx context.Context) ([]models.PaymentProvider, error)
	GetProviderByID(ctx context.Context, id string) (*models.PaymentProvider, error)
	UpdateProvider(ctx context.Context, id string, input UpdateProviderInput) (*models.PaymentProvider, error)
}

type UpdateProviderInput struct {
	Name                     *string `json:"name"`
	IsEnabled                *bool   `json:"isEnabled"`
	ConfigureURL             *string `json:"configureUrl"`
	LandingViewComponentName *string `json:"landingViewComponentName"`
	AdditionalSettings       *string `json:"additionalSettings"`
}

type paymentProviderService struct {
	providerRepo repositories.PaymentProviderRepository
}

func NewPaymentProviderService(providerRepo repositories.PaymentProviderRepository) PaymentProviderService {
	return &paymentProviderService{providerRepo: providerRepo}
}

func (s *paymentProviderService) GetAllProviders(ctx context.Context) ([]models.PaymentProvider, error) {
	return s.providerRepo.FindAll(ctx)
}

func (s *paymentProviderService) GetEnabledProviders(ctx context.Context) ([]models.PaymentProvider, error) {
	return s.providerRepo.FindEnabled(ctx)
}

func (s *paymentProviderService) GetProviderByID(ctx context.Context, id string) (*models.PaymentProvider, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment provider not found")
		}
		return nil, err
	}
	return provider, nil
}

func (s *paymentProviderService) UpdateProvider(ctx context.Context, id string, input UpdateProviderInput) (*models.PaymentProvider, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment provider not found")
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
	if input.LandingViewComponentName != nil {
		provider.LandingViewComponentName = *input.LandingViewComponentName
	}
	if input.AdditionalSettings != nil {
		provider.AdditionalSettings = *input.AdditionalSettings
	}

	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}
