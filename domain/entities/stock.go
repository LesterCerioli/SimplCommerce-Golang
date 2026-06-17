package entities

type Stock struct {
	BaseEntity
	ProductID        uint `gorm:"not null;uniqueIndex:idx_product_warehouse"`
	WarehouseID      uint `gorm:"not null;uniqueIndex:idx_product_warehouse"`
	Quantity         int  `gorm:"default:0"`
	ReservedQuantity int  `gorm:"default:0"`
	Warehouse        Warehouse `gorm:"foreignKey:WarehouseID"`
}

func (Stock) TableName() string { return "inventory_stocks" }

type StockHistory struct {
	BaseEntity
	ProductID        uint   `gorm:"not null"`
	WarehouseID      uint   `gorm:"not null"`
	CreatedByID      uint
	AdjustedQuantity int
	Note             string `gorm:"size:1000"`
}

func (StockHistory) TableName() string { return "inventory_stock_histories" }

type Warehouse struct {
	BaseEntity
	Name      string `gorm:"size:450;not null"`
	AddressID uint   `gorm:"not null"`
	VendorID  *uint
}

func (Warehouse) TableName() string { return "inventory_warehouses" }

type ProductBackInStockSubscription struct {
	BaseEntity
	ProductID     uint   `gorm:"not null"`
	CustomerEmail string `gorm:"size:256;not null"`
}

func (ProductBackInStockSubscription) TableName() string { return "inventory_product_back_in_stock_subscriptions" }
