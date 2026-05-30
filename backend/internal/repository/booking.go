package repository

import (
	"backend/internal/config"
	"backend/internal/model"
	"backend/pkg/custom_error"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type BookingRepository interface {
	Save(ctx context.Context, booking *model.Booking) (*model.Booking, error)
	Update(ctx context.Context, bookingId int64, booking *model.UpdateBooking) (*model.Booking, error)
	Delete(ctx context.Context, bookingId int64) error
}

type bookingRepository struct {
	db     *config.NativeDatabase
	logger *slog.Logger
}

func NewBookingRepository(db *config.NativeDatabase, logger *slog.Logger) BookingRepository {
	return &bookingRepository{
		db:     db,
		logger: logger,
	}
}

func (b *bookingRepository) Save(ctx context.Context, booking *model.Booking) (*model.Booking, error) {
	tx, err := b.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	query := `INSERT INTO bookings(
                     user_id,
                     movie_id,
                     seats,
                     status,
                     created_at,
                     updated_at) VALUES ($1, $2, $3, $4, $5, $6) returning id`

	err = tx.QueryRowContext(ctx, query,
		booking.UserId,
		booking.MovieId,
		booking.Seats,
		booking.Status,
		now,
		now).Scan(&booking.Id)

	if err != nil {
		return nil, err
	}

	return booking, tx.Commit()
}

func (b *bookingRepository) Update(ctx context.Context, bookingId int64, updateData *model.UpdateBooking) (*model.Booking, error) {
	tx, err := b.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE bookings 
              SET seats = COALESCE($1, seats),
                  status = COALESCE($2, status),
                  updated_at = NOW()
              WHERE id = $3 
              RETURNING id, user_id, movie_id, seats, status, created_at, updated_at`

	var booking model.Booking
	err = tx.QueryRowContext(ctx, query,
		pq.Array(updateData.Seats),
		updateData.Status,
		bookingId,
	).Scan(
		&booking.Id,
		&booking.UserId,
		&booking.MovieId,
		pq.Array(&booking.Seats),
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_error.ErrNotFound
		}
		return nil, err
	}

	return &booking, tx.Commit()
}

func (b *bookingRepository) Delete(ctx context.Context, bookingId int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	res, err := b.db.DB.ExecContext(ctx, query, bookingId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_error.ErrNotFound
	}
	return nil
}
