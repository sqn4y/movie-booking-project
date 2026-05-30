package pkg

import (
	"backend/pkg/custom_error"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func ParseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		custom_error.WriteAppError(w, custom_error.NewAppError(http.StatusBadRequest, "validation_error", "id must be a positive integer"))
		return 0, false
	}
	return id, true
}
