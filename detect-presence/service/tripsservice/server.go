package tripsservice

import (
	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
)

type Server struct {
	db *otelsql.DB
	q  *database.Queries
}

func New(db *otelsql.DB) *Server {
	return &Server{
		db: db,
		q:  database.New(db),
	}
}
