package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ckeyer/bat"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	var addr string
	var debug bool

	app := &cli.App{
		Name:    "bat",
		Version: "0.1",
		Usage:   "bat [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				EnvVars:     []string{"ADDR"},
				Destination: &addr,
				DefaultText: ":8080",
				Value:       ":8080",
				Usage:       "server listenning addr, [host]:port",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"D"},
				Usage:       "debug model",
				Destination: &debug,
				DefaultText: "false",
				Value:       false,
			},
		},
		Before: func(ctx *cli.Context) error {
			log.SetFormatter(&log.JSONFormatter{})
			if debug {
				log.SetLevel(log.DebugLevel)
			}
			log.Debug("cli.main starting.")
			return nil
		},
		Action: func(ctx *cli.Context) error {
			bat.Serve(addr)
			return nil
		},
	}

	app.Run(os.Args)
}
