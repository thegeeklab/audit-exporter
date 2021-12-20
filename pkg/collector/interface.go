package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ICollector defines the interface that need to be implemented by each collector
type ICollector interface {
	Name() string
	Settings() Settings
	Scan(context.Context) error
	StartLoop(context.Context, time.Duration)
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

// Settings defines the required collector settings
type Settings struct {
	ClientConcurrency int64
	LoopInterval      time.Duration
}
