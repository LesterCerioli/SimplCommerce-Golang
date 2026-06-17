package entities

type ProductAttribute struct {
	BaseEntity
	Name    string `gorm:"size:450;not null"`
	GroupID string `gorm:"type:uuid"`
	Group   ProductAttributeGroup `gorm:"foreignKey:GroupID"`
}

func (ProductAttribute) TableName() string { return "catalog_product_attributes" }

type ProductAttributeGroup struct {
	BaseEntity
	Name       string             `gorm:"size:450;not null"`
	Attributes []ProductAttribute `gorm:"foreignKey:GroupID"`
}

func (ProductAttributeGroup) TableName() string { return "catalog_product_attribute_groups" }

type ProductAttributeValue struct {
	BaseEntity
	AttributeID string `gorm:"type:uuid;not null"`
	ProductID   string `gorm:"type:uuid;not null"`
	Value       string `gorm:"type:text"`
	Attribute   ProductAttribute `gorm:"foreignKey:AttributeID"`
	Product     Product          `gorm:"foreignKey:ProductID"`
}

func (ProductAttributeValue) TableName() string { return "catalog_product_attribute_values" }
