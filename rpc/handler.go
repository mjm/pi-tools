package rpc

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mjm/pi-tools/debug"
)

var httpPort *int

func SetDefaultHTTPPort(port int) {
	httpPort = flag.Int("http-port", port, "HTTP port to listen on for metrics and API requests")
}

func ListenAndServe(opts ...Option) {
	if httpPort == nil {
		log.Panicf("no default HTTP port configured")
	}

	h := NewHandler(opts...)
	addr := fmt.Sprintf(":%d", *httpPort)

	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("error listening to HTTP: %v", err)
		}
	}
}

func NewHandler(opts ...Option) http.Handler {
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

	return handler
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
