package telegram

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type SetMyCommandsRequest struct {
	Commands []BotCommand `json:"commands"`
}

func (c *Client) SetMyCommands(ctx context.Context, req SetMyCommandsRequest) error {
	ctx, span := tracer.Start(ctx, "telegram.SetMyCommands",
		trace.WithAttributes(
			label.Int("telegram.request.param.commands.count", len(req.Commands))))
	defer span.End()

	var resp VoidResponse
	if err := c.perform(ctx, "setMyCommands", req, &resp); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}
