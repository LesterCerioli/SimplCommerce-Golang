package entities

import (
	"time"
)

type CartRule struct {
	BaseEntity
	Name                 string  `gorm:"size:450;not null"`
	Description          string  `gorm:"type:text"`
	IsActive             bool    `gorm:"default:true"`
	StartOn              *time.Time
	EndOn                *time.Time
	IsCouponRequired     bool    `gorm:"default:false"`
	RuleToApply          string  `gorm:"size:50;default:'by_fixed'"`
	DiscountAmount       float64 `gorm:"not null"`
	MaxDiscountAmount    *float64
	DiscountStep         *int
	UsageLimitPerCoupon  *int
	UsageLimitPerCustomer *int

	Coupons        []Coupon                `gorm:"foreignKey:CartRuleID"`
	Categories     []CartRuleCategory      `gorm:"foreignKey:CartRuleID"`
	Products       []CartRuleProduct       `gorm:"foreignKey:CartRuleID"`
	CustomerGroups []CartRuleCustomerGroup `gorm:"foreignKey:CartRuleID"`
	Usages         []CartRuleUsage         `gorm:"foreignKey:CartRuleID"`
}

func (CartRule) TableName() string { return "pricing_cart_rules" }

type Coupon struct {
	BaseEntity
	CartRuleID string   `gorm:"type:uuid;not null"`
	Code       string   `gorm:"uniqueIndex;size:100;not null"`
	CartRule   CartRule `gorm:"foreignKey:CartRuleID"`
}

func (Coupon) TableName() string { return "pricing_coupons" }

type CartRuleCategory struct {
	CartRuleID string `gorm:"primaryKey;type:uuid"`
	CategoryID string `gorm:"primaryKey;type:uuid"`
}

func (CartRuleCategory) TableName() string { return "pricing_cart_rule_categories" }

type CartRuleProduct struct {
	CartRuleID string `gorm:"primaryKey;type:uuid"`
	ProductID  string `gorm:"primaryKey;type:uuid"`
}

func (CartRuleProduct) TableName() string { return "pricing_cart_rule_products" }

type CartRuleCustomerGroup struct {
	CartRuleID      string `gorm:"primaryKey;type:uuid"`
	CustomerGroupID string `gorm:"primaryKey;type:uuid"`
}

func (CartRuleCustomerGroup) TableName() string { return "pricing_cart_rule_customer_groups" }

type CartRuleUsage struct {
	BaseEntity
	CartRuleID string  `gorm:"type:uuid;not null"`
	CouponID   *string `gorm:"type:uuid"`
	UserID     string  `gorm:"type:uuid;not null"`
	OrderID    string  `gorm:"type:uuid;not null"`
	CartRule   CartRule `gorm:"foreignKey:CartRuleID"`
	Coupon     *Coupon  `gorm:"foreignKey:CouponID"`
}

func (CartRuleUsage) TableName() string { return "pricing_cart_rule_usages" }

type CatalogRule struct {
	BaseEntity
	Name              string  `gorm:"size:450;not null"`
	Description       string  `gorm:"type:text"`
	IsActive          bool    `gorm:"default:true"`
	StartOn           *time.Time
	EndOn             *time.Time
	RuleToApply       string  `gorm:"size:50;default:'by_fixed'"`
	DiscountAmount    float64 `gorm:"not null"`
	MaxDiscountAmount *float64

	CustomerGroups []CatalogRuleCustomerGroup `gorm:"foreignKey:CatalogRuleID"`
}

func (CatalogRule) TableName() string { return "pricing_catalog_rules" }

type CatalogRuleCustomerGroup struct {
	CatalogRuleID   string `gorm:"primaryKey;type:uuid"`
	CustomerGroupID string `gorm:"primaryKey;type:uuid"`
}

func (CatalogRuleCustomerGroup) TableName() string { return "pricing_catalog_rule_customer_groups" }
