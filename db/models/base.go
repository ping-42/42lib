package models

import (
	"time"

	"github.com/google/uuid"
)

type TsSensorTaskBase struct {
	Time     time.Time `gorm:"type:TIMESTAMPTZ;"`
	SensorID uuid.UUID `gorm:"type:uuid;"`
	TaskID   uuid.UUID `gorm:"type:uuid;"`
}

type LvProtocol struct {
	ID   int64
	Type string
}
