package tool

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadDotenvIfExists(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", filename, err)
	}
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("load %s: %w", filename, err)
	}
	return nil
}
