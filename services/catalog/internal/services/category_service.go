package services

import (
	"context"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category *catalogmodels.Category) error
	UpdateCategory(ctx context.Context, category *catalogmodels.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryByID(ctx context.Context, id uint) (*catalogmodels.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*catalogmodels.Category, error)
	GetCategoryTree(ctx context.Context) ([]catalogmodels.Category, error)
	GetCategories(ctx context.Context, page, pageSize int) ([]catalogmodels.Category, int64, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, category *catalogmodels.Category) error {
	if category.Slug == "" {
		category.Slug = generateSlug(category.Name)
	}
	return s.categoryRepo.Create(ctx, category)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *catalogmodels.Category) error {
	if category.Slug == "" {
		category.Slug = generateSlug(category.Name)
	}
	return s.categoryRepo.Update(ctx, category)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (*catalogmodels.Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*catalogmodels.Category, error) {
	return s.categoryRepo.FindBySlug(ctx, slug)
}

func (s *categoryService) GetCategoryTree(ctx context.Context) ([]catalogmodels.Category, error) {
	return s.categoryRepo.GetTree(ctx)
}

func (s *categoryService) GetCategories(ctx context.Context, page, pageSize int) ([]catalogmodels.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.categoryRepo.Paginate(ctx, page, pageSize)
}
