package messagesservice

import (
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
)

type Server struct {
	t      *telegram.Client
	trips  tripspb.TripsServiceClient
	chatID int
}

func New(t *telegram.Client, trips tripspb.TripsServiceClient, chatID int) *Server {
	return &Server{
		t:      t,
		trips:  trips,
		chatID: chatID,
	}
}
