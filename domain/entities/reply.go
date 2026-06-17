package entities

type Reply struct {
	BaseEntity
	ReviewID    string `gorm:"type:uuid;not null"`
	UserID      string `gorm:"type:uuid;not null"`
	Comment     string `gorm:"type:text"`
	ReplierName string `gorm:"size:450"`
	Status      string `gorm:"size:50;default:'Pending'"`
	Review      Review `gorm:"foreignKey:ReviewID"`
}

func (Reply) TableName() string { return "reviews_replies" }
