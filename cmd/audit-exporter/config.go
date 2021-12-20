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
			Usage:       "bind address for the metrics server",
			EnvVars:     []string{"AUDIT_EXPORTER_MONITOR_ADDRESS"},
			Value:       "127.0.0.1:9000",
			Destination: &settings.Monitor.Address,
		},
		&cli.Int64Flag{
			Name:        "monitor-max-connections",
			Usage:       "max connection limit for the metrics server",
			EnvVars:     []string{"AUDIT_EXPORTER_MONITOR_MAX_CONNECTIONS"},
			Value:       math.MaxInt64,
			Destination: &settings.Monitor.MaxConnections,
		},
		&cli.BoolFlag{
			Name:        "monitor-keepalived",
			Usage:       "enable keepalived for the metrics server",
			EnvVars:     []string{"AUDIT_EXPORTER_KEEPALIVED"},
			Value:       true,
			Destination: &settings.Monitor.KeepAlived,
		},
		&cli.BoolFlag{
			Name:        "monitor-reuse-port",
			Usage:       "enable the SO_REUSEPORT socket option for the metrics server",
			EnvVars:     []string{"AUDIT_EXPORTER_REUSE_PORT"},
			Value:       false,
			Destination: &settings.Monitor.ReUsePort,
		},
		&cli.DurationFlag{
			Name:        "monitor-tcp-keepalive-interval",
			Usage:       "set keepalived interval for the metrics server",
			EnvVars:     []string{"AUDIT_EXPORTER_TCP_KEEPALIVE_INTERVAL"},
			Value:       0,
			Destination: &settings.Monitor.TCPKeepAliveInterval,
		},

		// Settings for trivy collector
		&cli.Int64Flag{
			Name:        "trivy-concurrency",
			Usage:       "max concurrent trivy clients used to scan images",
			EnvVars:     []string{"AUDIT_EXPORTER_TRIVY_CONCURRENCY"},
			Value:       10,
			Destination: &settings.Trivy.ClientConcurrency,
		},
		&cli.DurationFlag{
			Name:        "trivy-loop-interval",
			Usage:       "interval to schedule trivy scans",
			EnvVars:     []string{"AUDIT_EXPORTER_TRIVY_LOOP_INTERVAL"},
			Value:       time.Second * 60,
			Destination: &settings.Trivy.LoopInterval,
		},
	}
}
