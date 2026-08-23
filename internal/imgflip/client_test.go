package imgflip

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetMemesResponseDecodesTemplates(t *testing.T) {
	const body = `{"success":true,"data":{"memes":[{"id":"181913649","name":"Drake Hotline Bling","box_count":2,"width":1200,"height":1200}]}}`

	var response GetMemesResponse
	if err := json.NewDecoder(strings.NewReader(body)).Decode(&response); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	template := response.Data.Memes[0]
	if !response.Success || template.ID != "181913649" || template.Name != "Drake Hotline Bling" || template.BoxCount != 2 || template.Width != 1200 || template.Height != 1200 {
		t.Fatalf("decoded response = %#v", response)
	}
}

func TestNewClientUsesDefaultTimeout(t *testing.T) {
	client := NewClient(nil)
	if client.httpClient.Timeout != defaultTimeout {
		t.Fatalf("default timeout = %s, want %s", client.httpClient.Timeout, defaultTimeout)
	}
}

func TestClientTemplates(t *testing.T) {
	tests := []struct {
		name      string
		transport http.RoundTripper
		want      []Template
		wantErr   string
	}{
		{
			name: "success",
			transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.String() != memesEndpoint {
					t.Errorf("request URL = %q, want %q", request.URL.String(), memesEndpoint)
				}
				return response(http.StatusOK, `{"success":true,"data":{"memes":[{"id":"1","name":"Template","box_count":2,"width":100,"height":200}]}}`), nil
			}),
			want: []Template{{ID: "1", Name: "Template", BoxCount: 2, Width: 100, Height: 200}},
		},
		{
			name: "transport failure",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network unavailable")
			}),
			wantErr: "retrieve templates",
		},
		{
			name: "unsuccessful response",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, `{"success":false}`), nil
			}),
			wantErr: "Imgflip reported an unsuccessful response",
		},
		{
			name: "malformed response",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, `{`), nil
			}),
			wantErr: "decode template response",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(&http.Client{Transport: test.transport})
			got, err := client.Templates(context.Background())
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("Templates() error = %v, want containing %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Templates() error = %v", err)
			}
			if len(got) != len(test.want) || got[0] != test.want[0] {
				t.Fatalf("Templates() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestClientCaptionImage(t *testing.T) {
	const password = "not-for-output"
	tests := []struct {
		name      string
		transport http.RoundTripper
		want      CaptionImageResult
		wantErr   string
	}{
		{
			name: "success",
			transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodPost {
					t.Errorf("request method = %q, want %q", request.Method, http.MethodPost)
				}
				if request.URL.String() != captionEndpoint {
					t.Errorf("request URL = %q, want %q", request.URL.String(), captionEndpoint)
				}
				if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
					t.Errorf("Content-Type = %q, want form encoding", request.Header.Get("Content-Type"))
				}
				if err := request.ParseForm(); err != nil {
					t.Fatalf("ParseForm() error = %v", err)
				}
				wantForm := map[string]string{
					"template_id":    "181913649",
					"username":       "meme-user",
					"password":       password,
					"boxes[0][text]": "first",
					"boxes[1][text]": "second",
				}
				for key, want := range wantForm {
					if got := request.Form.Get(key); got != want {
						t.Errorf("form %q = %q, want %q", key, got, want)
					}
				}
				return response(http.StatusOK, `{"success":true,"data":{"url":"https://i.imgflip.com/image.jpg","page_url":"https://imgflip.com/i/page"}}`), nil
			}),
			want: CaptionImageResult{ImageURL: "https://i.imgflip.com/image.jpg", PageURL: "https://imgflip.com/i/page"},
		},
		{
			name: "transport failure",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network unavailable")
			}),
			wantErr: "create captioned meme",
		},
		{
			name: "HTTP failure",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusBadGateway, ""), nil
			}),
			wantErr: "unexpected HTTP status",
		},
		{
			name: "unsuccessful response",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, `{"success":false,"error_message":"invalid template"}`), nil
			}),
			wantErr: "invalid template",
		},
		{
			name: "malformed response",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, "{"), nil
			}),
			wantErr: "decode caption response",
		},
		{
			name: "missing generated URLs",
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, `{"success":true,"data":{"url":"https://i.imgflip.com/image.jpg"}}`), nil
			}),
			wantErr: "missing generated URLs",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(&http.Client{Transport: test.transport})
			got, err := client.CaptionImage(context.Background(), CaptionImageRequest{
				TemplateID: "181913649",
				Username:   "meme-user",
				Password:   password,
				Texts:      []string{"first", "second"},
			})
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("CaptionImage() error = %v, want containing %q", err, test.wantErr)
				}
				if strings.Contains(err.Error(), password) {
					t.Fatalf("CaptionImage() error exposed password: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("CaptionImage() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("CaptionImage() = %#v, want %#v", got, test.want)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func response(status int, body string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
