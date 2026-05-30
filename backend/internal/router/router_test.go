package router

import (
	"backend/internal/api"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSwaggerDoc(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	r := Create(api.NewBookingHandler(nil, logger), api.NewMovieHandler(nil, logger), logger)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var doc map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&doc); err != nil {
		t.Fatalf("decode swagger doc: %v", err)
	}
	if doc["openapi"] != "3.0.3" {
		t.Fatalf("expected openapi 3.0.3, got %v", doc["openapi"])
	}
}
