package messagesservice

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/homebase/bot/telegram"
)

var commands = []telegram.BotCommand{
	{
		Command:     "tag",
		Description: "Add one or more tags to the most recent trip.",
	},
	{
		Command:     "untag",
		Description: "Remove one or more tags from the most recent trip.",
	},
	{
		Command:     "ignore",
		Description: "Ignore the most recent trip, hiding it from the list of trips in Homebase.",
	},
}

func (s *Server) RegisterCommands(ctx context.Context) error {
	if err := s.t.SetMyCommands(ctx, telegram.SetMyCommandsRequest{
		Commands: commands,
	}); err != nil {
		return status.Errorf(codes.Internal, "setting commands: %s", err)
	}

	return nil
}
