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

type XXXXX struct {
	test  string
	test2 string
}

type Instance struct {
	settings   Settings
	processors []IProcessor
	logger     logrus.Logger
}

type Settings struct {
	monitor MonitorSettings
	test    string
}

func NewInstance(settings Settings) *Instance {
	return &Instance{
		settings: settings,
		logger:   *logrus.New(),
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
					inst.logger.Errorf("panic: %+v\n", err)
					inst.logger.Debugf("%s\n", debug.Stack())
				}
			}()
			if err := processor.Start(); err != nil && err != http.ErrServerClosed {
				inst.logger.Errorf("Failed to listen: %s\n", err)
			}
		}(processor)
	}
}

func (inst *Instance) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(gracePeriod)*time.Second)
	defer cancel()
	for _, p := range inst.processors {
		if err := p.Stop(ctx); err != nil {
			inst.logger.Errorf("Failed to shutdown: %+v\n", err)
		}
	}
	select {
	case <-ctx.Done():
		inst.logger.Infof("Instance shutdown timed out in %d seconds\n", gracePeriod)
	default:
	}
	inst.logger.Infof("Instance has been shutdown\n")
}
