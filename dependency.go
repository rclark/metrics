package metrics

import "github.com/armon/go-metrics"

//go:generate mockery --name goMetrics --with-expecter --exported
type goMetrics interface {
	IncrCounterWithLabels(key []string, val float32, labels []metrics.Label)
	SetGaugeWithLabels(key []string, val float32, labels []metrics.Label)
	AddSampleWithLabels(key []string, val float32, labels []metrics.Label)
	Shutdown()
}
