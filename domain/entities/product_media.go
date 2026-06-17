package entities

type ProductMedia struct {
	BaseEntity
	ProductID    uint `gorm:"not null"`
	MediaID      uint `gorm:"not null"`
	DisplayOrder int  `gorm:"default:0"`
	Product      Product `gorm:"foreignKey:ProductID"`
	Media        Media   `gorm:"foreignKey:MediaID"`
}

func (ProductMedia) TableName() string { return "catalog_product_medias" }
