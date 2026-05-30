package router

import (
	"backend/internal/api"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	m.Use(loggingMiddleware(logger))
	return m
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
