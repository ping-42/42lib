package models

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID    `gorm:"type:uuid"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Name           string
	Location       string
	Secret         string
	IsActive       bool
	CreatedAt      time.Time
}

// type SensorSupportedTaskTypes struct {
// 	SensorID   uuid.UUID  `gorm:"type:uuid"` //FK to Sensor.id
// 	Sensor     Sensor     `gorm:"foreignKey:SensorID"`
// 	TaskTypeID uint64     //FK to TaskType.id
// 	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`
// }
