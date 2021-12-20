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
	"github.com/thegeeklab/audit-exporter/pkg/client"
	"github.com/thegeeklab/audit-exporter/pkg/collector"
	"golang.org/x/net/netutil"
	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
)

const (
	metricsPath = "/metrics"
)

type MonitorSettings struct {
	Address               string
	MaxConnections        int64
	KeepAlived            bool
	ReUsePort             bool
	TCPKeepAliveInterval  time.Duration
	TrivyConcurrency      int64
	CollectorLoopInterval time.Duration
}

type Monitor struct {
	maxConnections int64
	listener       net.Listener
	server         *http.Server
}

func NewMonitor(settings Settings, logger logrus.Logger) (*Monitor, error) {
	router := mux.NewRouter()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	trivyCollector := collector.NewTrivyCollector(
		client.TrivyClient{},
		settings.Monitor.TrivyConcurrency,
		logger,
	)
	registry.MustRegister(trivyCollector)
	ctx := context.Background()
	if err := trivyCollector.Scan(ctx); err != nil {
		return nil, xerrors.Errorf("failed scan of trivy collector: %w", err)
	}
	trivyCollector.StartLoop(ctx, settings.Monitor.CollectorLoopInterval)

	router.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	var listener net.Listener
	var err error
	if settings.Monitor.ReUsePort {
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
			KeepAlive: settings.Monitor.TCPKeepAliveInterval,
		}
		listener, err = listenConfig.Listen(context.Background(), "tcp", settings.Monitor.Address)
	} else {
		listener, err = net.Listen("tcp", settings.Monitor.Address)
	}

	if err != nil {
		return nil, xerrors.Errorf("could not listen %s: %w", settings.Monitor.Address, err)
	}

	server := &http.Server{
		Handler: router,
	}
	server.SetKeepAlivesEnabled(settings.Monitor.KeepAlived)
	return &Monitor{
		maxConnections: settings.Monitor.MaxConnections,
		listener:       listener,
		server:         server,
	}, nil
}

func (m *Monitor) Start() error {
	return m.server.Serve(netutil.LimitListener(m.listener, int(m.maxConnections)))
}

func (m *Monitor) Stop(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
