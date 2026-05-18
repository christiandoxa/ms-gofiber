package handler

import (
	"context"

	"github.com/christiandoxa/welog"
	"github.com/christiandoxa/welog/pkg/model"
	"github.com/gofiber/fiber/v2"

	"ms-gofiber/pkg/httpx"
)

type fiberHTTPLogger struct {
	ctx *fiber.Ctx
}

func (l fiberHTTPLogger) Log(_ context.Context, req httpx.RequestLog, res httpx.ResponseLog) {
	if l.ctx == nil {
		return
	}
	welog.LogFiberClient(l.ctx, toTargetRequest(req), toTargetResponse(res))
}

func toTargetRequest(in httpx.RequestLog) model.TargetRequest {
	return model.TargetRequest{
		URL:         in.URL,
		Method:      in.Method,
		ContentType: in.ContentType,
		Header:      in.Header,
		Body:        in.Body,
		Timestamp:   in.Timestamp,
	}
}

func toTargetResponse(in httpx.ResponseLog) model.TargetResponse {
	return model.TargetResponse{
		Header:  in.Header,
		Body:    in.Body,
		Status:  in.Status,
		Latency: in.Latency,
	}
}
