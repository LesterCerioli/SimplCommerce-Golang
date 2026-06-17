package entities

type CartItem struct {
	BaseEntity
	ProductID  string  `gorm:"type:uuid;not null"`
	CustomerID string  `gorm:"type:uuid;not null;index"`
	Quantity   int     `gorm:"not null;default:1"`
	VendorID   *string `gorm:"type:uuid"`
}

func (CartItem) TableName() string { return "cart_cart_items" }
