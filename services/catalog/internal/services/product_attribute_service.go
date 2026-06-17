package services

import (
	"context"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
)

type ProductAttributeService interface {
	CreateAttribute(ctx context.Context, attr *catalogmodels.ProductAttribute) error
	UpdateAttribute(ctx context.Context, attr *catalogmodels.ProductAttribute) error
	DeleteAttribute(ctx context.Context, id uint) error
	GetAttributeByID(ctx context.Context, id uint) (*catalogmodels.ProductAttribute, error)
	GetAttributes(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttribute, int64, error)
}

type productAttributeService struct {
	attrRepo repositories.ProductAttributeRepository
}

func NewProductAttributeService(attrRepo repositories.ProductAttributeRepository) ProductAttributeService {
	return &productAttributeService{attrRepo: attrRepo}
}

func (s *productAttributeService) CreateAttribute(ctx context.Context, attr *catalogmodels.ProductAttribute) error {
	return s.attrRepo.Create(ctx, attr)
}

func (s *productAttributeService) UpdateAttribute(ctx context.Context, attr *catalogmodels.ProductAttribute) error {
	return s.attrRepo.Update(ctx, attr)
}

func (s *productAttributeService) DeleteAttribute(ctx context.Context, id uint) error {
	return s.attrRepo.Delete(ctx, id)
}

func (s *productAttributeService) GetAttributeByID(ctx context.Context, id uint) (*catalogmodels.ProductAttribute, error) {
	return s.attrRepo.FindByID(ctx, id)
}

func (s *productAttributeService) GetAttributes(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttribute, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.attrRepo.Paginate(ctx, page, pageSize)
}

type ProductAttributeGroupService interface {
	CreateGroup(ctx context.Context, group *catalogmodels.ProductAttributeGroup) error
	UpdateGroup(ctx context.Context, group *catalogmodels.ProductAttributeGroup) error
	DeleteGroup(ctx context.Context, id uint) error
	GetGroupByID(ctx context.Context, id uint) (*catalogmodels.ProductAttributeGroup, error)
	GetGroups(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttributeGroup, int64, error)
}

type productAttributeGroupService struct {
	groupRepo repositories.ProductAttributeGroupRepository
}

func NewProductAttributeGroupService(groupRepo repositories.ProductAttributeGroupRepository) ProductAttributeGroupService {
	return &productAttributeGroupService{groupRepo: groupRepo}
}

func (s *productAttributeGroupService) CreateGroup(ctx context.Context, group *catalogmodels.ProductAttributeGroup) error {
	return s.groupRepo.Create(ctx, group)
}

func (s *productAttributeGroupService) UpdateGroup(ctx context.Context, group *catalogmodels.ProductAttributeGroup) error {
	return s.groupRepo.Update(ctx, group)
}

func (s *productAttributeGroupService) DeleteGroup(ctx context.Context, id uint) error {
	return s.groupRepo.Delete(ctx, id)
}

func (s *productAttributeGroupService) GetGroupByID(ctx context.Context, id uint) (*catalogmodels.ProductAttributeGroup, error) {
	return s.groupRepo.FindByID(ctx, id)
}

func (s *productAttributeGroupService) GetGroups(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttributeGroup, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.groupRepo.Paginate(ctx, page, pageSize)
}
