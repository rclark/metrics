package metrics_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rclark/metrics"
)

func ExampleNewClient() {
	host, _ := os.Hostname()

	client, err := metrics.NewClient(
		metrics.WithPersistentTags("env:testing", fmt.Sprintf("host:%s", host)),
		metrics.WithSinkAddress("127.0.0.1:8125"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}

func ExampleGlobalConfig() {
	host, _ := os.Hostname()

	err := metrics.GlobalConfig(
		metrics.WithPersistentTags("env:testing", fmt.Sprintf("host:%s", host)),
		metrics.WithSinkAddress("127.0.0.1:8125"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer metrics.GlobalClose()
}

func ExampleGlobalEmit() {
	var BackgroundJobsFailed = metrics.CountMetric("background.jobs.failed")

	err := metrics.GlobalEmit(
		BackgroundJobsFailed, 1,
		metrics.WithTags("job-name:failure"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleEmit() {
	client, err := metrics.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var BackgroundJobsFailed = metrics.CountMetric("background.jobs.failed")
	err = metrics.Emit(
		client, BackgroundJobsFailed, 1,
		metrics.WithTags("job-name:failure"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
