package tool

import (
	"fmt"
	"os"
	"strconv"
)

func StringFromEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func IntFromEnv(key string, fallback int) (int, error) {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
		}
		return n, nil
	}
	return fallback, nil
}
