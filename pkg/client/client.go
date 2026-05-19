package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/christiandoxa/welog/pkg/infrastructure/logger"
)

const defaultTimeout = 5 * time.Second

type Request struct {
	Method  string
	URL     string
	Header  map[string]string
	Body    []byte
	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func Do(ctx context.Context, request Request) (*Response, error) {
	timeout := request.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	method := request.Method
	if method == "" {
		method = http.MethodGet
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, request.URL, bytes.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	for key, value := range request.Header {
		httpRequest.Header.Set(key, value)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Logger().WithError(err).Warn("failed to close response body")
		}
	}(httpResponse.Body) //nolint:errcheck

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: httpResponse.StatusCode,
		Header:     httpResponse.Header,
		Body:       body,
	}, nil
}
