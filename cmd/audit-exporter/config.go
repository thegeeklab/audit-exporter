package main

import (
	"github.com/thegeeklab/audit-exporter/pkg/server"
	"github.com/urfave/cli/v2"
)

func globalFlags(settings *server.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Usage:       "source url to parse",
			EnvVars:     []string{"URL_PARSER_URL"},
			Destination: &settings.test,
		},
	}
}
