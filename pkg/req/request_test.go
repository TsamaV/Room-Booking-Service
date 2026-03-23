package req_test

import (
	"Avito/pkg/req"
	"bytes"
	"net/http/httptest"
	"testing"
)

type TestPayload struct {
	Name string `json:"name"`
}

func TestDecodeSuccess(t *testing.T) {
	body := bytes.NewBufferString(`{"name": "test"}`)
	r := httptest.NewRequest("POST", "/", body)

	payload, err := req.Decode[TestPayload](r)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Name != "test" {
		t.Fatalf("expected test, got %s", payload.Name)
	}
}

func TestDecodeInvalid(t *testing.T) {
	body := bytes.NewBufferString(`invalid json`)
	r := httptest.NewRequest("POST", "/", body)

	_, err := req.Decode[TestPayload](r)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}