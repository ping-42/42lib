package models

type TsHostRuntimeStat struct {
	TsSensorTaskBase
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
