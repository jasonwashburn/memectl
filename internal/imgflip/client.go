package imgflip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	memesEndpoint  = "https://api.imgflip.com/get_memes"
	defaultTimeout = 30 * time.Second
)

// Template is a public Imgflip meme template.
type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BoxCount int    `json:"box_count"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// GetMemesResponse is the response returned by Imgflip's get_memes endpoint.
type GetMemesResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Memes []Template `json:"memes"`
	} `json:"data"`
}

// Client retrieves data from Imgflip.
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// NewClient returns an Imgflip client using httpClient. A nil client uses a bounded default client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{httpClient: httpClient, endpoint: memesEndpoint}
}

// Templates retrieves the public meme template list.
func (c *Client) Templates(ctx context.Context) ([]Template, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create template request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("retrieve templates: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("retrieve templates: unexpected HTTP status %s", response.Status)
	}

	var result GetMemesResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode template response: %w", err)
	}
	if !result.Success {
		return nil, fmt.Errorf("retrieve templates: Imgflip reported an unsuccessful response")
	}
	if result.Data.Memes == nil {
		return nil, fmt.Errorf("decode template response: missing meme list")
	}

	return result.Data.Memes, nil
}
