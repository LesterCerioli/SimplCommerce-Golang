package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/cart/internal/models"
	"github.com/simplcommerce-go/services/cart/internal/repositories"
	"gorm.io/gorm"
)

type CartService interface {
	AddItem(ctx context.Context, customerID, productID uint, quantity int) error
	UpdateQuantity(ctx context.Context, customerID, itemID uint, quantity int) error
	RemoveItem(ctx context.Context, customerID, itemID uint) error
	GetCart(ctx context.Context, customerID uint) ([]models.CartItem, error)
	ClearCart(ctx context.Context, customerID uint) error
}

type cartService struct {
	repo repositories.CartItemRepository
}

func NewCartService(repo repositories.CartItemRepository) CartService {
	return &cartService{repo: repo}
}

func (s *cartService) AddItem(ctx context.Context, customerID, productID uint, quantity int) error {
	existing, err := s.repo.FindByCustomerAndProduct(ctx, customerID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item := &models.CartItem{
				CustomerID: customerID,
				ProductID:  productID,
				Quantity:   quantity,
			}
			return s.repo.Create(ctx, item)
		}
		return err
	}
	existing.Quantity += quantity
	return s.repo.Update(ctx, existing)
}

func (s *cartService) UpdateQuantity(ctx context.Context, customerID, itemID uint, quantity int) error {
	item, err := s.repo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.CustomerID != customerID {
		return errors.New("cart item does not belong to customer")
	}
	if quantity <= 0 {
		return s.repo.Delete(ctx, itemID)
	}
	item.Quantity = quantity
	return s.repo.Update(ctx, item)
}

func (s *cartService) RemoveItem(ctx context.Context, customerID, itemID uint) error {
	item, err := s.repo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.CustomerID != customerID {
		return errors.New("cart item does not belong to customer")
	}
	return s.repo.Delete(ctx, itemID)
}

func (s *cartService) GetCart(ctx context.Context, customerID uint) ([]models.CartItem, error) {
	return s.repo.FindByCustomer(ctx, customerID)
}

func (s *cartService) ClearCart(ctx context.Context, customerID uint) error {
	return s.repo.ClearCustomerCart(ctx, customerID)
}
