package metrics

import (
	"time"

	"github.com/armon/go-metrics"
)

var global, _ = NewClient()

// Client is configured to emit metrics to a statsd sink.
type Client struct {
	addr   string
	tags   []string
	statsd *metrics.Metrics
}

func (c Client) count(name string, val int64, opts ...metricOption) error {
	mo := applyOptions(opts...)
	c.statsd.IncrCounterWithLabels([]string{name}, float32(val), mo.Labels())
	return nil
}

func (c Client) gauge(name string, val float32, opts ...metricOption) error {
	mo := applyOptions(opts...)
	c.statsd.SetGaugeWithLabels([]string{name}, val, mo.Labels())
	return nil
}

func (c Client) distribution(name string, val float32, opts ...metricOption) error {
	mo := applyOptions(opts...)
	c.statsd.AddSampleWithLabels([]string{name}, val, mo.Labels())
	return nil
}

func (c Client) timing(name string, val time.Duration, opts ...metricOption) error {
	return c.distribution(name, float32(val), opts...)
}

type clientOption func(*Client)

// WithPersistentTags configures a set of tags that should be applied to all
// metrics emitted by the Client. Tags should be strings of the form key:value.
func WithPersistentTags(tags ...string) clientOption {
	return func(c *Client) {
		c.tags = tags
	}
}

// WithSinkAddress sets the address of the statsd sink to which the Client will
// emit metrics. The default address is 127.0.0.1:8125.
func WithSinkAddress(addr string) clientOption {
	return func(c *Client) {
		c.addr = addr
	}
}

// NewClient creates a new Client.
func NewClient(opts ...clientOption) (Client, error) {
	var none Client

	c := Client{tags: []string{}, addr: "127.0.0.1:8125"}
	for _, opt := range opts {
		opt(&c)
	}

	sink, err := metrics.NewStatsdSink(c.addr)
	if err != nil {
		return none, err
	}

	statsd, err := metrics.New(&metrics.Config{}, sink)
	if err != nil {
		return none, err
	}

	c.statsd = statsd
	return c, nil
}

// GlobalConfig configures the package's global client, allowing for the use
// of the GlobalCollect function.
func GlobalConfig(opts ...clientOption) (err error) {
	global, err = NewClient(opts...)
	return
}
