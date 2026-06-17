package services

import (
	"context"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
)

type ProductOptionService interface {
	CreateOption(ctx context.Context, option *catalogmodels.ProductOption) error
	UpdateOption(ctx context.Context, option *catalogmodels.ProductOption) error
	DeleteOption(ctx context.Context, id uint) error
	GetOptionByID(ctx context.Context, id uint) (*catalogmodels.ProductOption, error)
	GetOptions(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductOption, int64, error)
}

type productOptionService struct {
	optionRepo repositories.ProductOptionRepository
}

func NewProductOptionService(optionRepo repositories.ProductOptionRepository) ProductOptionService {
	return &productOptionService{optionRepo: optionRepo}
}

func (s *productOptionService) CreateOption(ctx context.Context, option *catalogmodels.ProductOption) error {
	return s.optionRepo.Create(ctx, option)
}

func (s *productOptionService) UpdateOption(ctx context.Context, option *catalogmodels.ProductOption) error {
	return s.optionRepo.Update(ctx, option)
}

func (s *productOptionService) DeleteOption(ctx context.Context, id uint) error {
	return s.optionRepo.Delete(ctx, id)
}

func (s *productOptionService) GetOptionByID(ctx context.Context, id uint) (*catalogmodels.ProductOption, error) {
	return s.optionRepo.FindByID(ctx, id)
}

func (s *productOptionService) GetOptions(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductOption, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.optionRepo.Paginate(ctx, page, pageSize)
}

type ProductTemplateService interface {
	CreateTemplate(ctx context.Context, template *catalogmodels.ProductTemplate) error
	UpdateTemplate(ctx context.Context, template *catalogmodels.ProductTemplate) error
	DeleteTemplate(ctx context.Context, id uint) error
	GetTemplateByID(ctx context.Context, id uint) (*catalogmodels.ProductTemplate, error)
	GetTemplates(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductTemplate, int64, error)
}

type productTemplateService struct {
	templateRepo repositories.ProductTemplateRepository
}

func NewProductTemplateService(templateRepo repositories.ProductTemplateRepository) ProductTemplateService {
	return &productTemplateService{templateRepo: templateRepo}
}

func (s *productTemplateService) CreateTemplate(ctx context.Context, template *catalogmodels.ProductTemplate) error {
	return s.templateRepo.Create(ctx, template)
}

func (s *productTemplateService) UpdateTemplate(ctx context.Context, template *catalogmodels.ProductTemplate) error {
	return s.templateRepo.Update(ctx, template)
}

func (s *productTemplateService) DeleteTemplate(ctx context.Context, id uint) error {
	return s.templateRepo.Delete(ctx, id)
}

func (s *productTemplateService) GetTemplateByID(ctx context.Context, id uint) (*catalogmodels.ProductTemplate, error) {
	return s.templateRepo.FindByID(ctx, id)
}

func (s *productTemplateService) GetTemplates(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductTemplate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.templateRepo.Paginate(ctx, page, pageSize)
}
