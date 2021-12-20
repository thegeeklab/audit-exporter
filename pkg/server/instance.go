package server

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	gracePeriod = 10
)

type IProcessor interface {
	Start() error
	Stop(context.Context) error
}

type Instance struct {
	Settings   Settings
	processors []IProcessor
	Logger     logrus.Logger
}

type Settings struct {
	Monitor MonitorSettings
}

func NewInstance(settings Settings) *Instance {
	return &Instance{
		Settings: settings,
		Logger:   *logrus.New(),
	}
}

func (inst *Instance) AddProcessor(processor IProcessor) {
	inst.processors = append(inst.processors, processor)
}

func (inst *Instance) Start() {
	for _, processor := range inst.processors {
		go func(processor IProcessor) {
			defer func() {
				if err := recover(); err != nil {
					inst.Logger.Errorf("panic: %+v\n", err)
					inst.Logger.Debugf("%s\n", debug.Stack())
				}
			}()
			if err := processor.Start(); err != nil && err != http.ErrServerClosed {
				inst.Logger.Errorf("Failed to listen: %s\n", err)
			}
		}(processor)
	}
}

func (inst *Instance) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(gracePeriod)*time.Second)
	defer cancel()
	for _, p := range inst.processors {
		if err := p.Stop(ctx); err != nil {
			inst.Logger.Errorf("Failed to shutdown: %+v\n", err)
		}
	}
	select {
	case <-ctx.Done():
		inst.Logger.Infof("Instance shutdown timed out in %d seconds\n", gracePeriod)
	default:
	}
	inst.Logger.Infof("Instance has been shutdown\n")
}
