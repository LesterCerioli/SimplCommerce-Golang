package models

import "github.com/simplcommerce-go/pkg/models"

type Activity struct {
	models.BaseEntity
	ActivityTypeID uint         `gorm:"not null;index"`
	UserID         uint         `gorm:"index"`
	EntityID       uint
	EntityTypeID   string       `gorm:"size:450"`
	ActivityType   ActivityType `gorm:"foreignKey:ActivityTypeID"`
}

func (Activity) TableName() string { return "activity_log_activities" }

type ActivityType struct {
	models.BaseEntity
	Name       string     `gorm:"size:450;uniqueIndex;not null"`
	Activities []Activity `gorm:"foreignKey:ActivityTypeID"`
}

func (ActivityType) TableName() string { return "activity_log_activity_types" }
