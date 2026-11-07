package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondError(t *testing.T) {
	response := httptest.NewRecorder()
	respondError(response, http.StatusBadRequest, "invalid request")

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != false || body["error"] != "REQUEST_FAILED" || body["message"] != "invalid request" {
		t.Fatalf("response = %#v", body)
	}
}
