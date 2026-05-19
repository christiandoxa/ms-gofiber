package tool

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func HeaderToMap(header any, keys ...string) map[string]interface{} {
	headers := map[string]interface{}{}

	switch h := header.(type) {
	case *fiber.Ctx:
		for _, key := range keys {
			if value := h.Get(key); value != "" {
				headers[key] = value
			}
		}
	case *fasthttp.RequestHeader:
		for key, value := range h.All() {
			headers[string(key)] = string(value)
		}
	case *fasthttp.ResponseHeader:
		for key, value := range h.All() {
			headers[string(key)] = string(value)
		}
	}

	return headers
}
