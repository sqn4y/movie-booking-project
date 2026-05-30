package api

import (
	"backend/internal/model"
	servicemock "backend/internal/service/mock"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestBookingHandlerSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingService := servicemock.NewMockBookingService(ctrl)
	handler := NewBookingHandler(bookingService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	userID := uuid.New()
	bookingService.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, booking *model.Booking) (*model.Booking, error) {
			if booking.MovieId != 7 {
				t.Fatalf("expected movie id 7, got %d", booking.MovieId)
			}
			if booking.Status != "pending" {
				t.Fatalf("expected pending status, got %q", booking.Status)
			}
			booking.Id = 1
			booking.UserId = userID
			return booking, nil
		})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/booking", bytes.NewBufferString(`{"movie_id":7,"seats":["A1"]}`))

	handler.Save(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got model.Booking
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != 1 || got.UserId != userID {
		t.Fatalf("unexpected booking: %+v", got)
	}
}

func TestBookingHandlerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingService := servicemock.NewMockBookingService(ctrl)
	handler := NewBookingHandler(bookingService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	status := "approved"
	bookingService.EXPECT().
		Update(gomock.Any(), int64(5), gomock.Any()).
		DoAndReturn(func(_ context.Context, id int64, update *model.UpdateBooking) (*model.Booking, error) {
			if id != 5 {
				t.Fatalf("expected id 5, got %d", id)
			}
			if update.Status == nil || *update.Status != status {
				t.Fatalf("expected status %q, got %+v", status, update.Status)
			}
			return &model.Booking{Id: id, Status: status, Seats: update.Seats}, nil
		})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/booking/5", bytes.NewBufferString(`{"seats":["B1"],"status":"approved"}`))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestBookingHandlerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingService := servicemock.NewMockBookingService(ctrl)
	handler := NewBookingHandler(bookingService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	bookingService.EXPECT().
		Delete(gomock.Any(), int64(5)).
		Return(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/booking/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestBookingHandlerDeleteInvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingService := servicemock.NewMockBookingService(ctrl)
	handler := NewBookingHandler(bookingService, slog.New(slog.NewTextHandler(io.Discard, nil)))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/booking/bad", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "bad"})

	handler.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
