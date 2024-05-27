package ranker

type Ranker interface {
	Rank(WeightGetter) float64
}

// WeightGetter returns the specific type for ranking
type WeightGetter func() interface{}

const (
	// TypeError represents wrong WeightGetter type
	TypeError float64 = -1

	// RankZero is assigned when the sensor is at full capacity in some aspect
	RankZero float64 = 0
)

// RankEnvelope wraps all criteria that builds a sensor rank
type RankEnvelope struct {
	RuntimeRank      float64
	DistributionRank float64
}

// Rank calculates the final rank for a rank envelope, for all the criteria
func (re RankEnvelope) Rank(getWeights WeightGetter) float64 {
	weights, ok := getWeights().(SensorRankWeights)
	if !ok {
		return TypeError
	}

	return re.RuntimeRank * weights.RuntimeRank
}

// SensorRankWeights represent the calculation weights for the final rank
type SensorRankWeights struct {
	RuntimeRank float64
}

// DefaultFinalRankWeights is the default configuration for the final sensor ranking
var DefaultFinalRankWeights = SensorRankWeights{
	RuntimeRank: 1,
}

// HostRuntimeWeights represent the calculation weights based off host runtime
type HostRuntimeWeights struct {
	MaxCpuUsage    float64
	MaxMemUsage    float64
	CpuUsage       float64
	MemUsage       float64
	GoRoutineCount float64
}

// DefaultHostRuntimeWeights is the default configuration for the runtime ranking
var DefaultHostRuntimeWeights = HostRuntimeWeights{
	MaxCpuUsage:    95,
	MaxMemUsage:    95,
	CpuUsage:       0.4,
	MemUsage:       0.4,
	GoRoutineCount: 0.2,
}
