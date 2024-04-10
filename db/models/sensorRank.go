package models

import (
	"time"

	"github.com/google/uuid"
)

type SensorRank struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SensorID uuid.UUID //FK to Client.id
	Sensor   Sensor    `gorm:"foreignKey:SensorID"`

	CurrentRank int8

	CreatedAt time.Time
}
