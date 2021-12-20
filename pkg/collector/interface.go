package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ICollector interface {
	Name() string
	Settings() CollectorSettings
	Scan(context.Context) error
	StartLoop(context.Context, time.Duration)
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

type CollectorSettings struct {
	ClientConcurrency int64
	LoopInterval      time.Duration
}
