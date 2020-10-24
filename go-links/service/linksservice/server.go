package linksservice

import (
	"github.com/mjm/pi-tools/go-links/database"
)

type Server struct {
	db *database.Client
}

func New(db *database.Client) *Server {
	return &Server{
		db: db,
	}
}
