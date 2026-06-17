package entities

type UserAddress struct {
	BaseEntity
	UserID      string `gorm:"type:uuid;not null"`
	AddressID   string `gorm:"type:uuid;not null"`
	AddressType int  `gorm:"default:0"`
	User        User    `gorm:"foreignKey:UserID"`
	Address     Address `gorm:"foreignKey:AddressID"`
}

func (UserAddress) TableName() string { return "identity_user_addresses" }
