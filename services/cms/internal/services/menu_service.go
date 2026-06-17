package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/cms/internal/models"
	"github.com/simplcommerce-go/services/cms/internal/repositories"
)

var (
	ErrMenuNotFound     = errors.New("menu not found")
	ErrMenuItemNotFound = errors.New("menu item not found")
)

type MenuService interface {
	Create(ctx context.Context, menu *models.Menu) error
	Update(ctx context.Context, menu *models.Menu) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Menu, error)
	FindAll(ctx context.Context) ([]models.Menu, error)
	FindItemByID(ctx context.Context, id uint) (*models.MenuItem, error)
	AddItem(ctx context.Context, item *models.MenuItem) error
	UpdateItem(ctx context.Context, item *models.MenuItem) error
	DeleteItem(ctx context.Context, id uint) error
}

type menuService struct {
	menuRepo     repositories.MenuRepository
	menuItemRepo repositories.MenuItemRepository
}

func NewMenuService(menuRepo repositories.MenuRepository, menuItemRepo repositories.MenuItemRepository) MenuService {
	return &menuService{menuRepo: menuRepo, menuItemRepo: menuItemRepo}
}

func (s *menuService) Create(ctx context.Context, menu *models.Menu) error {
	return s.menuRepo.Create(ctx, menu)
}

func (s *menuService) Update(ctx context.Context, menu *models.Menu) error {
	_, err := s.menuRepo.FindByID(ctx, menu.ID)
	if err != nil {
		return ErrMenuNotFound
	}
	return s.menuRepo.Update(ctx, menu)
}

func (s *menuService) Delete(ctx context.Context, id uint) error {
	_, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return ErrMenuNotFound
	}
	return s.menuRepo.Delete(ctx, id)
}

func (s *menuService) FindByID(ctx context.Context, id uint) (*models.Menu, error) {
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrMenuNotFound
	}
	return menu, nil
}

func (s *menuService) FindAll(ctx context.Context) ([]models.Menu, error) {
	return s.menuRepo.FindAll(ctx)
}

func (s *menuService) FindItemByID(ctx context.Context, id uint) (*models.MenuItem, error) {
	item, err := s.menuItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrMenuItemNotFound
	}
	return item, nil
}

func (s *menuService) AddItem(ctx context.Context, item *models.MenuItem) error {
	_, err := s.menuRepo.FindByID(ctx, item.MenuID)
	if err != nil {
		return ErrMenuNotFound
	}
	return s.menuItemRepo.Create(ctx, item)
}

func (s *menuService) UpdateItem(ctx context.Context, item *models.MenuItem) error {
	_, err := s.menuItemRepo.FindByID(ctx, item.ID)
	if err != nil {
		return ErrMenuItemNotFound
	}
	return s.menuItemRepo.Update(ctx, item)
}

func (s *menuService) DeleteItem(ctx context.Context, id uint) error {
	_, err := s.menuItemRepo.FindByID(ctx, id)
	if err != nil {
		return ErrMenuItemNotFound
	}
	return s.menuItemRepo.Delete(ctx, id)
}
