package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

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
			{
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Action: func(c *cli.Context) error {
					if c.NArg() > 0 {
						cli.ShowCommandHelpAndExit(c, "list", 1)
						return nil
					}

					var appNames []string
					for appName := range nomadic.Registry() {
						appNames = append(appNames, appName)
					}
					sort.Strings(appNames)

					for _, appName := range appNames {
						fmt.Println(appName)
					}
					return nil
				},
			},
			{
				Name: "images",
				Action: func(c *cli.Context) error {
					if c.NArg() > 0 {
						cli.ShowCommandHelpAndExit(c, "images", 1)
						return nil
					}

					for _, imageURI := range nomadic.RegisteredImageURIs() {
						fmt.Println(imageURI)
					}
					return nil
				},
			},
		},
	}

	apps.Load()

	ctx := nomadic.WithEvents(context.Background(), nomadic.NewLoggingEventReporter())

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
