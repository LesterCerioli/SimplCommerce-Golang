package entities

type Brand struct {
	BaseEntity
	Name        string    `gorm:"size:450;not null"`
	Slug        string    `gorm:"uniqueIndex;size:450;not null"`
	Description string    `gorm:"type:text"`
	IsPublished bool      `gorm:"default:false"`
	IsDeleted   bool      `gorm:"default:false"`
	Products    []Product `gorm:"foreignKey:BrandID"`
}

func (Brand) TableName() string { return "catalog_brands" }
