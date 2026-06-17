package entities

type ProductTemplate struct {
	BaseEntity
	Name       string                          `gorm:"size:450;not null"`
	Attributes []ProductTemplateProductAttribute `gorm:"foreignKey:ProductTemplateID"`
}

func (ProductTemplate) TableName() string { return "catalog_product_templates" }

type ProductTemplateProductAttribute struct {
	ProductTemplateID  uint `gorm:"primaryKey"`
	ProductAttributeID uint `gorm:"primaryKey"`
}

func (ProductTemplateProductAttribute) TableName() string { return "catalog_product_template_product_attributes" }
