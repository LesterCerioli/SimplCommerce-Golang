package entities

type Order struct {
	BaseEntity
	CustomerID           string  `gorm:"type:uuid;not null;index"`
	VendorID             *string `gorm:"type:uuid"`
	CreatedByID          string  `gorm:"type:uuid"`
	UpdatedByID          *string `gorm:"type:uuid"`
	CouponCode           string  `gorm:"size:100"`
	CouponRuleName       string  `gorm:"size:450"`
	DiscountAmount       float64 `gorm:"default:0"`
	SubTotal             float64 `gorm:"not null"`
	SubTotalWithDiscount float64 `gorm:"not null"`
	ShippingAddressID    string  `gorm:"type:uuid;not null"`
	BillingAddressID     string  `gorm:"type:uuid;not null"`
	OrderStatus          string  `gorm:"size:50;default:'New'"`
	OrderNote            string  `gorm:"size:1000"`
	ParentID             *string `gorm:"type:uuid"`
	IsMasterOrder        bool    `gorm:"default:false"`
	ShippingMethod       string  `gorm:"size:450"`
	ShippingFeeAmount    float64 `gorm:"default:0"`
	TaxAmount            float64 `gorm:"default:0"`
	OrderTotal           float64 `gorm:"not null"`
	PaymentMethod        string  `gorm:"size:450"`
	PaymentFeeAmount     float64 `gorm:"default:0"`

	Items     []OrderItem    `gorm:"foreignKey:OrderID"`
	Histories []OrderHistory `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string { return "orders_orders" }
