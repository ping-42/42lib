package sensorTask

type HostTelemetry struct {
	Cpu        Cpu
	Memory     Memory
	GoRoutines int
	Network    []Network
}

type Cpu struct {
	ModelName string
	Cores     uint16
	CpuUsage  float64
}

type Memory struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type Network struct {
	Name      string
	BytesSent uint64
	BytesRecv uint64
}
