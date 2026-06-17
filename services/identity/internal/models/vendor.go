package models

import "github.com/simplcommerce-go/pkg/models"

type Vendor struct {
	models.BaseEntity
	Name        string `gorm:"size:450;not null"`
	Slug        string `gorm:"uniqueIndex;size:450"`
	Description string `gorm:"type:text"`
	Email       string `gorm:"size:256"`
	IsActive    bool   `gorm:"default:true"`
	IsDeleted   bool   `gorm:"default:false"`
	Users       []User `gorm:"foreignKey:VendorID"`
}

func (Vendor) TableName() string { return "identity_vendors" }
