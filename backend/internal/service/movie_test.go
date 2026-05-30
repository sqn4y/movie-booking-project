package service

import (
	"backend/internal/model"
	repomock "backend/internal/repository/mock"
	"context"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestMovieServiceFindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockMovieRepository(ctrl)
	svc := NewMovieService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	repo.EXPECT().
		FindAll(gomock.Any()).
		Return([]*model.Movie{{Id: 1, Title: "Movie"}}, nil)

	movies, err := svc.FindAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(movies) != 1 {
		t.Fatalf("expected 1 movie, got %d", len(movies))
	}
	if movies[0].Id != 1 || movies[0].Title != "Movie" {
		t.Fatalf("unexpected movie: %+v", movies[0])
	}
}

func TestMovieServiceSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockMovieRepository(ctrl)
	svc := NewMovieService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	genres := []*model.Genre{{Id: 1, Name: "Drama"}}
	genreIDs := []int64{1}

	repo.EXPECT().
		Save(gomock.Any(), gomock.Any(), genreIDs).
		DoAndReturn(func(_ context.Context, movie *model.Movie, ids []int64) (*model.Movie, error) {
			if !reflect.DeepEqual(ids, genreIDs) {
				t.Fatalf("expected genre ids %v, got %v", genreIDs, ids)
			}
			movie.Id = 15
			return movie, nil
		})
	repo.EXPECT().
		GetGenresByMovieId(gomock.Any(), int64(15)).
		Return(genres, nil)

	saved, err := svc.Save(context.Background(), &model.Movie{Title: "Movie"}, genreIDs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Id != 15 {
		t.Fatalf("expected saved id 15, got %d", saved.Id)
	}
	if !reflect.DeepEqual(saved.Genres, genres) {
		t.Fatalf("expected genres %+v, got %+v", genres, saved.Genres)
	}
}

func TestMovieServiceSaveReturnsRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockMovieRepository(ctrl)
	svc := NewMovieService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	wantErr := errors.New("save failed")
	repo.EXPECT().
		Save(gomock.Any(), gomock.Any(), []int64{1}).
		Return(nil, wantErr)

	_, err := svc.Save(context.Background(), &model.Movie{Title: "Movie"}, []int64{1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}
