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
					t.Errorf("request URL = %q, want %q", request.URL, memesEndpoint)
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
