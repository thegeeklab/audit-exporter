package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version of current build
var Version = "devel"

func main() {
	app := cli.NewApp()
	app.Name = "audit-exporter"
	app.Usage = "Prometheus exporter for various security tools."
	app.Version = Version
	app.Flags = globalFlags()

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
