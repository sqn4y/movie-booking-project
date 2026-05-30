package router

import (
	"backend/internal/api"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const URL = "/api/v1"

func Create(bookingHandler *api.BookingHandler, movieHandler *api.MovieHandler, logger *slog.Logger) *mux.Router {

	m := mux.NewRouter()
	router := m.PathPrefix(URL).Subrouter()

	router.HandleFunc("/movie", movieHandler.Save).Methods("POST")
	router.HandleFunc("/movies", movieHandler.FindAll).Methods("GET")

	router.HandleFunc("/booking", bookingHandler.Save).Methods("POST")
	router.HandleFunc("/booking/{id}", bookingHandler.Delete).Methods("DELETE")
	router.HandleFunc("/booking/{id}", bookingHandler.Update).Methods("PUT")

	m.HandleFunc("/swagger/doc.json", swaggerDoc).Methods("GET")
	m.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m.Use(loggingMiddleware(logger))
	return m
}

func swaggerDoc(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(swaggerDocPath())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func swaggerDocPath() string {
	if _, err := os.Stat("docs/openapi.json"); err == nil {
		return "docs/openapi.json"
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "docs/openapi.json"
	}

	return filepath.Join(filepath.Dir(file), "..", "..", "docs", "openapi.json")
}

func loggingMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(recorder, r)

			logger.Info("HTTP REQUEST:",
				"METHOD:", r.Method,
				"PATH:", r.URL.Path,
				"STATUS:", recorder.status,
				"DURATION:", time.Since(started).String(),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
