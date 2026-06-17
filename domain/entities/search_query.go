package entities

type SearchQuery struct {
	BaseEntity
	QueryText    string `gorm:"size:500;not null"`
	ResultsCount int    `gorm:"default:0"`
}

func (SearchQuery) TableName() string { return "search_queries" }
