package entities

type Shipment struct {
	BaseEntity
	OrderID        uint   `gorm:"not null;index"`
	TrackingNumber string `gorm:"size:450"`
	WarehouseID    uint   `gorm:"not null"`
	VendorID       *uint
	CreatedByID    uint
	Items          []ShipmentItem `gorm:"foreignKey:ShipmentID"`
}

func (Shipment) TableName() string { return "shipping_shipments" }

type ShipmentItem struct {
	BaseEntity
	ShipmentID  uint     `gorm:"not null;index"`
	OrderItemID uint     `gorm:"not null"`
	ProductID   uint     `gorm:"not null"`
	Quantity    int      `gorm:"not null"`
	Shipment    Shipment `gorm:"foreignKey:ShipmentID"`
}

func (ShipmentItem) TableName() string { return "shipping_shipment_items" }

type ShippingProvider struct {
	ID                                    string `gorm:"primaryKey;size:200"`
	Name                                  string `gorm:"size:450;not null"`
	IsEnabled                             bool   `gorm:"default:false"`
	ConfigureURL                          string `gorm:"size:450"`
	ToAllShippingEnabledCountries         bool   `gorm:"default:true"`
	OnlyCountryIDsString                  string `gorm:"size:1000"`
	ToAllShippingEnabledStatesOrProvinces bool   `gorm:"default:true"`
	OnlyStateOrProvinceIDsString          string `gorm:"size:1000"`
	AdditionalSettings                    string `gorm:"type:text"`
	ShippingPriceServiceTypeName          string `gorm:"size:450"`
}

func (ShippingProvider) TableName() string { return "shipping_shipping_providers" }

type PriceAndDestination struct {
	BaseEntity
	CountryID         string  `gorm:"size:10"`
	StateOrProvinceID *uint
	DistrictID        *uint
	ZipCode           string  `gorm:"size:20"`
	Note              string
	MinOrderSubtotal  float64 `gorm:"default:0"`
	ShippingPrice     float64 `gorm:"not null"`
}

func (PriceAndDestination) TableName() string { return "shipping_table_rate_price_and_destinations" }
