package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Id        int64     `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	MovieId   int64     `json:"movie_id"`
	Seats     []string  `json:"seats"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateBooking struct {
	Seats  []string `json:"seats"`
	Status *string  `json:"status"`
}

type RequestBooking struct {
	MovieId int64    `json:"movie_id"`
	Seats   []string `json:"seats"`
}

type ResponseBooking struct {
	Id        int64     `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	MovieId   int64     `json:"movie_id"`
	Seats     []string  `json:"seats"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Booking) Response() ResponseBooking {
	return ResponseBooking{
		Id:        b.Id,
		UserId:    b.UserId,
		MovieId:   b.MovieId,
		Seats:     b.Seats,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (r *RequestBooking) Booking() *Booking {

	userId := uuid.New()
	return &Booking{
		MovieId: r.MovieId,
		Seats:   r.Seats,
		Status:  "pending",
		UserId:  userId,
	}
}
