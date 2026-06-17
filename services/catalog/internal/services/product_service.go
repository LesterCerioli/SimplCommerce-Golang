package services

import (
	"context"
	"regexp"
	"strings"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *catalogmodels.Product) error
	UpdateProduct(ctx context.Context, product *catalogmodels.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	GetProductByID(ctx context.Context, id uint) (*catalogmodels.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*catalogmodels.Product, error)
	SearchProducts(ctx context.Context, query string, minPrice, maxPrice *float64, categoryID *uint, page, pageSize int) ([]catalogmodels.Product, int64, error)
	GetFeaturedProducts(ctx context.Context, count int) ([]catalogmodels.Product, error)
	GetProducts(ctx context.Context, page, pageSize int) ([]catalogmodels.Product, int64, error)
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, product *catalogmodels.Product) error {
	if product.Slug == "" {
		product.Slug = generateSlug(product.Name)
	}
	product.NormalizedName = strings.ToUpper(product.Name)
	return s.productRepo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product *catalogmodels.Product) error {
	if product.Slug == "" {
		product.Slug = generateSlug(product.Name)
	}
	product.NormalizedName = strings.ToUpper(product.Name)
	return s.productRepo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) GetProductByID(ctx context.Context, id uint) (*catalogmodels.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *productService) GetProductBySlug(ctx context.Context, slug string) (*catalogmodels.Product, error) {
	return s.productRepo.FindBySlug(ctx, slug)
}

func (s *productService) SearchProducts(ctx context.Context, query string, minPrice, maxPrice *float64, categoryID *uint, page, pageSize int) ([]catalogmodels.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.productRepo.Search(ctx, query, minPrice, maxPrice, categoryID, page, pageSize)
}

func (s *productService) GetFeaturedProducts(ctx context.Context, count int) ([]catalogmodels.Product, error) {
	if count < 1 {
		count = 10
	}
	return s.productRepo.GetFeatured(ctx, count)
}

func (s *productService) GetProducts(ctx context.Context, page, pageSize int) ([]catalogmodels.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.productRepo.Paginate(ctx, page, pageSize)
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")
	re := regexp.MustCompile(`\-+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
