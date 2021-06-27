package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/apps"
	"github.com/mjm/pi-tools/pkg/nomadic"
	nomadicpb "github.com/mjm/pi-tools/pkg/nomadic/proto/nomadic"
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
			{
				Name: "perform-deploy",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "server-socket-path",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					eventCh := make(chan *deploypb.ReportEvent, 10)
					ctx := nomadic.WithEvents(c.Context, nomadic.NewChannelEventReporter(eventCh))
					doneCh := make(chan struct{})

					clients, err := nomadic.DefaultClients()
					if err != nil {
						return err
					}

					serverSocketPath := c.Path("server-socket-path")
					conn, err := grpc.DialContext(ctx, "unix://"+serverSocketPath,
						grpc.WithInsecure(),
						grpc.WithBlock())
					if err != nil {
						return err
					}

					client := nomadicpb.NewNomadicClient(conn)
					stream, err := client.StreamEvents(c.Context)
					if err != nil {
						return err
					}

					var wg sync.WaitGroup
					var errored int32
					for appName, app := range nomadic.Registry() {
						wg.Add(1)

						go func(app nomadic.Deployable) {
							defer wg.Done()

							if err := app.Install(ctx, clients); err != nil {
								nomadic.Events(ctx).Error("App %s failed to install", appName, nomadic.WithError(err))
								atomic.AddInt32(&errored, 1)
								return
							}
						}(app)
					}

					go func() {
						for evt := range eventCh {
							stream.Send(&nomadicpb.StreamEventsRequest{
								Event: evt,
							})
						}
						close(doneCh)
					}()

					wg.Wait()
					close(eventCh)
					<-doneCh

					if _, err := stream.CloseAndRecv(); err != nil {
						return err
					}

					if errored != 0 {
						return fmt.Errorf("one or more apps failed to install")
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
