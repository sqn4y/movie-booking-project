package model

import "testing"

func TestRequestBookingBooking(t *testing.T) {
	req := RequestBooking{
		MovieId: 10,
		Seats:   []string{"A1", "A2"},
	}

	booking := req.Booking()

	if booking.MovieId != req.MovieId {
		t.Fatalf("expected movie id %d, got %d", req.MovieId, booking.MovieId)
	}
	if booking.Status != "pending" {
		t.Fatalf("expected status pending, got %q", booking.Status)
	}
	if booking.UserId.String() == "" {
		t.Fatal("expected generated user id")
	}
	if len(booking.Seats) != len(req.Seats) {
		t.Fatalf("expected %d seats, got %d", len(req.Seats), len(booking.Seats))
	}
}
