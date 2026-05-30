package model

import "time"

type Movie struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Duration    int64     `json:"duration"`
	Description string    `json:"description"`
	Genres      []Genre   `json:"genres"`
	AgeRating   int       `json:"age_rating"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Genre struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateMovie struct {
	Title       *string    `json:"title"`
	Director    *string    `json:"director"`
	Duration    *int64     `json:"duration"`
	Description *string    `json:"description"`
	GenreIDs    []int64    `json:"genre_ids"`
	AgeRating   *int       `json:"age_rating"`
	ReleaseDate *time.Time `json:"release_date"`
}

type RequestMovie struct {
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Duration    int64     `json:"duration"`
	Description string    `json:"description"`
	GenreIDs    []int64   `json:"genre_ids"`
	AgeRating   int       `json:"age_rating"`
	ReleaseDate time.Time `json:"release_date"`
}

type ResponseMovie struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Duration    int64     `json:"duration"`
	Description string    `json:"description"`
	Genres      []*Genre  `json:"genres"`
	AgeRating   int       `json:"age_rating"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *RequestMovie) Movie() *Movie {
	return &Movie{
		Title:       r.Title,
		Director:    r.Director,
		Duration:    r.Duration,
		Description: r.Description,
		AgeRating:   r.AgeRating,
		ReleaseDate: r.ReleaseDate,
	}
}

func (m *Movie) ResponseMovie() *ResponseMovie {
	return &ResponseMovie{
		Id:          m.Id,
		Title:       m.Title,
		Director:    m.Director,
		Duration:    m.Duration,
		Description: m.Description,
		AgeRating:   m.AgeRating,
		ReleaseDate: m.ReleaseDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToResponseMovies(movies []*Movie) []*ResponseMovie {
	res := make([]*ResponseMovie, 0, len(movies))
	for _, movie := range movies {
		res = append(res, movie.ResponseMovie())
	}
	return res
}
