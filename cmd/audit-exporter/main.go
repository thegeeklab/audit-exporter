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

		monitor, err := server.NewMonitor(*settings)
		if err != nil {
			return xerrors.Errorf("failed to create monitor: %w", err)
		}

		inst.AddProcessor(monitor)
		inst.Start()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM)
		<-quit
		inst.logger.Infof("Attempt to shutdown instance...\n")

		inst.Shutdown(context.Background())

		return nil
	}
}
