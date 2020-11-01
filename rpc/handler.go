package rpc

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mjm/pi-tools/debug"
)

var (
	httpPort *int
	grpcPort *int
)

func SetDefaultHTTPPort(port int) {
	httpPort = flag.Int("http-port", port, "HTTP port to listen on for metrics and API requests")
}

func SetDefaultGRPCPort(port int) {
	grpcPort = flag.Int("grpc-port", port, "gRPC port to listen on for non-grpc-web API requests")
}

func ListenAndServe(opts ...Option) {
	if httpPort == nil {
		log.Panicf("no default HTTP port configured")
	}

	h, g := newHandler(opts...)
	addr := fmt.Sprintf(":%d", *httpPort)

	go func() {
		log.Printf("Listening on %s for HTTP", addr)
		if err := http.ListenAndServe(addr, h); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Panicf("error listening to HTTP: %v", err)
			}
		}
	}()

	if grpcPort != nil {
		go func() {
			grpcAddr := fmt.Sprintf(":%d", *grpcPort)
			lis, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				log.Panicf("error listening to gRPC: %v", err)
			}

			log.Printf("Listening on %s for gRPC", grpcAddr)
			if err := g.Serve(lis); err != nil {
				log.Panicf("error serving gRPC: %v", err)
			}
		}()
	}
}

func newHandler(opts ...Option) (http.Handler, *grpc.Server) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	reflection.Register(grpcServer)

	if cfg.registerFn != nil {
		cfg.registerFn(grpcServer)
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		if debug.IsEnabled() {
			return true
		}

		log.Printf("Rejecting unknown origin: %s. Requests in production should be proxied by Homebase.", origin)
		return false
	}))

	handler := otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpc.IsAcceptableGrpcCorsRequest(r) || wrappedGrpc.IsGrpcWebRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
				return
			}

			if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
				return
			}

			http.DefaultServeMux.ServeHTTP(w, r)
		}),
		"Server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
		}))

	return handler, grpcServer
}

type Option func(*config)

type config struct {
	registerFn func(*grpc.Server)
}

func WithRegisteredServices(f func(server *grpc.Server)) Option {
	return func(c *config) {
		c.registerFn = f
	}
}
