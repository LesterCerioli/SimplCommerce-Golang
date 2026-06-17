package entities

type Activity struct {
	BaseEntity
	ActivityTypeID string       `gorm:"type:uuid;not null;index"`
	UserID         string       `gorm:"type:uuid;index"`
	EntityID       uint
	EntityTypeID   string       `gorm:"size:450"`
	ActivityType   ActivityType `gorm:"foreignKey:ActivityTypeID"`
}

func (Activity) TableName() string { return "activity_log_activities" }

type ActivityType struct {
	BaseEntity
	Name       string     `gorm:"size:450;uniqueIndex;not null"`
	Activities []Activity `gorm:"foreignKey:ActivityTypeID"`
}

func (ActivityType) TableName() string { return "activity_log_activity_types" }
