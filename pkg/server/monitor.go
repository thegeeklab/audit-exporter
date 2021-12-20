package server

import (
	"context"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/audit-exporter/pkg/collector"
	"golang.org/x/net/netutil"
	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
)

const (
	metricsPath = "/metrics"
)

// MonitorSettings defines required attributes
type MonitorSettings struct {
	Address              string
	MaxConnections       int64
	KeepAlived           bool
	ReUsePort            bool
	TCPKeepAliveInterval time.Duration
}

// Monitor defines the http metrics server
type Monitor struct {
	maxConnections int64
	listener       net.Listener
	server         *http.Server
}

// NewMonitor creates a new Monitor instance providing http metrics server
func NewMonitor(settings MonitorSettings, logger *logrus.Logger, mcollectors []collector.ICollector) (*Monitor, error) {
	router := mux.NewRouter()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	for _, collector := range mcollectors {
		registry.MustRegister(collector)
		ctx := context.Background()
		if err := collector.Scan(ctx); err != nil {
			return nil, xerrors.Errorf("failed scan of %s collector: %w", collector.Name(), err)
		}
		collector.StartLoop(ctx, collector.Settings().LoopInterval)
	}

	router.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	var listener net.Listener
	var err error
	if settings.ReUsePort {
		listenConfig := &net.ListenConfig{
			Control: func(network string, address string, c syscall.RawConn) error {
				var innerErr error
				if err := c.Control(func(s uintptr) {
					innerErr = unix.SetsockoptInt(int(s), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				}); err != nil {
					return err
				}
				if innerErr != nil {
					return innerErr
				}
				return nil
			},
			KeepAlive: settings.TCPKeepAliveInterval,
		}
		listener, err = listenConfig.Listen(context.Background(), "tcp", settings.Address)
	} else {
		listener, err = net.Listen("tcp", settings.Address)
	}

	if err != nil {
		return nil, xerrors.Errorf("could not listen %s: %w", settings.Address, err)
	}

	server := &http.Server{
		Handler: router,
	}
	server.SetKeepAlivesEnabled(settings.KeepAlived)
	return &Monitor{
		maxConnections: settings.MaxConnections,
		listener:       listener,
		server:         server,
	}, nil
}

// Start starts the http metrics server of the Monitor
func (m *Monitor) Start() error {
	return m.server.Serve(netutil.LimitListener(m.listener, int(m.maxConnections)))
}

// Stop stops the http metrics server of the Monitor
func (m *Monitor) Stop(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
