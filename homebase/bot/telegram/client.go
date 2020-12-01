package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

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
	metrics metrics
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

	meter := otel.Meter(instrumentationName)
	return &Client{
		cfg:     cfg,
		baseURL: baseURL,
		metrics: newMetrics(meter),
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
		c.measureRequest(ctx, action, err, 0, 0)
		return spanerr.RecordError(ctx, fmt.Errorf("creating request: %w", err))
	}

	r.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	res, err := c.cfg.HTTPClient.Do(r)
	duration := time.Now().Sub(startTime)
	if err != nil {
		c.measureRequest(ctx, action, err, 0, duration)
		return spanerr.RecordError(ctx, fmt.Errorf("performing request: %w", err))
	}
	defer res.Body.Close()

	span.SetAttributes(label.Int("telegram.response.status", res.StatusCode))

	if res.StatusCode != http.StatusOK {
		var resErr Error
		if err := json.NewDecoder(res.Body).Decode(&resErr); err != nil {
			c.measureRequest(ctx, action, err, res.StatusCode, duration)
			return spanerr.RecordError(ctx, fmt.Errorf("decoding error response: %w", err))
		}

		c.measureRequest(ctx, action, err, res.StatusCode, duration)
		return spanerr.RecordError(ctx, resErr)
	}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		c.measureRequest(ctx, action, err, res.StatusCode, duration)
		return spanerr.RecordError(ctx, fmt.Errorf("decoding response: %w", err))
	}

	c.measureRequest(ctx, action, nil, res.StatusCode, duration)
	return nil
}

func (c *Client) measureRequest(ctx context.Context, action string, err error, status int, duration time.Duration) {
	if err != nil || status > 299 {
		c.metrics.RequestErrorsTotal.Add(ctx, 1,
			label.String("action", action),
			semconv.HTTPStatusCodeKey.Int(status))
	}
	c.metrics.RequestTotal.Add(ctx, 1,
		label.String("action", action),
		semconv.HTTPStatusCodeKey.Int(status))
	if duration != 0 {
		c.metrics.RequestDurationSeconds.Record(ctx, duration.Seconds(),
			label.String("action", action))
	}
}

type Error struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e Error) Error() string {
	return e.Description
}

type VoidResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}
