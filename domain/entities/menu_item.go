package entities

type MenuItem struct {
	BaseEntity
	ParentID     *string `gorm:"type:uuid"`
	MenuID       string  `gorm:"type:uuid;not null"`
	EntityID     *string `gorm:"type:uuid;index"`
	CustomLink   string `gorm:"size:450"`
	Name         string `gorm:"size:450;not null"`
	DisplayOrder int    `gorm:"default:0"`
	Parent       *MenuItem `gorm:"foreignKey:ParentID"`
	Menu         Menu      `gorm:"foreignKey:MenuID"`
}

func (MenuItem) TableName() string { return "cms_menu_items" }
