package messagesservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/metric/global"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/database"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/storage"
)

type Server struct {
	db     storage.DB
	q      *database.Queries
	t      *telegram.Client
	trips  tripspb.TripsServiceClient
	chatID int

	metrics metrics
}

func New(db storage.DB, t *telegram.Client, trips tripspb.TripsServiceClient, chatID int) *Server {
	meter := global.Meter(instrumentationName)
	return &Server{
		db:      db,
		q:       database.New(db),
		t:       t,
		trips:   trips,
		chatID:  chatID,
		metrics: newMetrics(meter),
	}
}

func (s *Server) buildTripMessage(ctx context.Context, tripID uuid.UUID) (string, *telegram.ReplyMarkup, error) {
	res, err := s.trips.GetTrip(ctx, &tripspb.GetTripRequest{
		Id: tripID.String(),
	})
	if err != nil {
		return "", nil, err
	}

	leftAt, err := time.Parse(time.RFC3339, res.GetTrip().GetLeftAt())
	if err != nil {
		return "", nil, status.Errorf(codes.InvalidArgument, "invalid left at timestamp: %s", err)
	}
	returnedAt, err := time.Parse(time.RFC3339, res.GetTrip().GetReturnedAt())
	if err != nil {
		return "", nil, status.Errorf(codes.InvalidArgument, "invalid returned at timestamp: %s", err)
	}
	duration := returnedAt.Sub(leftAt)

	// fetch the most popular tags for trips and offer them as inline-reply options
	tagsResp, err := s.trips.ListTags(ctx, &tripspb.ListTagsRequest{
		Limit: int32(len(res.GetTrip().GetTags())) + 3,
	})
	if err != nil {
		return "", nil, status.Errorf(codes.Internal, "fetching popular tags: %v", err)
	}

	tagSet := map[string]struct{}{}
	for _, tag := range res.GetTrip().GetTags() {
		tagSet[tag] = struct{}{}
	}

	var buttonRows [][]telegram.InlineKeyboardButton
	for _, tag := range tagsResp.GetTags() {
		if len(buttonRows) >= 3 {
			break
		}
		if _, ok := tagSet[tag.GetName()]; ok {
			continue
		}

		buttonRows = append(buttonRows, []telegram.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("🏷 %s", tag.GetName()),
				CallbackData: fmt.Sprintf("TAG_TRIP#%s", tag.GetName()),
			},
		})
	}
	buttonRows = append(buttonRows, []telegram.InlineKeyboardButton{
		{
			Text:         "🗑 Ignore",
			CallbackData: "IGNORE",
		},
	})
	replyMarkup := &telegram.ReplyMarkup{
		InlineKeyboard: buttonRows,
	}

	var text strings.Builder
	templateInput := &tripCompletedTemplateInput{
		ReturnedAt: returnedAt,
		Duration:   duration,
		Tags:       res.GetTrip().GetTags(),
	}
	if err := templates.ExecuteTemplate(&text, tripCompletedTemplate, templateInput); err != nil {
		return "", nil, status.Errorf(codes.Internal, "rendering message template: %s", err)
	}

	return text.String(), replyMarkup, nil
}
