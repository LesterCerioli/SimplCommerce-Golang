package entities

type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}

func (ProductCategory) TableName() string { return "catalog_product_categories" }
