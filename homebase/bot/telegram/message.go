package telegram

import (
	"context"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from"`
	Date      int    `json:"date"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type SendMessageRequest struct {
	ChatID                int          `json:"chat_id"`
	Text                  string       `json:"text"`
	ParseMode             ParseMode    `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview"`
	DisableNotification   bool         `json:"disable_notification"`
	ReplyToMessageID      int          `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           *ReplyMarkup `json:"reply_markup,omitempty"`
}

type EditMessageTextRequest struct {
	ChatID      int          `json:"chat_id"`
	MessageID   int          `json:"message_id"`
	Text        string       `json:"text"`
	ParseMode   ParseMode    `json:"parse_mode,omitempty"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
}

type DeleteMessageRequest struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

type ParseMode string

const (
	MarkdownV2Mode     ParseMode = "MarkdownV2"
	HTMLMode           ParseMode = "HTML"
	MarkdownLegacyMode ParseMode = "Markdown"
)

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type SendMessageResponse struct {
	OK          bool    `json:"ok"`
	Description string  `json:"description"`
	Result      Message `json:"result"`
}

func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*Message, error) {
	ctx, span := tracer.Start(ctx, "telegram.SendMessage",
		trace.WithAttributes(
			label.Int("telegram.request.param.chat_id", req.ChatID),
			label.Int("telegram.request.param.text.length", len(req.Text)),
			label.String("telegram.request.param.parse_mode", string(req.ParseMode)),
			label.Bool("telegram.request.param.disable_web_page_preview", req.DisableWebPagePreview),
			label.Bool("telegram.request.param.disable_notification", req.DisableNotification),
			label.Int("telegram.request.param.reply_to_message_id", req.ReplyToMessageID)))
	defer span.End()

	var resp SendMessageResponse
	if err := c.perform(ctx, "sendMessage", req, &resp); err != nil {
		return nil, spanerr.RecordError(ctx, err)
	}

	msg := resp.Result
	span.SetAttributes(
		label.Int("telegram.response.message_id", msg.MessageID))
	return &msg, nil
}

func (c *Client) EditMessageText(ctx context.Context, req EditMessageTextRequest) (*Message, error) {
	ctx, span := tracer.Start(ctx, "telegram.EditMessageText",
		trace.WithAttributes(
			label.Int("telegram.request.param.chat_id", req.ChatID),
			label.Int("telegram.request.param.message_id", req.MessageID),
			label.Int("telegram.request.param.text.length", len(req.Text)),
			label.String("telegram.request.param.parse_mode", string(req.ParseMode))))
	defer span.End()

	var resp SendMessageResponse
	if err := c.perform(ctx, "editMessageText", req, &resp); err != nil {
		return nil, spanerr.RecordError(ctx, err)
	}

	msg := resp.Result
	span.SetAttributes(
		label.Int("telegram.response.message_id", msg.MessageID))
	return &msg, nil
}

func (c *Client) DeleteMessage(ctx context.Context, req DeleteMessageRequest) error {
	ctx, span := tracer.Start(ctx, "telegram.DeleteMessage",
		trace.WithAttributes(
			label.Int("telegram.request.param.chat_id", req.ChatID),
			label.Int("telegram.request.param.message_id", req.MessageID)))
	defer span.End()

	var resp VoidResponse
	if err := c.perform(ctx, "deleteMessage", req, &resp); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}
