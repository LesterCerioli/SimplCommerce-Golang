package models

import "github.com/simplcommerce-go/pkg/models"

type Category struct {
	models.BaseEntity
	Name             string     `gorm:"size:450;not null"`
	Slug             string     `gorm:"uniqueIndex;size:450;not null"`
	Description      string     `gorm:"type:text"`
	DisplayOrder     int        `gorm:"default:0"`
	IsPublished      bool       `gorm:"default:false"`
	IncludeInMenu    bool       `gorm:"default:true"`
	IsDeleted        bool       `gorm:"default:false"`
	ParentID         *uint
	ThumbnailImageID *uint
	Parent           *Category `gorm:"foreignKey:ParentID"`
	Children         []Category `gorm:"foreignKey:ParentID"`
	Products         []Product `gorm:"many2many:catalog_product_categories;"`
}

func (Category) TableName() string { return "catalog_categories" }
