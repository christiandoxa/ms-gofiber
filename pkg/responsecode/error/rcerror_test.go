package rcerror

import (
	"errors"
	"net/http"
	"testing"
)

func TestResponseCodeErrors(t *testing.T) {
	if ErrInvalidFieldFormat.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected invalid field status")
	}
	if ErrInvalidMandatoryField.ResponseCode == "" {
		t.Fatalf("expected mandatory field response code")
	}
	if ErrDuplicateExternalID.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected duplicate status")
	}
	if ErrGeneral.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected general status")
	}
	if ErrTimeout.StatusCode != http.StatusRequestTimeout {
		t.Fatalf("unexpected timeout status")
	}
	if !errors.Is(ErrDuplicateCacheKey, ErrDuplicateCacheKey) || !errors.Is(ErrFailedToStoreData, ErrFailedToStoreData) {
		t.Fatalf("expected sentinel errors")
	}
}
