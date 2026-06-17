package entities

type StateOrProvince struct {
	BaseEntity
	CountryID string `gorm:"size:10;not null"`
	Code      string `gorm:"size:10"`
	Name      string `gorm:"size:450"`
	Type      string `gorm:"size:50"`
	Country   Country `gorm:"foreignKey:CountryID"`
}

func (StateOrProvince) TableName() string { return "identity_state_or_provinces" }
