package models

import "github.com/simplcommerce-go/pkg/models"

type CartItem struct {
	models.BaseEntity
	ProductID  uint   `gorm:"not null"`
	CustomerID uint   `gorm:"not null;index"`
	Quantity   int    `gorm:"not null;default:1"`
	VendorID   *uint
}

func (CartItem) TableName() string { return "cart_cart_items" }
