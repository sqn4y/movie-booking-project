package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestParseID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/booking/12", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "12"})
	rec := httptest.NewRecorder()

	id, ok := ParseID(rec, req)
	if !ok {
		t.Fatal("expected id to be parsed")
	}
	if id != 12 {
		t.Fatalf("expected id 12, got %d", id)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestParseIDInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/booking/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	id, ok := ParseID(rec, req)
	if ok {
		t.Fatal("expected parsing to fail")
	}
	if id != 0 {
		t.Fatalf("expected zero id, got %d", id)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
