package imgflip

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMemesResponseDecodesTemplates(t *testing.T) {
	const body = `{"success":true,"data":{"memes":[{"id":"181913649","name":"Drake Hotline Bling","box_count":2,"width":1200,"height":1200}]}}`

	var response GetMemesResponse
	require.NoError(t, json.NewDecoder(strings.NewReader(body)).Decode(&response))

	template := response.Data.Memes[0]
	assert.True(t, response.Success)
	assert.Equal(t, "181913649", template.ID)
	assert.Equal(t, "Drake Hotline Bling", template.Name)
	assert.Equal(t, 2, template.BoxCount)
	assert.Equal(t, 1200, template.Width)
	assert.Equal(t, 1200, template.Height)
}

func TestNewClientUsesDefaultTimeout(t *testing.T) {
	client := NewClient(nil)
	assert.Equal(t, defaultTimeout, client.httpClient.Timeout)
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
				assert.Equal(t, memesEndpoint, request.URL.String())
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
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
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
				assert.Equal(t, http.MethodPost, request.Method)
				assert.Equal(t, captionEndpoint, request.URL.String())
				assert.Equal(t, "application/x-www-form-urlencoded", request.Header.Get("Content-Type"))
				require.NoError(t, request.ParseForm())
				wantForm := map[string]string{
					"template_id":    "181913649",
					"username":       "meme-user",
					"password":       password,
					"boxes[0][text]": "first",
					"boxes[1][text]": "second",
				}
				for key, want := range wantForm {
					assert.Equal(t, want, request.Form.Get(key), "form %q", key)
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
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				assert.NotContains(t, err.Error(), password)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
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
