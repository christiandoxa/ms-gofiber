package rcsuccess

import (
	"net/http"
	"testing"
)

func TestGeneralSuccess(t *testing.T) {
	if GeneralSuccess.StatusCode != http.StatusOK || GeneralSuccess.ResponseCode != "2000000" {
		t.Fatalf("unexpected success response: %+v", GeneralSuccess)
	}
}
