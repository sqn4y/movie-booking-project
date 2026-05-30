package api

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg"
	"backend/pkg/custom_error"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type MovieHandler struct {
	logger       *slog.Logger
	movieService service.MovieService
}

func NewMovieHandler(movieService service.MovieService, logger *slog.Logger) *MovieHandler {
	return &MovieHandler{
		logger:       logger,
		movieService: movieService,
	}
}

func (h *MovieHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieService.FindAll(r.Context())
	if h.handleError(w, err) {
		return
	}

	pkg.WriteJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) Save(w http.ResponseWriter, r *http.Request) {

	req, ok := h.decodeMovie(w, r)
	if !ok {
		return
	}

	saved, err := h.movieService.Save(r.Context(), req.Movie(), req.GenreIDs)
	if h.handleError(w, err) {
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, saved)
}

func (h *MovieHandler) decodeMovie(w http.ResponseWriter, r *http.Request) (*model.RequestMovie, bool) {
	defer r.Body.Close()

	var req model.RequestMovie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("%w: %s", custom_error.ErrValidation, err.Error()))
		return nil, false
	}
	return &req, true
}

func (h *MovieHandler) handleError(w http.ResponseWriter, err error) bool {
	return custom_error.HandleError(w, h.logger, err)
}
