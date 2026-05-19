package tool

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

var randomInt = func(reader io.Reader, max *big.Int) (*big.Int, error) {
	return rand.Int(reader, max)
}

func GenerateRequestReference() string {
	maxNumber := big.NewInt(1_000_000_000_000)
	number, err := randomInt(rand.Reader, maxNumber)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%012d", number)
}

func GenerateRequestRefnum() string {
	return GenerateRequestReference()
}
