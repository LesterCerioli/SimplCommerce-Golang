package entities

import (
	"time"
)

type Product struct {
	BaseEntity
	Name                   string                      `gorm:"size:450;not null"`
	Slug                   string                      `gorm:"uniqueIndex;size:450;not null"`
	ShortDescription       string                      `gorm:"size:450"`
	Description            string                      `gorm:"type:text"`
	Specification          string                      `gorm:"type:text"`
	Price                  float64                     `gorm:"not null;default:0"`
	OldPrice               *float64
	SpecialPrice           *float64
	SpecialPriceStart      *time.Time
	SpecialPriceEnd        *time.Time
	HasOptions             bool                        `gorm:"default:false"`
	IsVisibleIndividually  bool                        `gorm:"default:true"`
	IsFeatured             bool                        `gorm:"default:false"`
	IsCallForPricing       bool                        `gorm:"default:false"`
	IsAllowToOrder         bool                        `gorm:"default:true"`
	StockTrackingIsEnabled bool                        `gorm:"default:true"`
	StockQuantity          int                         `gorm:"default:0"`
	SKU                    string                      `gorm:"size:200"`
	GTIN                   string                      `gorm:"size:50"`
	NormalizedName         string                      `gorm:"size:450"`
	DisplayOrder           int                         `gorm:"default:0"`
	ReviewsCount           int                         `gorm:"default:0"`
	RatingAverage          *float64
	VendorID               *string `gorm:"type:uuid"`
	BrandID                *string `gorm:"type:uuid"`
	TaxClassID             *string `gorm:"type:uuid"`
	ThumbnailImageID       *string `gorm:"type:uuid"`
	IsPublished            bool                        `gorm:"default:false"`
	PublishedOn            *time.Time

	Categories         []Category                   `gorm:"many2many:catalog_product_categories;"`
	Brand              *Brand                       `gorm:"foreignKey:BrandID"`
	AttributeValues    []ProductAttributeValue      `gorm:"foreignKey:ProductID"`
	OptionValues       []ProductOptionValue         `gorm:"foreignKey:ProductID"`
	OptionCombinations []ProductOptionCombination   `gorm:"foreignKey:ProductID"`
	LinkedProducts     []ProductLink                `gorm:"foreignKey:ProductID"`
	Medias             []ProductMedia               `gorm:"foreignKey:ProductID"`
	PriceHistories     []ProductPriceHistory        `gorm:"foreignKey:ProductID"`
	ThumbnailImage     *Media                       `gorm:"foreignKey:ThumbnailImageID"`
}

func (Product) TableName() string { return "catalog_products" }
