package entities

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AuditableEntity struct {
	CreatedByID string  `gorm:"type:uuid;not null" json:"createdById"`
	UpdatedByID *string `gorm:"type:uuid" json:"updatedById,omitempty"`
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
