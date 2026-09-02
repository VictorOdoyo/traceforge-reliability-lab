package domain

type SamplingPolicy struct {
	DefaultRate         float64 `json:"defaultRate"`
	ErrorRate           float64 `json:"errorRate"`
	SlowTraceRate       float64 `json:"slowTraceRate"`
	SlowThresholdMillis float64 `json:"slowThresholdMillis"`
}

func (p SamplingPolicy) RateFor(trace Trace) float64 {
	if trace.HasError {
		return clampRate(p.ErrorRate)
	}
	if float64(trace.Duration.Milliseconds()) >= p.SlowThresholdMillis {
		return clampRate(p.SlowTraceRate)
	}
	return clampRate(p.DefaultRate)
}

func clampRate(rate float64) float64 {
	if rate < 0 {
		return 0
	}
	if rate > 1 {
		return 1
	}
	return rate
}
