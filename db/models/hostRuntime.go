package models

import (
	"time"

	"github.com/google/uuid"
)

type TsHostRuntimeStat struct {
	Time           time.Time `gorm:"type:TIMESTAMPTZ;"`
	SensorID       uuid.UUID `gorm:"type:uuid;"`
	GoRoutineCount int
	//
	CpuCores     uint16
	CpuUsage     float64
	CpuModelName string
	//
	MemTotal       uint64
	MemUsed        uint64
	MemFree        uint64
	MemUsedPercent float64
	//
	Network string
}
