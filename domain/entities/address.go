package entities

type Address struct {
	BaseEntity
	ContactName       string `gorm:"size:450"`
	Phone             string `gorm:"size:50"`
	AddressLine1      string `gorm:"size:450"`
	AddressLine2      string `gorm:"size:450"`
	City              string `gorm:"size:200"`
	ZipCode           string `gorm:"size:20"`
	DistrictID        *string `gorm:"type:uuid"`
	StateOrProvinceID string  `gorm:"type:uuid;not null"`
	CountryID         string  `gorm:"size:10"`

	StateOrProvince StateOrProvince `gorm:"foreignKey:StateOrProvinceID"`
	Country         Country         `gorm:"foreignKey:CountryID"`
	District        *District       `gorm:"foreignKey:DistrictID"`
}

func (Address) TableName() string { return "identity_addresses" }
