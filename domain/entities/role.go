package entities

type Role struct {
	BaseEntity
	Name             string `gorm:"uniqueIndex;size:256"`
	NormalizedName   string `gorm:"uniqueIndex;size:256"`
	ConcurrencyStamp string `gorm:"size:256"`
	Users            []User `gorm:"many2many:identity_user_roles;"`
}

func (Role) TableName() string { return "identity_roles" }
