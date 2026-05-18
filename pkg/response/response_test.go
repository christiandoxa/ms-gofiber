package response

import "testing"

func TestEnvelope(t *testing.T) {
	if got := Success("ok"); got.Status != "success" || got.Data != "ok" {
		t.Fatalf("unexpected success envelope: %+v", got)
	}

	fields := map[string]string{"title": "required"}
	if got := Error("invalid", fields); got.Status != "error" || got.Fields["title"] != "required" {
		t.Fatalf("unexpected error envelope: %+v", got)
	}
}
