package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Scraper struct {
	BaseURL  string
	Password string
	HTTP     *http.Client
}

func New(baseURL string, password string) *Scraper {
	return &Scraper{
		BaseURL:  baseURL,
		Password: password,
		HTTP:     http.DefaultClient,
	}
}

func (s *Scraper) GetStats(ctx context.Context) (*Stats, error) {
	var stats Stats
	if err := s.get(ctx, "summaryRaw", &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *Scraper) get(ctx context.Context, method string, result interface{}) error {
	auth := ""
	if s.Password != "" {
		auth = "&auth=" + s.Password
	}
	u := fmt.Sprintf("%s/admin/api.php?%s%s", s.BaseURL, method, auth)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	res, err := s.HTTP.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("unexpected HTTP response %s", res.Status)
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}
