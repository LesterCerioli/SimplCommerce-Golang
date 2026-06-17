package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/tax/internal/models"
	"github.com/simplcommerce-go/services/tax/internal/repositories"
)

type TaxService interface {
	// TaxClass operations
	CreateTaxClass(ctx context.Context, name string) (*models.TaxClass, error)
	GetTaxClass(ctx context.Context, id uint) (*models.TaxClass, error)
	UpdateTaxClass(ctx context.Context, id uint, name string) (*models.TaxClass, error)
	DeleteTaxClass(ctx context.Context, id uint) error
	GetAllTaxClasses(ctx context.Context, page, pageSize int) ([]models.TaxClass, int64, error)

	// TaxRate operations
	CreateTaxRate(ctx context.Context, input CreateTaxRateInput) (*models.TaxRate, error)
	GetTaxRate(ctx context.Context, id uint) (*models.TaxRate, error)
	UpdateTaxRate(ctx context.Context, id uint, input UpdateTaxRateInput) (*models.TaxRate, error)
	DeleteTaxRate(ctx context.Context, id uint) error
	GetAllTaxRates(ctx context.Context, page, pageSize int) ([]models.TaxRate, int64, error)

	// Tax calculation
	CalculateTax(ctx context.Context, input CalculateTaxInput) (*TaxCalculationResult, error)
}

type CreateTaxRateInput struct {
	TaxClassID        uint   `json:"taxClassId"`
	CountryID         string `json:"countryId"`
	StateOrProvinceID *uint  `json:"stateOrProvinceId"`
	Rate              float64 `json:"rate"`
	ZipCode           string `json:"zipCode"`
}

type UpdateTaxRateInput struct {
	TaxClassID        *uint   `json:"taxClassId"`
	CountryID         *string `json:"countryId"`
	StateOrProvinceID *uint   `json:"stateOrProvinceId"`
	Rate              *float64 `json:"rate"`
	ZipCode           *string `json:"zipCode"`
}

type CalculateTaxInput struct {
	CountryID         string  `json:"countryId"`
	StateOrProvinceID *uint   `json:"stateOrProvinceId"`
	ZipCode           string  `json:"zipCode"`
	OrderAmount       float64 `json:"orderAmount"`
}

type TaxCalculationResult struct {
	TaxAmount float64          `json:"taxAmount"`
	TaxRates  []AppliedTaxRate `json:"taxRates"`
}

type AppliedTaxRate struct {
	TaxClassID   uint    `json:"taxClassId"`
	TaxClassName string  `json:"taxClassName"`
	Rate         float64 `json:"rate"`
	Amount       float64 `json:"amount"`
}

type taxService struct {
	taxClassRepo repositories.TaxClassRepository
	taxRateRepo  repositories.TaxRateRepository
}

func NewTaxService(taxClassRepo repositories.TaxClassRepository, taxRateRepo repositories.TaxRateRepository) TaxService {
	return &taxService{
		taxClassRepo: taxClassRepo,
		taxRateRepo:  taxRateRepo,
	}
}

func (s *taxService) CreateTaxClass(ctx context.Context, name string) (*models.TaxClass, error) {
	if name == "" {
		return nil, errors.New("tax class name is required")
	}

	entity := &models.TaxClass{Name: name}
	if err := s.taxClassRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *taxService) GetTaxClass(ctx context.Context, id uint) (*models.TaxClass, error) {
	entity, err := s.taxClassRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax class not found")
		}
		return nil, err
	}
	return entity, nil
}

func (s *taxService) UpdateTaxClass(ctx context.Context, id uint, name string) (*models.TaxClass, error) {
	if name == "" {
		return nil, errors.New("tax class name is required")
	}

	entity, err := s.taxClassRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax class not found")
		}
		return nil, err
	}

	entity.Name = name
	if err := s.taxClassRepo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *taxService) DeleteTaxClass(ctx context.Context, id uint) error {
	if err := s.taxClassRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tax class not found")
		}
		return err
	}
	return nil
}

func (s *taxService) GetAllTaxClasses(ctx context.Context, page, pageSize int) ([]models.TaxClass, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.taxClassRepo.Paginate(ctx, page, pageSize)
}

func (s *taxService) CreateTaxRate(ctx context.Context, input CreateTaxRateInput) (*models.TaxRate, error) {
	if input.TaxClassID == 0 {
		return nil, errors.New("tax class id is required")
	}
	if input.CountryID == "" {
		return nil, errors.New("country id is required")
	}
	if input.Rate < 0 {
		return nil, errors.New("tax rate must be non-negative")
	}

	_, err := s.taxClassRepo.FindByID(ctx, input.TaxClassID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax class not found")
		}
		return nil, err
	}

	entity := &models.TaxRate{
		TaxClassID:        input.TaxClassID,
		CountryID:         input.CountryID,
		StateOrProvinceID: input.StateOrProvinceID,
		Rate:              input.Rate,
		ZipCode:           input.ZipCode,
	}

	if err := s.taxRateRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *taxService) GetTaxRate(ctx context.Context, id uint) (*models.TaxRate, error) {
	entity, err := s.taxRateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax rate not found")
		}
		return nil, err
	}
	return entity, nil
}

func (s *taxService) UpdateTaxRate(ctx context.Context, id uint, input UpdateTaxRateInput) (*models.TaxRate, error) {
	entity, err := s.taxRateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax rate not found")
		}
		return nil, err
	}

	if input.TaxClassID != nil {
		_, err := s.taxClassRepo.FindByID(ctx, *input.TaxClassID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("tax class not found")
			}
			return nil, err
		}
		entity.TaxClassID = *input.TaxClassID
	}
	if input.CountryID != nil {
		entity.CountryID = *input.CountryID
	}
	if input.StateOrProvinceID != nil {
		entity.StateOrProvinceID = input.StateOrProvinceID
	}
	if input.Rate != nil {
		if *input.Rate < 0 {
			return nil, errors.New("tax rate must be non-negative")
		}
		entity.Rate = *input.Rate
	}
	if input.ZipCode != nil {
		entity.ZipCode = *input.ZipCode
	}

	if err := s.taxRateRepo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *taxService) DeleteTaxRate(ctx context.Context, id uint) error {
	if err := s.taxRateRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tax rate not found")
		}
		return err
	}
	return nil
}

func (s *taxService) GetAllTaxRates(ctx context.Context, page, pageSize int) ([]models.TaxRate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.taxRateRepo.Paginate(ctx, page, pageSize)
}

func (s *taxService) CalculateTax(ctx context.Context, input CalculateTaxInput) (*TaxCalculationResult, error) {
	if input.CountryID == "" {
		return nil, errors.New("country id is required")
	}

	rates, err := s.taxRateRepo.FindByLocation(ctx, input.CountryID, input.StateOrProvinceID, input.ZipCode)
	if err != nil {
		return nil, err
	}

	var totalTax float64
	var appliedRates []AppliedTaxRate

	for _, rate := range rates {
		taxAmount := input.OrderAmount * rate.Rate / 100.0
		totalTax += taxAmount

		className := ""
		if rate.TaxClass.Name != "" {
			className = rate.TaxClass.Name
		}

		appliedRates = append(appliedRates, AppliedTaxRate{
			TaxClassID:   rate.TaxClassID,
			TaxClassName: className,
			Rate:         rate.Rate,
			Amount:       taxAmount,
		})
	}

	return &TaxCalculationResult{
		TaxAmount: totalTax,
		TaxRates:  appliedRates,
	}, nil
}
