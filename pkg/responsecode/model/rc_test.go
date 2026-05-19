package rcmodel

import "testing"

func TestResponseCode(t *testing.T) {
	responseCode := NewResponseCode(200, "success", "2000000")
	if responseCode.StatusCode != 200 || responseCode.ResponseCode != "2000000" || responseCode.Error() != "success" {
		t.Fatalf("unexpected response code: %+v", responseCode)
	}

	responseCode.ResponseMessage = ""
	responseCode.ResponseDesc = "description"
	if responseCode.Error() != "description" {
		t.Fatalf("unexpected response description")
	}
}
