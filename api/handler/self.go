package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"

	"ms-gofiber/api/respond"
	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/httpx"
)

type Internal struct{}

// Echo GET /v1/internal/echo
func (h *Internal) Echo(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, _ := apm.StartSpan(ctx, "InternalHandler.Echo", "handler")
	defer span.End()

	msg := c.Query("msg", "pong")
	return c.JSON(fiber.Map{
		"echo": msg,
		"ts":   time.Now().UTC(),
	})
}

// SelfCall GET /v1/client/self-call
func (h *Internal) SelfCall(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "InternalHandler.SelfCall", "handler")
	defer span.End()

	target := fmt.Sprintf("%s/v1/internal/echo?msg=%s", c.BaseURL(), url.QueryEscape("hello-from-client"))
	res, err := httpx.Do(ctx, httpx.Request{
		Method:      "GET",
		URL:         target,
		ContentType: "application/json",
		Header:      map[string]string{"X-Demo": "self-call"},
		Timeout:     5 * time.Second,
	}, fiberHTTPLogger{ctx: c})
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return apperror.New(apperror.ErrInternal, "self call failed")
	}

	var upstream any
	if err := json.Unmarshal(res.Body, &upstream); err != nil {
		apm.CaptureError(ctx, err).Send()
		return apperror.New(apperror.ErrInternal, "invalid self call response")
	}

	return c.JSON(respond.SuccessEnvelope(fiber.Map{
		"upstream_status": res.StatusCode,
		"upstream_body":   upstream,
	}, nil))
}
