package repository

import (
	"backend/internal/config"
	"backend/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type MovieRepository interface {
	FindAll(ctx context.Context) ([]*model.Movie, error)
	Save(ctx context.Context, movie *model.Movie, genreIds []int64) (*model.Movie, error)
	GetGenresByMovieId(ctx context.Context, movieId int64) ([]*model.Genre, error)
}

type movieRepository struct {
	db     *config.NativeDatabase
	logger *slog.Logger
}

func NewMovieRepository(db *config.NativeDatabase, logger *slog.Logger) MovieRepository {
	return &movieRepository{
		db:     db,
		logger: logger,
	}
}

func (m *movieRepository) FindAll(ctx context.Context) ([]*model.Movie, error) {
	query := `
		SELECT 
			m.id,
			m.title,
			m.director,
			m.duration,
			m.description,
			m.image_url,
			m.age_rating,
			m.release_date,
			m.created_at,
			m.updated_at,
			COALESCE(
				(SELECT JSON_AGG(
					JSON_BUILD_OBJECT('id', g.id, 'name', g.name)
				)
				FROM movie_genre mg
				JOIN genre g ON mg.genre_id = g.id
				WHERE mg.movie_id = m.id),
				'[]'
			) as genres
		FROM movies m
		ORDER BY m.id
	`

	rows, err := m.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]*model.Movie, 0)

	for rows.Next() {
		movie := model.Movie{}
		var genresJSON []byte

		err = rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Director,
			&movie.Duration,
			&movie.Description,
			&movie.ImageURL,
			&movie.AgeRating,
			&movie.ReleaseDate,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&genresJSON,
		)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(genresJSON, &movie.Genres); err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *movieRepository) Save(ctx context.Context, movie *model.Movie, genreIds []int64) (*model.Movie, error) {
	tx, err := m.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	query := `insert into movies (
                    title,
                    director,
                    duration,
                    description,
                    image_url,
                    age_rating,
                    release_date,
                    created_at,
                    updated_at
		      ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err = tx.QueryRowContext(ctx, query,
		movie.Title,
		movie.Director,
		movie.Duration,
		movie.Description,
		movie.ImageURL,
		movie.AgeRating,
		movie.ReleaseDate,
		now,
		now).Scan(&movie.Id)

	if err != nil {
		return nil, err
	}

	err = addGenresToMovie(ctx, tx, movie.Id, genreIds)
	if err != nil {
		return nil, err
	}

	return movie, tx.Commit()
}

func (m *movieRepository) GetGenresByMovieId(ctx context.Context, movieId int64) ([]*model.Genre, error) {
	query := `SELECT g.* FROM movie_genre mg JOIN genre g on mg.genre_id = g.id WHERE mg.movie_id = $1`

	rows, err := m.db.DB.QueryContext(ctx, query, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]*model.Genre, 0)
	for rows.Next() {
		var genre model.Genre
		err = rows.Scan(
			&genre.Id,
			&genre.Name,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return genres, nil
}

func addGenresToMovie(ctx context.Context, tx *sql.Tx, movieId int64, genreIds []int64) error {
	if len(genreIds) == 0 {
		return nil
	}

	query := `INSERT INTO movie_genre (movie_id, genre_id) VALUES `

	var args []interface{}
	var valueStrings []string

	for i, genreID := range genreIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, movieId, genreID)
	}

	query += strings.Join(valueStrings, ", ")

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
