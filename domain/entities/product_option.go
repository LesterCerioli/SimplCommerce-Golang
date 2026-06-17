package entities

type ProductOption struct {
	BaseEntity
	Name string `gorm:"size:450;not null"`
}

func (ProductOption) TableName() string { return "catalog_product_options" }

type ProductOptionValue struct {
	BaseEntity
	OptionID    string `gorm:"type:uuid;not null"`
	ProductID   string `gorm:"type:uuid;not null"`
	Value       string `gorm:"size:450"`
	DisplayType string `gorm:"size:50"`
	SortIndex   int    `gorm:"default:0"`
	Option      ProductOption `gorm:"foreignKey:OptionID"`
	Product     Product       `gorm:"foreignKey:ProductID"`
}

func (ProductOptionValue) TableName() string { return "catalog_product_option_values" }

type ProductOptionCombination struct {
	BaseEntity
	ProductID string `gorm:"type:uuid;not null"`
	OptionID  string `gorm:"type:uuid;not null"`
	Value     string `gorm:"size:450"`
	SortIndex int    `gorm:"default:0"`
	Product   Product       `gorm:"foreignKey:ProductID"`
	Option    ProductOption `gorm:"foreignKey:OptionID"`
}

func (ProductOptionCombination) TableName() string { return "catalog_product_option_combinations" }
