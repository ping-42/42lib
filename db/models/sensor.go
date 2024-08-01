package models

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	OrganisationID uuid.UUID
	Organisation   Organisation `gorm:"foreignKey:OrganisationID"`
	Name           string
	Location       string
	Secret         string
	IsActive       bool
	CreatedAt      time.Time
}

// type SensorSupportedTaskTypes struct {
// 	SensorID   uuid.UUID  //FK to Sensor.id
// 	Sensor     Sensor     `gorm:"foreignKey:SensorID"`
// 	TaskTypeID uint64     //FK to TaskType.id
// 	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`
// }
