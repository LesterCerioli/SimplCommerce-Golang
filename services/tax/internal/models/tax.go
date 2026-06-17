package models

import "github.com/simplcommerce-go/pkg/models"

type TaxClass struct {
	models.BaseEntity
	Name  string    `gorm:"size:450;not null;unique"`
	Rates []TaxRate `gorm:"foreignKey:TaxClassID"`
}

func (TaxClass) TableName() string { return "tax_tax_classes" }

type TaxRate struct {
	models.BaseEntity
	TaxClassID        uint     `gorm:"not null"`
	CountryID         string   `gorm:"size:10;not null"`
	StateOrProvinceID *uint
	Rate              float64  `gorm:"not null"`
	ZipCode           string   `gorm:"size:20"`
	TaxClass          TaxClass `gorm:"foreignKey:TaxClassID"`
}

func (TaxRate) TableName() string { return "tax_tax_rates" }
