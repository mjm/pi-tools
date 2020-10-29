package linksservice

import (
	"github.com/mjm/pi-tools/go-links/database"
	"github.com/mjm/pi-tools/storage"
)

type Server struct {
	db *database.Queries
}

func New(db storage.DB) *Server {
	return &Server{
		db: database.New(db),
	}
}
