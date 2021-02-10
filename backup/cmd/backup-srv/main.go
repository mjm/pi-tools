package main

import (
	"flag"

	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/backup/borgbackup"
	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
	"github.com/mjm/pi-tools/backup/service/backupservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	//tarsnapKeyPath = flag.String("tarsnap-keyfile", "", "Path to the Tarsnap key for the backups")
	borgPath     = flag.String("borg-path", "borg", "Path to the Borg binary")
	borgRepoPath = flag.String("borg-repo-path", "/backup/borg", "Path to the Borg backup repository")
)

func main() {
	rpc.SetDefaultHTTPPort(2320)
	rpc.SetDefaultGRPCPort(2321)
	flag.Parse()

	stopObs := observability.MustStart("backup-srv")
	defer stopObs()

	b := borgbackup.New(*borgPath)
	backupService := backupservice.New(b, backupservice.Config{
		BorgRepoPath: *borgRepoPath,
	})

	rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		backuppb.RegisterBackupServiceServer(server, backupService)
	}))

	signal.Wait()
}
