package telegram

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, req AnswerCallbackQueryRequest) error {
	ctx, span := tracer.Start(ctx, "telegram.AnswerCallbackQuery",
		trace.WithAttributes(
			attribute.String("telegram.request.param.callback_query_id", req.CallbackQueryID),
			attribute.Int("telegram.request.param.text.length", len(req.Text)),
			attribute.Bool("telegram.request.param.show_alert", req.ShowAlert),
			attribute.String("telegram.request.param.url", req.URL),
			attribute.Int("telegram.request.param.cache_time", req.CacheTime)))
	defer span.End()

	var resp VoidResponse
	if err := c.perform(ctx, "answerCallbackQuery", req, &resp); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}
