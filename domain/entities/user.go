package entities

type User struct {
	BaseEntity
	UserGuid                 string          `gorm:"uniqueIndex;size:36"`
	FullName                 string          `gorm:"size:450"`
	Email                    string          `gorm:"uniqueIndex;size:256"`
	PasswordHash             string          `gorm:"size:500"`
	PhoneNumber              string          `gorm:"size:50"`
	IsDeleted                bool            `gorm:"default:false"`
	Culture                  string          `gorm:"size:10"`
	RefreshTokenHash         string          `gorm:"size:500"`
	VendorID                 *uint
	DefaultShippingAddressID *uint
	DefaultBillingAddressID  *uint

	Roles          []Role          `gorm:"many2many:identity_user_roles;"`
	CustomerGroups []CustomerGroup `gorm:"many2many:identity_customer_group_users;"`
	Vendor         *Vendor         `gorm:"foreignKey:VendorID"`
}

func (User) TableName() string { return "identity_users" }
