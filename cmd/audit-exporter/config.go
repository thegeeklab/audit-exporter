package main

import (
	"github.com/urfave/cli/v2"
)

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Usage:   "source url to parse",
			EnvVars: []string{"URL_PARSER_URL"},
		},
	}
}
