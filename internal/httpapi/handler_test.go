package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeStore struct {
	lastValue int64
	numbers   []int64
	err       error
}

func (f *fakeStore) AddAndList(ctx context.Context, value int64) ([]int64, error) {
	f.lastValue = value
	return f.numbers, f.err
}

func TestHandleNumbersOK(t *testing.T) {
	store := &fakeStore{numbers: []int64{1, 2, 3}}
	handler := NewRouter(store)

	reqBody := []byte(`{"number":2}`)
	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if store.lastValue != 2 {
		t.Fatalf("expected value 2, got %d", store.lastValue)
	}

	var resp struct {
		Numbers []int64 `json:"numbers"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Numbers) != 3 || resp.Numbers[0] != 1 || resp.Numbers[2] != 3 {
		t.Fatalf("unexpected numbers: %#v", resp.Numbers)
	}
}

func TestHandleNumbersBadJSON(t *testing.T) {
	store := &fakeStore{}
	handler := NewRouter(store)

	req := httptest.NewRequest(http.MethodPost, "/numbers", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestHandleNumbersMethodNotAllowed(t *testing.T) {
	store := &fakeStore{}
	handler := NewRouter(store)

	req := httptest.NewRequest(http.MethodGet, "/numbers", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
