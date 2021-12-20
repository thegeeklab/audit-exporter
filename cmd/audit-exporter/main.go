package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/audit-exporter/pkg/server"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

// Version of current build
var Version = "devel"

type Instance struct {
}

func main() {
	settings := &server.Settings{}

	app := cli.NewApp()
	app.Name = "audit-exporter"
	app.Usage = "Prometheus exporter for various security tools."
	app.Version = Version
	app.Flags = globalFlags(settings)
	app.Action = run(settings)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(settings *server.Settings) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		inst := server.NewInstance(*settings)

		inst.Logger.Infof("Start monitor at %s", settings.Monitor.Address)
		monitor, err := server.NewMonitor(*settings, inst.Logger)
		if err != nil {
			inst.Logger.Fatal(xerrors.Errorf("Failed to create monitor: %w", err))
		}

		inst.AddProcessor(monitor)
		inst.Start()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM)
		<-quit
		inst.Logger.Infof("Attempt to shutdown instance")

		inst.Shutdown(context.Background())

		return nil
	}
}
