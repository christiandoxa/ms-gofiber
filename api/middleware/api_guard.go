package middleware

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"ms-gofiber/pkg/apperror"
)

type RequestValidator interface {
	ValidateStruct(any) error
}

type RequestHeader struct {
	XPartnerID  string `reqHeader:"X-PARTNER-ID" json:"X-PARTNER-ID" validate:"required,alphanum,max=36"`
	ChannelID   string `reqHeader:"CHANNEL-ID" json:"CHANNEL-ID" validate:"required,alphanum,max=5"`
	XExternalID string `reqHeader:"X-EXTERNAL-ID" json:"X-EXTERNAL-ID" validate:"required,numeric,max=36"`
}

func skippedPaths() map[string]struct{} {
	return map[string]struct{}{
		"/v1/health":           {},
		"/v1/internal/echo":    {},
		"/v1/client/self-call": {},
	}
}

func HeaderGuard(validate RequestValidator) fiber.Handler {
	skipped := skippedPaths()
	return func(c *fiber.Ctx) error {
		if shouldSkipPath(c.Path(), skipped) {
			return c.Next()
		}

		var header RequestHeader
		if err := c.ReqHeaderParser(&header); err != nil {
			logMiddlewareError(c, err)
			return apperror.New(apperror.ErrBadRequest, "invalid request headers")
		}
		if err := validate.ValidateStruct(header); err != nil {
			logMiddlewareError(c, err)
			return err
		}

		c.Locals("request_header", header)
		return c.Next()
	}
}

func ExternalIDGuard(redisClient *redis.Client, ttl time.Duration) fiber.Handler {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}

	skipped := skippedPaths()
	return func(c *fiber.Ctx) error {
		if shouldSkipPath(c.Path(), skipped) {
			return c.Next()
		}

		externalID := c.Get("X-EXTERNAL-ID")
		if externalID == "" {
			return apperror.New(apperror.ErrBadRequest, "missing X-EXTERNAL-ID header")
		}
		if redisClient == nil {
			return c.Next()
		}

		key := "x-external-id:" + externalID
		err := redisClient.SetArgs(c.UserContext(), key, true, redis.SetArgs{Mode: "NX", TTL: ttl}).Err()
		if errors.Is(err, redis.Nil) {
			return apperror.New(apperror.ErrConflict, "duplicate X-EXTERNAL-ID")
		}
		if err != nil {
			logMiddlewareError(c, err)
			return apperror.Wrap(apperror.ErrInternal, "failed to store X-EXTERNAL-ID", err)
		}

		return c.Next()
	}
}

func shouldSkipPath(path string, skippedPaths map[string]struct{}) bool {
	_, ok := skippedPaths[path]
	return ok
}

func logMiddlewareError(c *fiber.Ctx, err error) {
	logFiberError(c, err, "middleware error")
}
