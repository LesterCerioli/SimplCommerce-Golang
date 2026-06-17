package entities

type CustomerGroup struct {
	BaseEntity
	Name        string `gorm:"size:450;uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	IsDeleted   bool   `gorm:"default:false"`
	Users       []User `gorm:"many2many:identity_customer_group_users;"`
}

func (CustomerGroup) TableName() string { return "identity_customer_groups" }
