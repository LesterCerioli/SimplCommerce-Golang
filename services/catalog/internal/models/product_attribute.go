package models

import "github.com/simplcommerce-go/pkg/models"

type ProductAttribute struct {
	models.BaseEntity
	Name    string `gorm:"size:450;not null"`
	GroupID uint
	Group   ProductAttributeGroup `gorm:"foreignKey:GroupID"`
}

func (ProductAttribute) TableName() string { return "catalog_product_attributes" }

type ProductAttributeGroup struct {
	models.BaseEntity
	Name       string             `gorm:"size:450;not null"`
	Attributes []ProductAttribute `gorm:"foreignKey:GroupID"`
}

func (ProductAttributeGroup) TableName() string { return "catalog_product_attribute_groups" }

type ProductAttributeValue struct {
	models.BaseEntity
	AttributeID uint   `gorm:"not null"`
	ProductID   uint   `gorm:"not null"`
	Value       string `gorm:"type:text"`
	Attribute   ProductAttribute `gorm:"foreignKey:AttributeID"`
	Product     Product          `gorm:"foreignKey:ProductID"`
}

func (ProductAttributeValue) TableName() string { return "catalog_product_attribute_values" }
