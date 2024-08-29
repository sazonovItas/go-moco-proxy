package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	version string
	commit  string
	date    string
)

func cliVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Prints app version",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "short",
				Usage:   "show short version format",
				Aliases: []string{"s"},
				Value:   true,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println(ctx.App.Name, version)

			if !ctx.Bool("short") {
				fmt.Println("date:", date)
				fmt.Println("commit:", commit)
			}

			return nil
		},
	}
}
