package entities

type OrderAddress struct {
	BaseEntity
	ContactName       string `gorm:"size:450;not null"`
	Phone             string `gorm:"size:50"`
	AddressLine1      string `gorm:"size:450;not null"`
	AddressLine2      string `gorm:"size:450"`
	City              string `gorm:"size:200;not null"`
	ZipCode           string `gorm:"size:20"`
	DistrictID        *uint
	StateOrProvinceID uint   `gorm:"not null"`
	CountryID         string `gorm:"size:10;not null"`
}

func (OrderAddress) TableName() string { return "orders_order_addresses" }
