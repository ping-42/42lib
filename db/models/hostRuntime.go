package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/constants"
	"github.com/ping-42/42lib/ranker"
	"gorm.io/gorm"
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
}

// Rank implements the Ranker interface for the runtime statistics
func (hrs TsHostRuntimeStat) Rank(getWeights ranker.WeightGetter) float64 {
	weights, ok := getWeights().(ranker.HostRuntimeWeights)
	if !ok {
		return ranker.TypeError
	}

	// sensor at full capacity
	if hrs.MemUsedPercent > weights.MaxMemUsage ||
		hrs.CpuUsage > weights.MaxCpuUsage ||
		hrs.GoRoutineCount >= constants.MaxGoRoutineCountSensor {
		return ranker.RankZero
	}

	grUsage := (float64(hrs.GoRoutineCount) / constants.MaxGoRoutineCountSensor) * 100
	cpuFree := 100 - hrs.CpuUsage
	memFree := 100 - hrs.MemUsedPercent
	grFree := 100 - grUsage

	return (weights.CpuUsage * cpuFree) +
		(weights.MemUsage * memFree) +
		(weights.GoRoutineCount * grFree)
}

// GetRuntimeStats returns the results grouped by sensor id and the given interval
func GetRuntimeStats(db *gorm.DB, durationMinutes int, sensorIds []string) (stats []TsHostRuntimeStat, err error) {
	err = db.Raw(`SELECT sensor_id
		,CAST(avg(go_routine_count) AS INTEGER) AS go_routine_count
		,avg(cpu_usage) as cpu_usage
		,avg(mem_used_percent) as mem_used_percent
	FROM ts_host_runtime_stats
	WHERE "time" > now() - INTERVAL ?
	AND sensor_id IN ?
	GROUP BY sensor_id;`,
		gorm.Expr(fmt.Sprintf("'%v minutes'", durationMinutes)),
		sensorIds).Scan(&stats).Error
	return
}
