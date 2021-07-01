package nomadicservice

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/pkg/nomadic/proto/nomadic"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func DeployAll(ctx context.Context, binaryPath string, ch chan<- *deploypb.ReportEvent) error {
	sockDir, err := os.MkdirTemp("", "deploy-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(sockDir)
	sockPath := filepath.Join(sockDir, "deploy.sock")

	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		return err
	}
	defer lis.Close()

	doneCh := make(chan struct{})

	s := &server{ch: ch, doneCh: doneCh}
	srv := grpc.NewServer()
	nomadic.RegisterNomadicServer(srv, s)
	go srv.Serve(lis)

	sc := trace.SpanContextFromContext(ctx)
	traceID := sc.TraceID().String()
	spanID := sc.SpanID().String()

	cmd := exec.CommandContext(ctx, binaryPath, "--trace-id", traceID, "--parent-span-id", spanID, "perform-deploy", "--server-socket-path", sockPath)
	err = cmd.Run()
	<-doneCh
	srv.Stop()

	return err
}

type server struct {
	ch     chan<- *deploypb.ReportEvent
	doneCh chan struct{}
}

func (s *server) StreamEvents(server nomadic.Nomadic_StreamEventsServer) error {
	for {
		msg, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				close(s.ch)
				close(s.doneCh)
				return server.SendAndClose(&nomadic.StreamEventsResponse{})
			}
			return err
		}

		s.ch <- msg.GetEvent()
	}
}
