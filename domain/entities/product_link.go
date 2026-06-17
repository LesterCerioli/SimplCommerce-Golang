package entities

type ProductLink struct {
	BaseEntity
	ProductID       string `gorm:"type:uuid;not null"`
	LinkedProductID string `gorm:"type:uuid;not null"`
	LinkType        int  `gorm:"default:2"`
	Product         Product `gorm:"foreignKey:ProductID"`
	LinkedProduct   Product `gorm:"foreignKey:LinkedProductID"`
}

func (ProductLink) TableName() string { return "catalog_product_links" }
