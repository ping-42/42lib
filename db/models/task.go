package models

import "github.com/google/uuid"

type Task struct {
	ID                   uuid.UUID          `gorm:"type:uuid;primary_key;" json:"id"`
	TaskTypeID           uint64             //FK to TaskType.id
	TaskType             LvTaskType         `gorm:"foreignKey:TaskTypeID"`
	TaskStatusID         uint8              //FK to TaskType.id
	TaskStatus           LvTaskStatus       `gorm:"foreignKey:TaskStatusID"`
	SensorID             uuid.UUID          //FK to Sensor.id
	Sensor               Sensor             `gorm:"foreignKey:SensorID"`
	ClientSubscriptionID uint64             //FK to ClientSubscription.id
	ClientSubscription   ClientSubscription `gorm:"foreignKey:ClientSubscriptionID"`
	Opts                 []byte             `gorm:"type:jsonb"`
}

type LvTaskType struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Type string
}

type LvTaskStatus struct {
	ID     uint8 `gorm:"primaryKey;autoIncrement"`
	Status string
}
