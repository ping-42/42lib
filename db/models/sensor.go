package models

import "github.com/google/uuid"

type Sensor struct {
	ID       uuid.UUID `gorm:"primaryKey;autoIncrement"`
	Name     string
	Location string
	Secret   string
}

// type SensorSupportedTaskTypes struct {
// 	SensorID   uuid.UUID  //FK to Sensor.id
// 	Sensor     Sensor     `gorm:"foreignKey:SensorID"`
// 	TaskTypeID uint64     //FK to TaskType.id
// 	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`
// }
