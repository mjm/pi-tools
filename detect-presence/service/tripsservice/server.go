package tripsservice

import (
	"github.com/mjm/pi-tools/detect-presence/database"
)

type Server struct {
	db *database.Client
}

func New(db *database.Client) *Server {
	return &Server{
		db: db,
	}
}
