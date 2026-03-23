package res_test

import (
	"Avito/pkg/res"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	res.JSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatal("expected Content-Type application/json")
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var response map[string]map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"]["code"] != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %s", response["error"]["code"])
	}
}