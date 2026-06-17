package models

import "github.com/simplcommerce-go/pkg/models"

type Payment struct {
	models.BaseEntity
	OrderID              uint    `gorm:"not null;index"`
	Amount               float64 `gorm:"not null"`
	PaymentFee           float64 `gorm:"default:0"`
	PaymentMethod        string  `gorm:"size:450;not null"`
	GatewayTransactionID string  `gorm:"size:500"`
	Status               string  `gorm:"size:50;default:'Succeeded'"`
	FailureMessage       string  `gorm:"type:text"`
}

func (Payment) TableName() string { return "payments_payments" }

type PaymentProvider struct {
	ID                       string `gorm:"primaryKey;size:200"`
	Name                     string `gorm:"size:450;not null"`
	IsEnabled                bool   `gorm:"default:false"`
	ConfigureURL             string `gorm:"size:450"`
	LandingViewComponentName string `gorm:"size:450"`
	AdditionalSettings       string `gorm:"type:text"`
}

func (PaymentProvider) TableName() string { return "payments_payment_providers" }
