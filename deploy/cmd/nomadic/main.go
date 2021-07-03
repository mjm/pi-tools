package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/apps"
	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/pkg/nomadic"
	nomadicpb "github.com/mjm/pi-tools/pkg/nomadic/proto/nomadic"
	"github.com/mjm/pi-tools/pkg/nomadic/service/nomadicservice"
	"github.com/mjm/pi-tools/pkg/spanerr"
)

var tracer = otel.Tracer("github.com/mjm/pi-tools/deploy/cmd/nomadic")

func main() {
	var tp *sdktrace.TracerProvider

	app := &cli.App{
		Name: "nomadic",
		Authors: []*cli.Author{
			{
				Name:  "Matt Moriarity",
				Email: "matt@mattmoriarity.com",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "debug-tracing",
			},
			&cli.StringFlag{
				Name: "trace-id",
			},
			&cli.StringFlag{
				Name: "parent-span-id",
			},
		},
		Before: func(c *cli.Context) error {
			traceIDStr := c.String("trace-id")
			parentSpanIDStr := c.String("parent-span-id")
			if traceIDStr == "" || parentSpanIDStr == "" {
				return nil
			}

			traceID, err := trace.TraceIDFromHex(traceIDStr)
			if err != nil {
				return err
			}
			parentSpanID, err := trace.SpanIDFromHex(parentSpanIDStr)
			if err != nil {
				return err
			}

			sc := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     parentSpanID,
				TraceFlags: trace.FlagsSampled,
			})
			c.Context = trace.ContextWithRemoteSpanContext(c.Context, sc)

			if c.Bool("debug-tracing") {
				log.Printf("Setting up stdout exporter")
				exporter, err := stdouttrace.New()
				if err != nil {
					return fmt.Errorf("creating stdout exporter: %w", err)
				}

				tp = sdktrace.NewTracerProvider(
					sdktrace.WithBatcher(exporter))
				otel.SetTracerProvider(tp)
			} else {
				hostIP := os.Getenv("HOST_IP")
				exporter, err := otlptracegrpc.New(
					context.Background(),
					otlptracegrpc.WithInsecure(),
					otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:4317", hostIP)))
				if err != nil {
					return fmt.Errorf("creating otlp exporter: %w", err)
				}

				r, err := resource.New(context.Background(), resource.WithAttributes(
					semconv.ServiceNamespaceKey.String(os.Getenv("NOMAD_NAMESPACE")),
					semconv.ServiceNameKey.String("nomadic"),
					semconv.ServiceInstanceIDKey.String(os.Getenv("NOMAD_ALLOC_ID")),

					semconv.ContainerNameKey.String(os.Getenv("NOMAD_TASK_NAME")),

					semconv.HostNameKey.String(os.Getenv("HOSTNAME")),
					semconv.HostIDKey.String(os.Getenv("NOMAD_CLIENT_ID"))))
				if err != nil {
					return fmt.Errorf("creating telemetry resource: %w", err)
				}

				tp = sdktrace.NewTracerProvider(
					sdktrace.WithBatcher(exporter),
					sdktrace.WithResource(r))
				otel.SetTracerProvider(tp)
			}

			return nil
		},
		After: func(c *cli.Context) error {
			if tp != nil {
				log.Printf("Shutting down exporter")
				return tp.Shutdown(context.Background())
			}
			return nil
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

					//tracer := otel.Tracer("github.com/mjm/pi-tools/deploy/cmd/nomadic")
					ctx, span := tracer.Start(c.Context, "install",
						trace.WithAttributes(
							attribute.String("app.name", app.Name())))
					defer span.End()

					if err := app.Install(ctx, clients); err != nil {
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
					serverSocketPath := c.Path("server-socket-path")
					ctx, span := tracer.Start(c.Context, "perform-deploy",
						trace.WithAttributes(
							attribute.String("server.socket_path", serverSocketPath)))
					defer span.End()

					eventCh := make(chan *deploypb.ReportEvent, 10)
					ctx = nomadic.WithEvents(ctx, nomadic.NewChannelEventReporter(eventCh))
					doneCh := make(chan struct{})

					clients, err := nomadic.DefaultClients()
					if err != nil {
						return spanerr.RecordError(ctx, err)
					}

					conn, err := grpc.DialContext(ctx, "unix://"+serverSocketPath,
						grpc.WithInsecure(),
						grpc.WithBlock())
					if err != nil {
						return spanerr.RecordError(ctx, err)
					}

					client := nomadicpb.NewNomadicClient(conn)
					stream, err := client.StreamEvents(c.Context)
					if err != nil {
						return spanerr.RecordError(ctx, err)
					}

					var wg sync.WaitGroup
					var errored int32
					for appName, app := range nomadic.Registry() {
						wg.Add(1)

						go func(appName string, app nomadic.Deployable) {
							defer wg.Done()

							ctx, span := tracer.Start(ctx, "Deployable.Install",
								trace.WithAttributes(
									attribute.String("app.name", app.Name())))
							defer span.End()

							if err := app.Install(ctx, clients); err != nil {
								_ = spanerr.RecordError(ctx, err)
								nomadic.Events(ctx).Error("App %s failed to install", appName, nomadic.WithError(err))
								atomic.AddInt32(&errored, 1)
								return
							}
						}(appName, app)
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
						return spanerr.RecordError(ctx, err)
					}

					if errored != 0 {
						err := fmt.Errorf("one or more apps failed to install")
						return spanerr.RecordError(ctx, err)
					}
					return nil
				},
			},
			{
				Name: "test-server",
				Action: func(c *cli.Context) error {
					binaryPath := c.Args().First()

					doneCh := make(chan struct{})
					eventCh := make(chan *deploypb.ReportEvent)
					go func() {
						for evt := range eventCh {
							log.Println(evt)
						}
						close(doneCh)
					}()

					deployErr := nomadicservice.DeployAll(c.Context, binaryPath, eventCh)
					<-doneCh

					if deployErr == nil {
						log.Println("All apps finished deploying successfully")
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
