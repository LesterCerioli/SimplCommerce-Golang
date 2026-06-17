package entities

import (
	"time"
)

type ProductPriceHistory struct {
	BaseEntity
	ProductID         string `gorm:"type:uuid"`
	CreatedByID       string `gorm:"type:uuid"`
	Price             *float64
	OldPrice          *float64
	SpecialPrice      *float64
	SpecialPriceStart *time.Time
	SpecialPriceEnd   *time.Time
	Product           Product `gorm:"foreignKey:ProductID"`
}

func (ProductPriceHistory) TableName() string { return "catalog_product_price_histories" }
