package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"ms-gofiber/internal/dto"
	"ms-gofiber/pkg/apperror"
)

type requestValidator interface {
	ValidateStruct(any) error
}

func defaultSkippedPaths() map[string]struct{} {
	return map[string]struct{}{
		"/v1/health":           {},
		"/v1/internal/echo":    {},
		"/v1/client/self-call": {},
	}
}

func HeaderGuard(validate requestValidator, skippedPaths map[string]struct{}) fiber.Handler {
	if skippedPaths == nil {
		skippedPaths = defaultSkippedPaths()
	}

	return func(c *fiber.Ctx) error {
		if shouldSkipPath(c.Path(), skippedPaths) {
			return c.Next()
		}

		var header dto.RequestHeader
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

func ExternalIDGuard(redisClient *redis.Client, ttl time.Duration, skippedPaths map[string]struct{}) fiber.Handler {
	if skippedPaths == nil {
		skippedPaths = defaultSkippedPaths()
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
		ok, err := redisClient.SetNX(c.UserContext(), key, true, ttl).Result()
		if err != nil {
			logMiddlewareError(c, err)
			return apperror.Wrap(apperror.ErrInternal, "failed to store X-EXTERNAL-ID", err)
		}
		if !ok {
			return apperror.New(apperror.ErrConflict, "duplicate X-EXTERNAL-ID")
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
