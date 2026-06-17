package entities

type Media struct {
	BaseEntity
	Caption   string `gorm:"size:450"`
	FileSize  int64
	FileName  string `gorm:"size:450"`
	MediaType int    `gorm:"default:0"`
}

func (Media) TableName() string { return "identity_media" }
