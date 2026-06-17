package entities

type Stock struct {
	BaseEntity
	ProductID        string `gorm:"type:uuid;not null;uniqueIndex:idx_product_warehouse"`
	WarehouseID      string `gorm:"type:uuid;not null;uniqueIndex:idx_product_warehouse"`
	Quantity         int  `gorm:"default:0"`
	ReservedQuantity int  `gorm:"default:0"`
	Warehouse        Warehouse `gorm:"foreignKey:WarehouseID"`
}

func (Stock) TableName() string { return "inventory_stocks" }

type StockHistory struct {
	BaseEntity
	ProductID        string `gorm:"type:uuid;not null"`
	WarehouseID      string `gorm:"type:uuid;not null"`
	CreatedByID      string `gorm:"type:uuid"`
	AdjustedQuantity int
	Note             string `gorm:"size:1000"`
}

func (StockHistory) TableName() string { return "inventory_stock_histories" }

type Warehouse struct {
	BaseEntity
	Name      string `gorm:"size:450;not null"`
	AddressID string  `gorm:"type:uuid;not null"`
	VendorID  *string `gorm:"type:uuid"`
}

func (Warehouse) TableName() string { return "inventory_warehouses" }

type ProductBackInStockSubscription struct {
	BaseEntity
	ProductID     string `gorm:"type:uuid;not null"`
	CustomerEmail string `gorm:"size:256;not null"`
}

func (ProductBackInStockSubscription) TableName() string { return "inventory_product_back_in_stock_subscriptions" }
