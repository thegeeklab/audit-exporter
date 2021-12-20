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

type IProcessor interface {
	Start() error
	Stop(context.Context) error
}

type instance struct {
	Settings   Settings
	processors []IProcessor
	Logger     logrus.Logger
}

// Settings defines the global available settings
type Settings struct {
	Monitor MonitorSettings
	Trivy   collector.CollectorSettings
}

// NewInstance creates a new prometheus exporter instance
func NewInstance(settings Settings) *instance {
	return &instance{
		Settings: settings,
		Logger:   *logrus.New(),
	}
}

func (inst *instance) AddProcessor(processor IProcessor) {
	inst.processors = append(inst.processors, processor)
}

func (inst *instance) Start() {
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

func (inst *instance) Shutdown(ctx context.Context) {
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
