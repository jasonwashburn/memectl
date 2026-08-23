package imgflip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	memesEndpoint   = "https://api.imgflip.com/get_memes"
	captionEndpoint = "https://api.imgflip.com/caption_image"
	defaultTimeout  = 30 * time.Second
)

// Template is a public Imgflip meme template.
type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BoxCount int    `json:"box_count"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	URL      string `json:"url"`
}

// GetMemesResponse is the response returned by Imgflip's get_memes endpoint.
type GetMemesResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Memes []Template `json:"memes"`
	} `json:"data"`
}

// CaptionImageRequest contains the inputs required to caption an Imgflip template.
type CaptionImageRequest struct {
	TemplateID string
	Username   string
	Password   string
	Texts      []string
}

// CaptionImageResult identifies the hosted image and its Imgflip page.
type CaptionImageResult struct {
	ImageURL string
	PageURL  string
}

type captionImageResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
	Data         struct {
		URL     string `json:"url"`
		PageURL string `json:"page_url"`
	} `json:"data"`
}

// Client retrieves data from Imgflip.
type Client struct {
	httpClient      *http.Client
	memesEndpoint   string
	captionEndpoint string
}

// NewClient returns an Imgflip client using httpClient. A nil client uses a bounded default client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{httpClient: httpClient, memesEndpoint: memesEndpoint, captionEndpoint: captionEndpoint}
}

// Templates retrieves the public meme template list.
func (c *Client) Templates(ctx context.Context) ([]Template, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.memesEndpoint, nil)
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

// CaptionImage creates a captioned meme from an Imgflip template.
func (c *Client) CaptionImage(ctx context.Context, input CaptionImageRequest) (CaptionImageResult, error) {
	form := url.Values{
		"template_id": {input.TemplateID},
		"username":    {input.Username},
		"password":    {input.Password},
	}
	for index, text := range input.Texts {
		form.Set("boxes["+strconv.Itoa(index)+"][text]", text)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.captionEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return CaptionImageResult{}, fmt.Errorf("create caption request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return CaptionImageResult{}, fmt.Errorf("create captioned meme: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return CaptionImageResult{}, fmt.Errorf("create captioned meme: unexpected HTTP status %s", response.Status)
	}

	var result captionImageResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return CaptionImageResult{}, fmt.Errorf("decode caption response: %w", err)
	}
	if !result.Success {
		if result.ErrorMessage != "" {
			return CaptionImageResult{}, fmt.Errorf("create captioned meme: Imgflip reported an unsuccessful response: %s", result.ErrorMessage)
		}
		return CaptionImageResult{}, fmt.Errorf("create captioned meme: Imgflip reported an unsuccessful response")
	}
	if result.Data.URL == "" || result.Data.PageURL == "" {
		return CaptionImageResult{}, fmt.Errorf("decode caption response: missing generated URLs")
	}

	return CaptionImageResult{ImageURL: result.Data.URL, PageURL: result.Data.PageURL}, nil
}
