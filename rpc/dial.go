package rpc

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Dial creates a connection to a gRPC service in the cluster. The service is assumed to not
// use TLS. The client is automatically set up for distributed tracing.
func Dial(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts = append([]grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}, opts...)

	return grpc.DialContext(ctx, target, opts...)
}

// MustDial calls Dial, but panics if there is an error creating the connection. Use this in main functions,
// when panicking was going to be your solution to an error anyway.
func MustDial(ctx context.Context, target string, opts ...grpc.DialOption) *grpc.ClientConn {
	conn, err := Dial(ctx, target, opts...)
	if err != nil {
		log.Panicf("Error dialing gRPC service at %q: %v", target, err)
	}
	return conn
}
