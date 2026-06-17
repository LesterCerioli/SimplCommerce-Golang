package entities

type ProductLink struct {
	BaseEntity
	ProductID       uint `gorm:"not null"`
	LinkedProductID uint `gorm:"not null"`
	LinkType        int  `gorm:"default:2"`
	Product         Product `gorm:"foreignKey:ProductID"`
	LinkedProduct   Product `gorm:"foreignKey:LinkedProductID"`
}

func (ProductLink) TableName() string { return "catalog_product_links" }
