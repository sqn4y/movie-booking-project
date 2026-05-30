package api

import (
	"backend/internal/model"
	servicemock "backend/internal/service/mock"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestMovieHandlerFindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	movieService := servicemock.NewMockMovieService(ctrl)
	handler := NewMovieHandler(movieService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	want := []*model.ResponseMovie{{Id: 1, Title: "Movie"}}
	movieService.EXPECT().
		FindAll(gomock.Any()).
		Return(want, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)

	handler.FindAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []*model.ResponseMovie
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestMovieHandlerSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	movieService := servicemock.NewMockMovieService(ctrl)
	handler := NewMovieHandler(movieService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	releaseDate := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	body := []byte(`{
		"title": "Movie",
		"director": "Director",
		"duration": 120,
		"description": "Description",
		"genre_ids": [1, 3],
		"age_rating": 18,
		"release_date": "2024-12-01T00:00:00Z"
	}`)
	want := &model.ResponseMovie{Id: 10, Title: "Movie", ReleaseDate: releaseDate}

	movieService.EXPECT().
		Save(gomock.Any(), gomock.Any(), []int64{1, 3}).
		DoAndReturn(func(_ context.Context, movie *model.Movie, genreIDs []int64) (*model.ResponseMovie, error) {
			if movie.Title != "Movie" {
				t.Fatalf("expected title Movie, got %q", movie.Title)
			}
			if !movie.ReleaseDate.Equal(releaseDate) {
				t.Fatalf("expected release date %s, got %s", releaseDate, movie.ReleaseDate)
			}
			return want, nil
		})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movie", bytes.NewReader(body))

	handler.Save(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got model.ResponseMovie
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != want.Id || got.Title != want.Title {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestMovieHandlerSaveInvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	movieService := servicemock.NewMockMovieService(ctrl)
	handler := NewMovieHandler(movieService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movie", bytes.NewBufferString("{"))

	handler.Save(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
