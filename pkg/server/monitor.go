package server

import (
	"context"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thegeeklab/audit-exporter/client"
	"github.com/thegeeklab/audit-exporter/collector"
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

func NewMonitor(settings MonitorSettings) (*Monitor, error) {
	router := mux.NewRouter()

	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())
	trivyCollector := collector.NewTrivyCollector(
		&client.TrivyClient{},
		settings.TrivyConcurrency,
	)
	registry.MustRegister(trivyCollector)
	ctx := context.Background()
	if err := trivyCollector.Scan(ctx); err != nil {
		return nil, xerrors.Errorf("failed to scan of trivy collector: %w", err)
	}
	trivyCollector.StartLoop(ctx, settings.CollectorLoopInterval)

	router.Handle(metricsPath, promhttp.Handler())

	var listener net.Listener
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
		Handler: &promhttp.Handler{
			Handler: router,
			FormatSpanName: func(r *http.Request) string {
				return r.Method + " " + r.URL.Path
			},
		},
	}
	server.SetKeepAlivesEnabled(settings.KeepAlived)
	return &Monitor{
		maxConnections: settings.MaxConnections,
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
