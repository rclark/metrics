package metrics

import (
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/rclark/metrics/mocks"
	"github.com/stretchr/testify/require"
)

func TestEmit(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err, "NewClient should not fail")

	t.Cleanup(func() {
		err := client.Close()
		require.NoError(t, err, "Close should not fail")
	})

	mockGoMetrics := func(statsd *mocks.GoMetrics) func() {
		original := client.statsd
		client.statsd = statsd
		return func() {
			client.statsd = original
		}
	}

	t.Run("count", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := Emit(client, CountMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("gauge", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			SetGaugeWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := Emit(client, GaugeMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("distribution", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			AddSampleWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := Emit(client, DistributionMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("timing", func(t *testing.T) {
		duration := time.Since(time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC))
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			AddSampleWithLabels([]string{"name"}, float32(duration), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := Emit(client, TimingMetric("name"), duration)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("metric options", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{{Name: "key", Value: "val"}}).
			Return()
		defer mockGoMetrics(statsd)()

		err := Emit(client, CountMetric("name"), 1, WithTags("key:val"))
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("client options", func(t *testing.T) {
		client, err := NewClient(
			WithPersistentTags("key:val"),
			WithSinkAddress("1.2.3.4:5678"),
		)
		require.NoError(t, err, "NewClient with options should not fail")

		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{{Name: "key", Value: "val"}}).
			Return()
		client.statsd = statsd

		err = Emit(client, CountMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})
}

func TestGlobalEmit(t *testing.T) {
	err := GlobalConfig()
	require.NoError(t, err, "NewClient should not fail")

	t.Cleanup(func() {
		global, err = NewClient()
		require.NoError(t, err, "resetting global client should not fail")
	})

	mockGoMetrics := func(statsd *mocks.GoMetrics) func() {
		original := global.statsd
		global.statsd = statsd
		return func() {
			global.statsd = original
		}
	}

	t.Run("count", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := GlobalEmit(CountMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("gauge", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			SetGaugeWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := GlobalEmit(GaugeMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("distribution", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			AddSampleWithLabels([]string{"name"}, float32(1), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := GlobalEmit(DistributionMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("timing", func(t *testing.T) {
		duration := time.Since(time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC))
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			AddSampleWithLabels([]string{"name"}, float32(duration), []metrics.Label{}).
			Return()
		defer mockGoMetrics(statsd)()

		err := GlobalEmit(TimingMetric("name"), duration)
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("metric options", func(t *testing.T) {
		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{{Name: "key", Value: "val"}}).
			Return()
		defer mockGoMetrics(statsd)()

		err := GlobalEmit(CountMetric("name"), 1, WithTags("key:val"))
		require.NoError(t, err, "emit should not fail")
	})

	t.Run("client options", func(t *testing.T) {
		err := GlobalConfig(
			WithPersistentTags("key:val"),
			WithSinkAddress("1.2.3.4:5678"),
		)
		require.NoError(t, err, "NewClient with options should not fail")

		statsd := mocks.NewGoMetrics(t)
		statsd.EXPECT().
			IncrCounterWithLabels([]string{"name"}, float32(1), []metrics.Label{{Name: "key", Value: "val"}}).
			Return()
		global.statsd = statsd

		err = GlobalEmit(CountMetric("name"), 1)
		require.NoError(t, err, "emit should not fail")
	})
}
