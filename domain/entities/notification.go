package entities

type Notification struct {
	BaseEntity
	UserID     string  `gorm:"type:uuid;not null;index"`
	Title      string  `gorm:"size:450;not null"`
	Body       string  `gorm:"type:text"`
	IsRead     bool    `gorm:"default:false"`
	EntityID   *string `gorm:"type:uuid"`
	EntityType string `gorm:"size:450"`
}

func (Notification) TableName() string { return "notifications_notifications" }
