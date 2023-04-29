package metrics

import (
	"errors"
	"strings"
	"time"

	"github.com/armon/go-metrics"
)

type valueType interface {
	int32 | float32 | time.Duration
}

type metricType[V valueType] int

const (
	count        metricType[int32]         = 1
	gauge        metricType[float32]       = 2
	distribution metricType[float32]       = 3
	timing       metricType[time.Duration] = 4
)

type metricOptions struct {
	tags []string
}

func (mo metricOptions) Labels() []metrics.Label {
	labels := make([]metrics.Label, len(mo.tags))

	for i, tag := range mo.tags {
		parts := strings.Split(tag, ":")
		labels[i] = metrics.Label{Name: parts[0]}
		if len(parts) > 1 {
			labels[i].Value = parts[1]
		}
	}

	return labels
}

type metricOption func(*metricOptions)

// WithTags applies a set of tags to the metric being emitted. These are
// appended to the Client's persistent tags (see WithPersistentTags). Tags
// should be strings of the form key:value.
func WithTags(tags ...string) metricOption {
	return func(mo *metricOptions) {
		mo.tags = append(mo.tags, tags...)
	}
}

func applyOptions(opts ...metricOption) metricOptions {
	mo := metricOptions{tags: []string{}}
	for _, opt := range opts {
		opt(&mo)
	}
	return mo
}

// Emit emits a value for the provided metric Identifier using the provided
// Client.
func Emit[V valueType](client Client, metric Identifier[V], val V, opts ...metricOption) error {
	switch int(metric.mt) {
	case 1:
		return client.count(metric.name, int64(val), opts...)
	case 2:
		return client.gauge(metric.name, float32(val), opts...)
	case 3:
		return client.distribution(metric.name, float32(val), opts...)
	case 4:
		return client.timing(metric.name, time.Duration(val), opts...)
	default:
		return errors.New("invalid metric identifier")
	}
}

// EmitGlobal emits a value for the provided metric Identifier using the
// client configured for the package globally, see GlobalConfig.
func EmitGlobal[V valueType](metric Identifier[V], val V, opts ...metricOption) error {
	return Emit(global, metric, val, opts...)
}
