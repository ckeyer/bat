package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ckeyer/bat"
	"github.com/ckeyer/bat/store"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	var addr, redisAddr, redisAuth string
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
				EnvVars:     []string{"DEBUG"},
				Usage:       "debug model",
				Destination: &debug,
				DefaultText: "false",
				Value:       false,
			},

			// db: redis
			&cli.StringFlag{
				Name:        "redis_addr",
				EnvVars:     []string{"REDIS_ADDR", "REDIS"},
				Destination: &redisAddr,
				DefaultText: "",
				Value:       "",
				Usage:       "redis db addr, [host]:port",
			},
			&cli.StringFlag{
				Name:        "redis_auth",
				EnvVars:     []string{"REDIS_AUTH"},
				Destination: &redisAuth,
				DefaultText: "",
				Value:       "",
				Usage:       "redis Auth, string",
			},
		},
		Before: func(ctx *cli.Context) error {
			log.SetFormatter(&log.JSONFormatter{})
			if debug {
				log.SetLevel(log.DebugLevel)
			}
			if redisAddr == "" {
				return fmt.Errorf("required REDIS_ADDR")
			}
			log.Debug("cli.main starting.")
			return nil
		},
		Action: func(ctx *cli.Context) error {
			cli := store.NewRedisCli(redisAddr, redisAuth)
			if err := cli.Ping().Err(); err != nil {
				return err
			}

			bat.Serve(addr, cli)
			return nil
		},
	}

	app.Run(os.Args)
}
