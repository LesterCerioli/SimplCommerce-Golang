package models

import "github.com/simplcommerce-go/pkg/models"

type MenuItem struct {
	models.BaseEntity
	ParentID     *uint
	MenuID       uint   `gorm:"not null"`
	EntityID     *uint  `gorm:"index"`
	CustomLink   string `gorm:"size:450"`
	Name         string `gorm:"size:450;not null"`
	DisplayOrder int    `gorm:"default:0"`
	Parent       *MenuItem `gorm:"foreignKey:ParentID"`
	Menu         Menu      `gorm:"foreignKey:MenuID"`
}

func (MenuItem) TableName() string { return "cms_menu_items" }
