package linksservice

import (
	"github.com/mjm/pi-tools/go-links/database"
)

type Server struct {
	db *database.Queries
}

func New(db *database.Queries) *Server {
	return &Server{
		db: db,
	}
}
