package models

import "github.com/simplcommerce-go/pkg/models"

type Review struct {
	models.BaseEntity
	UserID       uint   `gorm:"not null;index"`
	Title        string `gorm:"size:450"`
	Comment      string `gorm:"type:text"`
	Rating       int    `gorm:"not null"`
	ReviewerName string `gorm:"size:450"`
	Status       string `gorm:"size:50;default:'Pending'"`
	EntityTypeID string `gorm:"size:450;not null"`
	EntityID     uint   `gorm:"not null"`
	Replies      []Reply `gorm:"foreignKey:ReviewID"`
}

func (Review) TableName() string { return "reviews_reviews" }
