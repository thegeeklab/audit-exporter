package server

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/audit-exporter/pkg/collector"
)

const (
	gracePeriod = 10
)

// IProcessor defines the interface that need to be implemented by each monitor
type IProcessor interface {
	Start() error
	Stop(context.Context) error
}

// Instance defines the exporter instance
type Instance struct {
	Settings   Settings
	processors []IProcessor
	Logger     logrus.Logger
}

// Settings defines the global server settings
type Settings struct {
	Monitor MonitorSettings
	Trivy   collector.Settings
}

// NewInstance creates a new exporter instance
func NewInstance(settings Settings) *Instance {
	return &Instance{
		Settings: settings,
		Logger:   *logrus.New(),
	}
}

// AddProcessor allows to add monitors to the exporter instance
func (inst *Instance) AddProcessor(processor IProcessor) {
	inst.processors = append(inst.processors, processor)
}

// Start brings up the http server for each monitor
func (inst *Instance) Start() {
	for _, processor := range inst.processors {
		go func(processor IProcessor) {
			defer func() {
				if err := recover(); err != nil {
					inst.Logger.Errorf("panic: %+v", err)
					inst.Logger.Debugf("%s", debug.Stack())
				}
			}()
			if err := processor.Start(); err != nil && err != http.ErrServerClosed {
				inst.Logger.Errorf("Failed to listen: %s", err)
			}
		}(processor)
	}
}

// Shutdown gracefully shutdown of the http server for each monitor
func (inst *Instance) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(gracePeriod)*time.Second)
	defer cancel()
	for _, p := range inst.processors {
		if err := p.Stop(ctx); err != nil {
			inst.Logger.Errorf("Failed to shutdown: %+v", err)
		}
	}
	select {
	case <-ctx.Done():
		inst.Logger.Infof("Instance shutdown timed out in %d seconds", gracePeriod)
	default:
	}
	inst.Logger.Infof("Instance has been shutdown")
}
