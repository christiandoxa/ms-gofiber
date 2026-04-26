package middleware

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"ms-gofiber/internal/app/adapter/dto"
	"ms-gofiber/pkg/apperror"
)

type RequestValidator interface {
	ValidateStruct(any) error
}

var reqHeaderParser = func(c *fiber.Ctx, out *dto.RequestHeader) error {
	return c.ReqHeaderParser(out)
}

func DefaultSkippedPaths() map[string]struct{} {
	return map[string]struct{}{
		"/v1/health":           {},
		"/v1/internal/echo":    {},
		"/v1/client/self-call": {},
	}
}

func HeaderGuard(validate RequestValidator, skippedPaths map[string]struct{}) fiber.Handler {
	if skippedPaths == nil {
		skippedPaths = DefaultSkippedPaths()
	}

	return func(c *fiber.Ctx) error {
		if shouldSkipPath(c.Path(), skippedPaths) {
			return c.Next()
		}

		var header dto.RequestHeader
		if err := reqHeaderParser(c, &header); err != nil {
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

func ExternalIDGuard(redisClient *redis.Client, ttl time.Duration, skippedPaths map[string]struct{}) fiber.Handler {
	if skippedPaths == nil {
		skippedPaths = DefaultSkippedPaths()
	}
	if ttl <= 0 {
		ttl = 60 * time.Second
	}

	return func(c *fiber.Ctx) error {
		if shouldSkipPath(c.Path(), skippedPaths) {
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
	if v := c.Locals("logger"); v != nil {
		if logger, ok := v.(*logrus.Entry); ok && logger != nil {
			logger.WithError(err).Error("middleware error")
		}
	}
}
