package main

import (
	"math"
	"time"

	"github.com/thegeeklab/audit-exporter/pkg/server"
	"github.com/urfave/cli/v2"
)

func globalFlags(settings *server.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "monitor-address",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_MONITOR_ADDRESS"},
			Value:       "127.0.0.1:9000",
			Destination: &settings.Monitor.Address,
		},
		&cli.Int64Flag{
			Name:        "monitor-max-connections",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_MONITOR_MAX_CONNECTIONS"},
			Value:       math.MaxInt64,
			Destination: &settings.Monitor.MaxConnections,
		},
		&cli.BoolFlag{
			Name:        "keepalived",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_KEEPALIVED"},
			Value:       true,
			Destination: &settings.Monitor.KeepAlived,
		},
		&cli.BoolFlag{
			Name:        "reuse-port",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_REUSE_PORT"},
			Value:       false,
			Destination: &settings.Monitor.ReUsePort,
		},
		&cli.DurationFlag{
			Name:        "tcp-keepalive-interval",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_TCP_KEEPALIVE_INTERVAL"},
			Value:       0,
			Destination: &settings.Monitor.TCPKeepAliveInterval,
		},
		&cli.Int64Flag{
			Name:        "trivy-concurrency",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_TRIVY_CONCURRENCY"},
			Value:       10,
			Destination: &settings.Monitor.TrivyConcurrency,
		},
		&cli.DurationFlag{
			Name:        "collector-loop-interval",
			Usage:       "source url to parse",
			EnvVars:     []string{"AUDIT_EXPORTER_COLLECTOR_LOOP_INTERVAL"},
			Value:       time.Second * 60,
			Destination: &settings.Monitor.CollectorLoopInterval,
		},
	}
}
