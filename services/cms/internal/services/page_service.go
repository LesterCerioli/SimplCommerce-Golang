package services

import (
	"context"
	"errors"
	"time"

	"github.com/simplcommerce-go/services/cms/internal/models"
	"github.com/simplcommerce-go/services/cms/internal/repositories"
)

var (
	ErrPageNotFound = errors.New("page not found")
	ErrSlugExists   = errors.New("slug already exists")
)

type PageService interface {
	Create(ctx context.Context, page *models.Page) error
	Update(ctx context.Context, page *models.Page) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Page, error)
	FindBySlug(ctx context.Context, slug string) (*models.Page, error)
	FindAll(ctx context.Context) ([]models.Page, error)
	Publish(ctx context.Context, id uint) error
	Unpublish(ctx context.Context, id uint) error
}

type pageService struct {
	pageRepo repositories.PageRepository
}

func NewPageService(pageRepo repositories.PageRepository) PageService {
	return &pageService{pageRepo: pageRepo}
}

func (s *pageService) Create(ctx context.Context, page *models.Page) error {
	existing, _ := s.pageRepo.FindBySlug(ctx, page.Slug)
	if existing != nil {
		return ErrSlugExists
	}
	return s.pageRepo.Create(ctx, page)
}

func (s *pageService) Update(ctx context.Context, page *models.Page) error {
	_, err := s.pageRepo.FindByID(ctx, page.ID)
	if err != nil {
		return ErrPageNotFound
	}
	return s.pageRepo.Update(ctx, page)
}

func (s *pageService) Delete(ctx context.Context, id uint) error {
	_, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPageNotFound
	}
	return s.pageRepo.Delete(ctx, id)
}

func (s *pageService) FindByID(ctx context.Context, id uint) (*models.Page, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPageNotFound
	}
	return page, nil
}

func (s *pageService) FindBySlug(ctx context.Context, slug string) (*models.Page, error) {
	page, err := s.pageRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, ErrPageNotFound
	}
	return page, nil
}

func (s *pageService) FindAll(ctx context.Context) ([]models.Page, error) {
	return s.pageRepo.FindAll(ctx)
}

func (s *pageService) Publish(ctx context.Context, id uint) error {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPageNotFound
	}
	now := time.Now()
	page.IsPublished = true
	page.PublishedOn = &now
	return s.pageRepo.Update(ctx, page)
}

func (s *pageService) Unpublish(ctx context.Context, id uint) error {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPageNotFound
	}
	page.IsPublished = false
	return s.pageRepo.Update(ctx, page)
}
