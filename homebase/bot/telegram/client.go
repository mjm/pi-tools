package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type Config struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

type Client struct {
	cfg     Config
	baseURL *url.URL
}

func New(cfg Config) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("no telegram token provided")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.telegram.org"
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &Client{
		cfg:     cfg,
		baseURL: baseURL,
	}, nil
}

func (c *Client) requestURL(action string) *url.URL {
	u := *c.baseURL
	u.Path = fmt.Sprintf("bot%s/%s", c.cfg.Token, action)
	return &u
}

func (c *Client) perform(ctx context.Context, action string, params interface{}, resp interface{}) error {
	ctx, span := tracer.Start(ctx, "telegram.perform",
		trace.WithAttributes(
			label.String("telegram.action", action),
			label.String("telegram.request.type", fmt.Sprintf("%T", params)),
			label.String("telegram.response.type", fmt.Sprintf("%T", resp))))
	defer span.End()

	body, err := json.Marshal(params)
	if err != nil {
		return spanerr.RecordError(ctx, fmt.Errorf("serializing request params: %w", err))
	}

	span.SetAttributes(label.Int("telegram.request.length", len(body)))

	u := c.requestURL(action)
	span.SetAttributes(label.String("telegram.url", u.String()))

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return spanerr.RecordError(ctx, fmt.Errorf("creating request: %w", err))
	}

	r.Header.Set("Content-Type", "application/json")

	res, err := c.cfg.HTTPClient.Do(r)
	if err != nil {
		return spanerr.RecordError(ctx, fmt.Errorf("performing request: %w", err))
	}
	defer res.Body.Close()

	span.SetAttributes(label.Int("telegram.response.status", res.StatusCode))

	if res.StatusCode != http.StatusOK {
		var resErr Error
		if err := json.NewDecoder(res.Body).Decode(&resErr); err != nil {
			return spanerr.RecordError(ctx, fmt.Errorf("decoding error response: %w", err))
		}

		return spanerr.RecordError(ctx, resErr)
	}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return spanerr.RecordError(ctx, fmt.Errorf("decoding response: %w", err))
	}

	return nil
}

type Error struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e Error) Error() string {
	return e.Description
}
