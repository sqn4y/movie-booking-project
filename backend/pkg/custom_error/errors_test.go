package custom_error

import (
	"errors"
	"net/http"
	"testing"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{name: "not found", err: ErrNotFound, status: http.StatusNotFound, code: "not_found"},
		{name: "validation", err: ErrValidation, status: http.StatusBadRequest, code: "validation_error"},
		{name: "internal", err: errors.New("db failed"), status: http.StatusInternalServerError, code: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := MapError(tt.err)
			if appErr.Status != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, appErr.Status)
			}
			if appErr.Code != tt.code {
				t.Fatalf("expected code %q, got %q", tt.code, appErr.Code)
			}
		})
	}
}
