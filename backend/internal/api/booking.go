package api

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg"
	"backend/pkg/custom_error"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type BookingHandler struct {
	logger         *slog.Logger
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{
		logger:         logger,
		bookingService: bookingService,
	}
}

func (h *BookingHandler) Save(w http.ResponseWriter, r *http.Request) {
	request, ok := h.decodeBooking(w, r)
	if !ok {
		return
	}

	saved, err := h.bookingService.Save(r.Context(), request.Booking())
	if h.handleError(w, err) {
		return
	}
	pkg.WriteJSON(w, http.StatusCreated, saved)
}

func (h *BookingHandler) Update(w http.ResponseWriter, r *http.Request) {
	bookingId, ok := pkg.ParseID(w, r)
	if !ok {
		return
	}

	var req model.UpdateBooking
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("%w: %s", custom_error.ErrValidation, err.Error()))
		return
	}

	updated, err := h.bookingService.Update(r.Context(), bookingId, &req)
	if h.handleError(w, err) {
		return
	}
	pkg.WriteJSON(w, http.StatusOK, updated.Response())
}

func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	bookingId, ok := pkg.ParseID(w, r)
	if !ok {
		return
	}
	err := h.bookingService.Delete(r.Context(), bookingId)
	if h.handleError(w, err) {
		return
	}
	pkg.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *BookingHandler) decodeBooking(w http.ResponseWriter, r *http.Request) (*model.RequestBooking, bool) {
	defer r.Body.Close()

	var req model.RequestBooking
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("%w: %s", custom_error.ErrValidation, err.Error()))
		return nil, false
	}
	return &req, true
}

func (h *BookingHandler) handleError(w http.ResponseWriter, err error) bool {
	return custom_error.HandleError(w, h.logger, err)
}
