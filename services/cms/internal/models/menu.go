package models

import "github.com/simplcommerce-go/pkg/models"

type Menu struct {
	models.BaseEntity
	Name        string     `gorm:"size:450;not null"`
	IsPublished bool       `gorm:"default:false"`
	IsSystem    bool       `gorm:"default:false"`
	Items       []MenuItem `gorm:"foreignKey:MenuID"`
}

func (Menu) TableName() string { return "cms_menus" }
