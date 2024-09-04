package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	OrganizationID uuid.UUID    //FK to Client.id
	Organization   Organization `gorm:"foreignKey:OrganizationID"`

	TaskTypeID uint64     //FK to TaskType.id
	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`

	TestsCountSubscribed int
	TestsCountExecuted   int

	// the tasks should be assined on each X period
	Period time.Duration

	LastExecutionCompleted time.Time

	Opts []byte `gorm:"type:jsonb"`

	IsActive bool
}

// // TODO this needs to be deprecated. Now will user Subscription
// type ClientSubscription struct {
// 	ID uint64 `gorm:"primaryKey;autoIncrement"`

// 	ClientID uint64 //FK to Client.id
// 	Client   Client `gorm:"foreignKey:ClientID"`

// 	TaskTypeID uint64     //FK to TaskType.id
// 	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`

// 	TestsCountSubscribed int
// 	TestsCountExecuted   int

// 	// the tasks should be assined on each X period
// 	Period time.Duration

// 	LastExecutionCompleted time.Time

// 	Opts []byte `gorm:"type:jsonb"`

// 	IsActive bool
// }
