package entities

import (
	"time"
)

type Page struct {
	BaseEntity
	Name        string     `gorm:"size:450;not null"`
	Slug        string     `gorm:"uniqueIndex;size:450;not null"`
	Body        string     `gorm:"type:text"`
	IsPublished bool       `gorm:"default:false"`
	PublishedOn *time.Time
}

func (Page) TableName() string { return "cms_pages" }
