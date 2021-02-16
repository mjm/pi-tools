package backupservice

import (
	"github.com/mjm/pi-tools/backup/borgbackup"
	"github.com/mjm/pi-tools/backup/tarsnap"
)

type Config struct {
	BorgRepoPath   string
	TarsnapKeyPath string
}

type Server struct {
	borg    *borgbackup.Borg
	tarsnap *tarsnap.Tarsnap
	cfg     Config
}

func New(borg *borgbackup.Borg, ts *tarsnap.Tarsnap, cfg Config) *Server {
	return &Server{
		borg:    borg,
		tarsnap: ts,
		cfg:     cfg,
	}
}
