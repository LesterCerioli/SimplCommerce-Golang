package models

import "github.com/simplcommerce-go/pkg/models"

type ProductLink struct {
	models.BaseEntity
	ProductID       uint `gorm:"not null"`
	LinkedProductID uint `gorm:"not null"`
	LinkType        int  `gorm:"default:2"`
	Product         Product `gorm:"foreignKey:ProductID"`
	LinkedProduct   Product `gorm:"foreignKey:LinkedProductID"`
}

func (ProductLink) TableName() string { return "catalog_product_links" }
