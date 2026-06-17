package entities

type OrderHistory struct {
	BaseEntity
	OrderID     uint   `gorm:"not null;index"`
	OldStatus   string `gorm:"size:50"`
	NewStatus   string `gorm:"size:50;not null"`
	Note        string `gorm:"size:1000"`
	CreatedByID uint
	Order       Order  `gorm:"foreignKey:OrderID"`
}

func (OrderHistory) TableName() string { return "orders_order_histories" }
