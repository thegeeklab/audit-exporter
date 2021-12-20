package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/audit-exporter/pkg/client"
	"github.com/thegeeklab/audit-exporter/pkg/collector"
	"github.com/thegeeklab/audit-exporter/pkg/server"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

// Version of current build
var Version = "devel"

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
		collectors := []collector.ICollector{}
		inst := server.NewInstance(*settings)
		inst.Logger.SetFormatter(&logrus.TextFormatter{
			// DisableColors: true,
			FullTimestamp: true,
		})
		inst.Logger.Infof("Start monitor at %s", settings.Monitor.Address)

		trivyCollector := collector.NewTrivyCollector(
			client.TrivyClient{},
			settings.Trivy,
			&inst.Logger,
		)

		collectors = append(collectors, trivyCollector)

		monitor, err := server.NewMonitor(settings.Monitor, &inst.Logger, collectors)
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
