package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/v2"

	"ms-gofiber/pkg/apperror"
	"ms-gofiber/pkg/httpx"
	"ms-gofiber/pkg/respond"
)

type InternalHandler struct{}

func NewInternalHandler() *InternalHandler { return &InternalHandler{} }

// Echo GET /v1/internal/echo
func (h *InternalHandler) Echo(c *fiber.Ctx) error {
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
func (h *InternalHandler) SelfCall(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span, ctx := apm.StartSpan(ctx, "InternalHandler.SelfCall", "handler")
	defer span.End()

	target := fmt.Sprintf("%s/v1/internal/echo?msg=%s", c.BaseURL(), url.QueryEscape("hello-from-client"))
	status, body, _, err := httpx.Do(
		c,
		"GET",
		target,
		"application/json",
		map[string]string{"X-Demo": "self-call"},
		nil,
		5*time.Second,
	)
	if err != nil {
		c.Locals("logger").(*logrus.Entry).Error("Something went wrong")
		apm.CaptureError(ctx, err).Send()
		return apperror.New(apperror.ErrInternal, "self call failed")
	}

	var upstream any
	_ = json.Unmarshal(body, &upstream)

	return c.JSON(respond.SuccessEnvelope(fiber.Map{
		"upstream_status": status,
		"upstream_body":   upstream,
	}, nil))
}
