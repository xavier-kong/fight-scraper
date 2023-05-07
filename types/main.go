package types

type Event struct {
	ID uint
	Name string `gorm:"size: 255; not null;" json:"name"`
	TimestampSeconds int `gorm:"type: numeric; not null;" json:"timestamp_seconds"`
	Headline string `gorm:"size: 255; not null;" json:"headline"`
	Url string `gorm:"size: 255; not null;" json:"url"`
	Org string `gorm:"size: 255; not null;" json:"org"`
}
