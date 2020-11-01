package telegram

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type GetUpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type GetUpdatesResponse struct {
	OK          bool     `json:"ok"`
	Description string   `json:"description"`
	Result      []Update `json:"result"`
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type UpdateOrError struct {
	Update *Update
	Err    error
}

func (c *Client) WatchUpdates(ctx context.Context, ch chan<- UpdateOrError, req GetUpdatesRequest) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := c.GetUpdates(ctx, req)
				if err != nil {
					ch <- UpdateOrError{Err: err}
				} else {
					for _, update := range updates {
						ch <- UpdateOrError{Update: &update}

						// Update the offset in the request so that it fetches newer updates next time
						if update.UpdateID >= req.Offset {
							req.Offset = update.UpdateID + 1
						}
					}
				}
			}
		}
	}()
}

func (c *Client) GetUpdates(ctx context.Context, req GetUpdatesRequest) ([]Update, error) {
	ctx, span := tracer.Start(ctx, "telegram.GetUpdates",
		trace.WithAttributes(
			label.Int("telegram.request.param.offset", req.Offset),
			label.Int("telegram.request.param.limit", req.Limit),
			label.Int("telegram.request.param.timeout", req.Timeout),
			label.Array("telegram.request.param.allowed_updates", req.AllowedUpdates)))
	defer span.End()

	var resp GetUpdatesResponse
	if err := c.perform(ctx, "getUpdates", req, &resp); err != nil {
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Int("telegram.response.update_count", len(resp.Result)))
	return resp.Result, nil
}
