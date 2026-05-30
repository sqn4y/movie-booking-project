package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"log/slog"
)

type BookingService interface {
	Save(ctx context.Context, booking *model.Booking) (*model.Booking, error)
	Update(ctx context.Context, bookingId int64, booking *model.UpdateBooking) (*model.Booking, error)
	Delete(ctx context.Context, bookingId int64) error
}

type bookingService struct {
	logger            *slog.Logger
	bookingRepository repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, logger *slog.Logger) BookingService {
	return &bookingService{
		logger:            logger,
		bookingRepository: bookingRepo,
	}
}

func (b *bookingService) Save(ctx context.Context, booking *model.Booking) (*model.Booking, error) {
	return b.bookingRepository.Save(ctx, booking)
}

func (b *bookingService) Update(ctx context.Context, bookingId int64, booking *model.UpdateBooking) (*model.Booking, error) {
	return b.bookingRepository.Update(ctx, bookingId, booking)
}

func (b *bookingService) Delete(ctx context.Context, bookingId int64) error {
	return b.bookingRepository.Delete(ctx, bookingId)
}
