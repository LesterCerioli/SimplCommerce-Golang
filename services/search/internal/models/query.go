package models

import "github.com/simplcommerce-go/pkg/models"

type Query struct {
	models.BaseEntity
	QueryText    string `gorm:"size:500;not null"`
	ResultsCount int    `gorm:"default:0"`
}

func (Query) TableName() string { return "search_queries" }
