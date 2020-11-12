package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

func main() {
	app := &cli.App{
		Name:    "telepush",
		Usage:   "Telegram Bot push tool",
		Version: fmt.Sprintf("%s %s %s", version, commitID, buildDate),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server listen address",
				EnvVars: []string{"TELEPUSH_ADDRESS"},
				Value:   "0.0.0.0:8080",
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Server push token",
				EnvVars:  []string{"TELEPUSH_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "bot-api",
				Usage:   "Telegram api address",
				EnvVars: []string{"TELEPUSH_BOT_API"},
				Value:   "https://api.telegram.org",
			},
			&cli.StringFlag{
				Name:     "bot-token",
				Usage:    "Telegram api token",
				EnvVars:  []string{"TELEPUSH_BOT_TOKEN"},
				Required: true,
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "mritd",
				Email: "mritd@linux.com",
			},
		},
		Action: func(c *cli.Context) error {
			conf = Config{
				Addr:        c.String("addr"),
				Token:       c.String("token"),
				BotApiAddr:  c.String("bot-api"),
				BotApiToken: c.String("bot-token"),
			}
			var err error
			bot, err = NewTelegram(conf.BotApiAddr, conf.BotApiToken)
			if err != nil {
				return err
			}
			logger.Info("Telegram Bot init success...")
			Serve()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
