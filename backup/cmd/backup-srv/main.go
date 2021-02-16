package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/backup/borgbackup"
	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
	"github.com/mjm/pi-tools/backup/service/backupservice"
	"github.com/mjm/pi-tools/backup/tarsnap"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	tarsnapPath    = flag.String("tarsnap-path", "tarsnap", "Path to the Tarsnap binary")
	tarsnapKeyPath = flag.String("tarsnap-keyfile", "", "Path to the Tarsnap key for the backups")
	borgPath       = flag.String("borg-path", "borg", "Path to the Borg binary")
	borgRepoPath   = flag.String("borg-repo-path", "/backup/borg/backup", "Path to the Borg backup repository")
)

func main() {
	rpc.SetDefaultHTTPPort(2320)
	rpc.SetDefaultGRPCPort(2321)
	flag.Parse()

	stopObs := observability.MustStart("backup-srv")
	defer stopObs()

	b := borgbackup.New(*borgPath)
	t := tarsnap.New(*tarsnapPath)
	backupService := backupservice.New(b, t, backupservice.Config{
		BorgRepoPath:   *borgRepoPath,
		TarsnapKeyPath: *tarsnapKeyPath,
	})

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		backuppb.RegisterBackupServiceServer(server, backupService)
	}))

	signal.Wait()
}
