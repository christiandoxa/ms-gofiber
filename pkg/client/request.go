package client

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
)

func RequestJSON(ctx context.Context, request any, response any, url string, header map[string]string) (int, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	result, err := Do(ctx, Request{
		Method: http.MethodPost,
		URL:    url,
		Header: header,
		Body:   body,
	})
	if err != nil {
		return 0, err
	}

	if response != nil {
		if err := json.Unmarshal(result.Body, response); err != nil {
			return 0, err
		}
	}
	return result.StatusCode, nil
}
