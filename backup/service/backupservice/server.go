package backupservice

import (
	"github.com/mjm/pi-tools/backup/borgbackup"
)

type Config struct {
	BorgRepoPath string
}

type Server struct {
	borg *borgbackup.Borg
	cfg  Config
}

func New(borg *borgbackup.Borg, cfg Config) *Server {
	return &Server{
		borg: borg,
		cfg:  cfg,
	}
}
