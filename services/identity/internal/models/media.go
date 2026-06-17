package models

import "github.com/simplcommerce-go/pkg/models"

type Media struct {
	models.BaseEntity
	Caption   string `gorm:"size:450"`
	FileSize  int64
	FileName  string `gorm:"size:450"`
	MediaType int    `gorm:"default:0"`
}

func (Media) TableName() string { return "identity_media" }
