package tripsservice

import (
	"github.com/mjm/pi-tools/detect-presence/database"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/storage"
)

type Server struct {
	db       storage.DB
	q        *database.Queries
	messages messagespb.MessagesServiceClient
}

func New(db storage.DB, messages messagespb.MessagesServiceClient) *Server {
	return &Server{
		db:       db,
		q:        database.New(db),
		messages: messages,
	}
}
