package nomadic

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/pkg/nomadic/proto/nomadic"
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

	err = exec.CommandContext(ctx, binaryPath, "--server-socket-path", sockPath).Run()
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
