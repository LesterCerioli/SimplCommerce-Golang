package models

import "github.com/simplcommerce-go/pkg/models"

type UserAddress struct {
	models.BaseEntity
	UserID      uint `gorm:"not null"`
	AddressID   uint `gorm:"not null"`
	AddressType int  `gorm:"default:0"`
	User        User    `gorm:"foreignKey:UserID"`
	Address     Address `gorm:"foreignKey:AddressID"`
}

func (UserAddress) TableName() string { return "identity_user_addresses" }
