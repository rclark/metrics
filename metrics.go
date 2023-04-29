package metrics

import "time"

// Identifier holds the name and type of a single metric.
type Identifier[V valueType] struct {
	mt   metricType[V]
	name string
}

func new[V valueType](name string, mt metricType[V]) Identifier[V] {
	return Identifier[V]{mt: mt, name: name}
}

// CountMetric creates a metric Identifier for a count metric.
func CountMetric(name string) Identifier[int32] {
	return new(name, count)
}

// GaugeMetric creates a metric Identifier for a gauge metric.
func GaugeMetric(name string) Identifier[float32] {
	return new(name, gauge)
}

// DistributionMetric creates a metric Identifier for a distribution metric.
func DistributionMetric(name string) Identifier[float32] {
	return new(name, distribution)
}

// TimingMetric creates a metric Identifier for a distribution metric that
// captures durations.
func TimingMetric(name string) Identifier[time.Duration] {
	return new(name, timing)
}
