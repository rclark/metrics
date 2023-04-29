# metrics

A library for collecting [statsd metrics](https://github.com/statsd/statsd/blob/master/docs/metric_types.md).

## Install

```
go get github.com/rclark/metrics
```

## Configure a client

The library defines a client struct that you configure with
- a list of tags that should be applied to all metrics the client emits, and
- the address for a statsd sink.

You can choose whether to configure a package-global client, or pass the a client through your application.

```go
// To create a client struct to pass through your application:
client, err := metrics.NewClient(metrics.WithSinkAddress("0.0.0.0:8125"))

// Or to configure the global client:
err := metrics.GlobalConfig(metrics.WithSinkAddress("0.0.0.0:8125"))
```

## Define metric identifiers

In your codebase, define identifiers for each metric that your application generates. An identifier is just the metric's name, and its type. There are 4 types:

- Count: Used to track the number of events per unit time.
- Gauge: Used to track a measured value at a point in time.
- Distribution: Used to track the statistical distribution of a value over time.
- Timing: Used to track the statistical distribution of a duration over time.

Use a type-specific function to generate metric identifiers in your application. By doing this in one place in your program, you will have a catalog of all the metrics that are important to your application.

```go
var BackgroundJobsFailed = metrics.CountMetric("background.jobs.failed")
var BackgroundJobsDuration = metrics.TimingMetric("background.jobs.duration")
var BackgroundJobsQueueDepth = metrics.GaugeMetric("background.jobs.queue")
```

## Emit metrics

Use the `Emit()` function if you are passing a client through your application:

```go
err := metrics.Emit(client, BackgroundJobsFailed, 1)
```

Or use the `EmitGlobal()` function if you are relying on a package-global client:

```go
err := metrics.Emit(BackgroundJobsFailed, 1)
```
