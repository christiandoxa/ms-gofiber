package httpx

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"go.elastic.co/apm/module/apmhttp/v2"

	"ms-gofiber/pkg/logging"
)

const defaultTimeout = 10 * time.Second

type Request struct {
	Method      string
	URL         string
	ContentType string
	Header      map[string]string
	Body        []byte
	Timeout     time.Duration
}

type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

type Logger interface {
	Log(ctx context.Context, req RequestLog, res ResponseLog)
}

type RequestLog struct {
	URL         string
	Method      string
	ContentType string
	Header      map[string]interface{}
	Body        []byte
	Timestamp   time.Time
}

type ResponseLog struct {
	Header  map[string]interface{}
	Body    []byte
	Status  int
	Latency time.Duration
}

func Do(ctx context.Context, req Request, logger Logger) (*Response, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	}
	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}
	httpReq = httpReq.WithContext(ctx)

	reqLog := RequestLog{
		URL:         req.URL,
		Method:      req.Method,
		ContentType: httpReq.Header.Get("Content-Type"),
		Header:      headerToInterface(httpReq.Header),
		Body:        req.Body,
		Timestamp:   time.Now(),
	}

	start := time.Now()
	httpClient := apmhttp.WrapClient(&http.Client{Timeout: timeout})
	httpRes, err := httpClient.Do(httpReq)

	if err != nil {
		resLog := ResponseLog{
			Header:  map[string]interface{}{},
			Body:    nil,
			Status:  0,
			Latency: time.Since(start),
		}
		logRequest(ctx, logger, reqLog, resLog)
		return nil, err
	}
	defer closeBody(ctx, httpRes.Body)

	body, err := io.ReadAll(httpRes.Body)

	resLog := ResponseLog{
		Header:  headerToInterface(httpRes.Header),
		Body:    body,
		Status:  httpRes.StatusCode,
		Latency: time.Since(start),
	}
	logRequest(ctx, logger, reqLog, resLog)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: httpRes.StatusCode,
		Body:       body,
		Header:     httpRes.Header,
	}, nil
}

func headerToInterface(h http.Header) map[string]interface{} {
	m := make(map[string]interface{}, len(h))
	for k, v := range h {
		if len(v) == 1 {
			m[k] = v[0]
		} else {
			m[k] = v
		}
	}
	return m
}

func logRequest(ctx context.Context, logger Logger, req RequestLog, res ResponseLog) {
	if logger == nil {
		return
	}
	logger.Log(ctx, req, res)
}

func closeBody(ctx context.Context, body io.Closer) {
	if err := body.Close(); err != nil {
		logging.Error(ctx, err, "failed to close response body", nil)
	}
}
