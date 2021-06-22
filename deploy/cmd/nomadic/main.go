package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mjm/pi-tools/apps"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

func main() {
	app := &cli.App{
		Name: "nomadic",
		Authors: []*cli.Author{
			{
				Name:  "Matt Moriarity",
				Email: "matt@mattmoriarity.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						cli.ShowCommandHelpAndExit(c, "install", 1)
						return nil
					}

					clients, err := nomadic.DefaultClients()
					if err != nil {
						return err
					}

					// TODO support multiple apps
					appName := c.Args().First()
					app := nomadic.Find(appName)

					if app == nil {
						return fmt.Errorf("Unknown application name %q", appName)
					}

					log.Printf("Installing %s", appName)
					if err := app.Install(c.Context, clients); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	apps.Load()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
