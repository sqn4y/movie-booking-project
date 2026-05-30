package model

import (
	"testing"
	"time"
)

func TestRequestMovieMovie(t *testing.T) {
	releaseDate := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	req := RequestMovie{
		Title:       "Oppenheimer",
		Director:    "Christopher Nolan",
		Duration:    10800,
		Description: "Movie description",
		ImageURL:    "/static/images/oppenheimer.svg",
		GenreIDs:    []int64{1, 3},
		AgeRating:   18,
		ReleaseDate: releaseDate,
	}

	movie := req.Movie()

	if movie.Title != req.Title {
		t.Fatalf("expected title %q, got %q", req.Title, movie.Title)
	}
	if movie.Director != req.Director {
		t.Fatalf("expected director %q, got %q", req.Director, movie.Director)
	}
	if movie.Duration != req.Duration {
		t.Fatalf("expected duration %d, got %d", req.Duration, movie.Duration)
	}
	if movie.ImageURL != req.ImageURL {
		t.Fatalf("expected image url %q, got %q", req.ImageURL, movie.ImageURL)
	}
	if movie.AgeRating != req.AgeRating {
		t.Fatalf("expected age rating %d, got %d", req.AgeRating, movie.AgeRating)
	}
	if !movie.ReleaseDate.Equal(releaseDate) {
		t.Fatalf("expected release date %s, got %s", releaseDate, movie.ReleaseDate)
	}
}

func TestToResponseMovies(t *testing.T) {
	movies := []*Movie{
		{Id: 1, Title: "First"},
		{Id: 2, Title: "Second"},
	}

	responses := ToResponseMovies(movies)

	if len(responses) != len(movies) {
		t.Fatalf("expected %d responses, got %d", len(movies), len(responses))
	}
	if responses[0] == nil || responses[1] == nil {
		t.Fatal("expected responses without nil elements")
	}
	if responses[0].Id != 1 || responses[1].Id != 2 {
		t.Fatalf("unexpected response ids: %+v", responses)
	}
}
