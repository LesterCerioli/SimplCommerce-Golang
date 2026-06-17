package entities

type Notification struct {
	BaseEntity
	UserID     uint   `gorm:"not null;index"`
	Title      string `gorm:"size:450;not null"`
	Body       string `gorm:"type:text"`
	IsRead     bool   `gorm:"default:false"`
	EntityID   *uint
	EntityType string `gorm:"size:450"`
}

func (Notification) TableName() string { return "notifications_notifications" }
