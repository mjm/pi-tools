package tripsservice

import (
	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/storage"
)

type Server struct {
	db storage.DB
	q  *database.Queries
}

func New(db storage.DB) *Server {
	return &Server{
		db: db,
		q:  database.New(db),
	}
}
