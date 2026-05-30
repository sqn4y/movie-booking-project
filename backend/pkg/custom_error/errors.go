package custom_error

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var ErrValidation = errors.New("validation error")
var ErrNotFound = errors.New("not found")

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAppError(status int, code, message string) AppError {
	return AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func MapError(err error) AppError {
	switch {
	case errors.Is(err, ErrNotFound):
		return NewAppError(http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrValidation):
		return NewAppError(http.StatusBadRequest, "validation_error", err.Error())
	default:
		return NewAppError(http.StatusInternalServerError, "internal_error", "internal app error")
	}
}

func WriteAppError(w http.ResponseWriter, appErr AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	_ = json.NewEncoder(w).Encode(appErr)
}

func HandleError(w http.ResponseWriter, logger *slog.Logger, err error) bool {
	if err == nil {
		return false
	}

	appErr := MapError(err)
	if appErr.Status == http.StatusInternalServerError {
		logger.Error("request failed", "error", err)
	}

	WriteAppError(w, appErr)
	return true
}
