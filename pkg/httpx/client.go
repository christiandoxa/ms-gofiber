package httpx

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/christiandoxa/welog"
	"github.com/christiandoxa/welog/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp/v2"
)

func Do(
	c *fiber.Ctx,
	method, url, contentType string,
	hdr map[string]string,
	body []byte,
	timeout time.Duration,
) (status int, respBody []byte, respHeader http.Header, err error) {

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range hdr {
		req.Header.Set(k, v)
	}

	// propagate APM context dari Fiber
	req = req.WithContext(c.UserContext())

	reqModel := model.TargetRequest{
		URL:         url,
		Method:      method,
		ContentType: req.Header.Get("Content-Type"),
		Header:      headerToInterface(req.Header),
		Body:        body,
		Timestamp:   time.Now(),
	}

	start := time.Now()

	// http client dengan APM instrumentation
	httpClient := apmhttp.WrapClient(&http.Client{Timeout: timeout})
	res, err := httpClient.Do(req)

	var resModel model.TargetResponse
	if err != nil {
		resModel = model.TargetResponse{
			Header:  map[string]interface{}{},
			Body:    nil,
			Status:  0,
			Latency: time.Since(start),
		}
		welog.LogFiberClient(c, reqModel, resModel)
		return 0, nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Locals("logger").(*logrus.Entry).Error("Something went wrong", err.Error())
		}
	}(res.Body)

	rb, _ := io.ReadAll(res.Body)
	status = res.StatusCode
	respHeader = res.Header
	respBody = rb

	resModel = model.TargetResponse{
		Header:  headerToInterface(res.Header),
		Body:    rb,
		Status:  res.StatusCode,
		Latency: time.Since(start),
	}
	welog.LogFiberClient(c, reqModel, resModel)
	return status, respBody, respHeader, nil
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
