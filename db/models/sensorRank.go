package models

import (
	"time"

	"github.com/google/uuid"
)

type SensorRank struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SensorID uuid.UUID `gorm:"type:uuid"` //FK to Client.id
	Sensor   Sensor    `gorm:"foreignKey:SensorID"`

	Rank             float64
	DistributionRank float64

	CreatedAt time.Time
}
