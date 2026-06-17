package entities

type District struct {
	BaseEntity
	StateOrProvinceID string          `gorm:"type:uuid;not null"`
	Name              string          `gorm:"size:450"`
	Type              string          `gorm:"size:50"`
	Location          string          `gorm:"type:text"`
	StateOrProvince   StateOrProvince `gorm:"foreignKey:StateOrProvinceID"`
}

func (District) TableName() string { return "identity_districts" }
