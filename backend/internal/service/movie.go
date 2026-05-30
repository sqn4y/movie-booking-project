package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"log/slog"
)

type MovieService interface {
	FindAll(ctx context.Context) ([]*model.ResponseMovie, error)
	Save(ctx context.Context, movie *model.Movie, genreIds []int64) (*model.ResponseMovie, error)
}

type movieService struct {
	logger          *slog.Logger
	movieRepository repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository, logger *slog.Logger) MovieService {
	return &movieService{
		logger:          logger,
		movieRepository: movieRepo,
	}
}

func (m *movieService) FindAll(ctx context.Context) ([]*model.ResponseMovie, error) {

	movies, err := m.movieRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return model.ToResponseMovies(movies), nil
}

func (m *movieService) Save(ctx context.Context, movie *model.Movie, genreIds []int64) (*model.ResponseMovie, error) {
	saved, err := m.movieRepository.Save(ctx, movie, genreIds)
	if err != nil {
		return nil, err
	}
	genresByMovieId, err := m.movieRepository.GetGenresByMovieId(ctx, movie.Id)
	if err != nil {
		return nil, err
	}

	response := saved.ResponseMovie()
	response.Genres = genresByMovieId

	return response, nil
}
