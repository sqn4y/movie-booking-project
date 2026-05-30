package service

import (
	"backend/internal/model"
	repomock "backend/internal/repository/mock"
	"context"
	"io"
	"log/slog"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestBookingServiceSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockBookingRepository(ctrl)
	svc := NewBookingService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	booking := &model.Booking{MovieId: 7, Seats: []string{"A1"}}
	repo.EXPECT().
		Save(gomock.Any(), booking).
		Return(&model.Booking{Id: 1, MovieId: 7, Seats: []string{"A1"}}, nil)

	saved, err := svc.Save(context.Background(), booking)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Id != 1 {
		t.Fatalf("expected id 1, got %d", saved.Id)
	}
}

func TestBookingServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockBookingRepository(ctrl)
	svc := NewBookingService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	status := "approved"
	update := &model.UpdateBooking{Seats: []string{"B1"}, Status: &status}
	repo.EXPECT().
		Update(gomock.Any(), int64(3), update).
		Return(&model.Booking{Id: 3, Status: status}, nil)

	updated, err := svc.Update(context.Background(), 3, update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != status {
		t.Fatalf("expected status %q, got %q", status, updated.Status)
	}
}

func TestBookingServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockBookingRepository(ctrl)
	svc := NewBookingService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))

	repo.EXPECT().
		Delete(gomock.Any(), int64(3)).
		Return(nil)

	if err := svc.Delete(context.Background(), 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
