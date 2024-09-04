package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;" json:"id"`
	TaskTypeID     uint64       //FK to TaskType.id
	TaskType       LvTaskType   `gorm:"foreignKey:TaskTypeID"`
	TaskStatusID   uint8        //FK to TaskType.id
	TaskStatus     LvTaskStatus `gorm:"foreignKey:TaskStatusID"`
	SensorID       uuid.UUID    //FK to Sensor.id
	Sensor         Sensor       `gorm:"foreignKey:SensorID"`
	SubscriptionID uint64       //FK to Subscription.id
	Subscription   Subscription `gorm:"foreignKey:SubscriptionID"`
	CreatedAt      time.Time    `gorm:"type:TIMESTAMPTZ;"`
	Opts           []byte       `gorm:"type:jsonb"`
}

type LvTaskType struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Type string
}

type LvTaskStatus struct {
	ID     uint8 `gorm:"primaryKey;autoIncrement"`
	Status string
}

func GetLatestSensorTasks(db *gorm.DB, sensorIds []string) (tasks []Task, err error) {
	err = db.Raw(`SELECT sensor_id
	,max(created_at) AS created_at
	FROM tasks
	WHERE sensor_id IN ?
	GROUP BY sensor_id;`, sensorIds).Scan(&tasks).Error
	return
}
