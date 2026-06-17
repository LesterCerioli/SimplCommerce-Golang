package entities

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AuditableEntity struct {
	CreatedByID uint  `gorm:"not null" json:"createdById"`
	UpdatedByID *uint `json:"updatedById,omitempty"`
}

type SEOEntity struct {
	MetaTitle       string `gorm:"size:450" json:"metaTitle,omitempty"`
	MetaKeywords    string `gorm:"size:450" json:"metaKeywords,omitempty"`
	MetaDescription string `gorm:"type:text" json:"metaDescription,omitempty"`
}

type PublishableEntity struct {
	IsPublished bool       `gorm:"default:false" json:"isPublished"`
	PublishedOn *time.Time `json:"publishedOn,omitempty"`
}
