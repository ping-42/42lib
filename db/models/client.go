package models

type Client struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Email string
}
