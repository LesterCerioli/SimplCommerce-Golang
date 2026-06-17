package entities

type Country struct {
	ID                string `gorm:"primaryKey;size:10"`
	Name              string `gorm:"size:450"`
	Code3             string `gorm:"size:3"`
	IsBillingEnabled  bool   `gorm:"default:true"`
	IsShippingEnabled bool   `gorm:"default:true"`
	IsCityEnabled     bool   `gorm:"default:true"`
	IsZipCodeEnabled  bool   `gorm:"default:true"`
	IsDistrictEnabled bool   `gorm:"default:false"`
}

func (Country) TableName() string { return "identity_countries" }
