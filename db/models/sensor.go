package models

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID //FK to User.id // TODO: since User is created after Sensor the migrations failing...
	User      User      `gorm:"foreignKey:UserID"`
	Name      string
	Location  string
	Secret    string
	IsActive  bool
	CreatedAt time.Time
}

// type SensorSupportedTaskTypes struct {
// 	SensorID   uuid.UUID  //FK to Sensor.id
// 	Sensor     Sensor     `gorm:"foreignKey:SensorID"`
// 	TaskTypeID uint64     //FK to TaskType.id
// 	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`
// }
