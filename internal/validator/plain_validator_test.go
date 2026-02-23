package validator

import "testing"

func TestValidatePagination(t *testing.T) {
	limit, offset, err := ValidatePagination("", "", 100)
	if err != nil || limit != 10 || offset != 0 {
		t.Fatalf("unexpected default pagination: %d %d %v", limit, offset, err)
	}

	limit, offset, err = ValidatePagination("20", "3", 100)
	if err != nil || limit != 20 || offset != 3 {
		t.Fatalf("unexpected pagination: %d %d %v", limit, offset, err)
	}

	limit, offset, err = ValidatePagination("999", "1", 100)
	if err != nil || limit != 100 || offset != 1 {
		t.Fatalf("unexpected capped limit: %d %d %v", limit, offset, err)
	}

	if _, _, err = ValidatePagination("0", "0", 100); err == nil {
		t.Fatalf("expected limit error")
	}
	if _, _, err = ValidatePagination("abc", "0", 100); err == nil {
		t.Fatalf("expected limit parse error")
	}
	if _, _, err = ValidatePagination("1", "-1", 100); err == nil {
		t.Fatalf("expected offset error")
	}
	if _, _, err = ValidatePagination("1", "abc", 100); err == nil {
		t.Fatalf("expected offset parse error")
	}
}
