package entities

type ProductCategory struct {
	ProductID  string `gorm:"primaryKey;type:uuid"`
	CategoryID string `gorm:"primaryKey;type:uuid"`
}

func (ProductCategory) TableName() string { return "catalog_product_categories" }
