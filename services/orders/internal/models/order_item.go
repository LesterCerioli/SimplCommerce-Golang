package models

import "github.com/simplcommerce-go/pkg/models"

type OrderItem struct {
	models.BaseEntity
	OrderID        uint    `gorm:"not null;index"`
	ProductID      uint    `gorm:"not null"`
	ProductName    string  `gorm:"size:450"`
	ProductSKU     string  `gorm:"size:200"`
	ProductPrice   float64 `gorm:"not null"`
	Quantity       int     `gorm:"not null"`
	DiscountAmount float64 `gorm:"default:0"`
	TaxAmount      float64 `gorm:"default:0"`
	TaxPercent     float64 `gorm:"default:0"`
	Order          Order   `gorm:"foreignKey:OrderID"`
}

func (OrderItem) TableName() string { return "orders_order_items" }
