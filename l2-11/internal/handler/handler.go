package handler

import (
	"calendar-server/internal/errs"
	"calendar-server/internal/model"
	"calendar-server/internal/service"
	"calendar-server/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service service.Service
	utils   utils.Utils
}

func New(service service.Service, utils utils.Utils) Handler {
	return Handler{
		service: service,
		utils:   utils,
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON decode error:", err)
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if req.UserID == 0 || req.Title == "" || req.Date == "" {
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}

	_, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	id, err := h.service.CreateEvent(r.Context(), req)
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"result": "event created", "id": id})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var event model.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println("JSON decode error:", err)
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if event.Title == "" || event.Date == "" {
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}

	_, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	err = h.service.UpdateEvent(r.Context(), event)
	if err != nil {
		log.Println(err)
		if errors.Is(err, errs.ErrEventNotFound) {
			h.utils.WriteErr(w, http.StatusNotFound, "event not found")
			return
		}
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"result": "event updated", "id": event.ID})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.DeleteEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON decode error:", err)
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if req.ID == 0 {
		h.utils.WriteErr(w, http.StatusBadRequest, "invalid json body")
		return
	}

	err := h.service.DeleteEvent(r.Context(), req.ID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, errs.ErrEventNotFound) {
			h.utils.WriteErr(w, http.StatusNotFound, "event not found")
			return
		}
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"result": "event deleted", "id": req.ID})
}

func (h *Handler) EventsForToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusBadRequest, "error user id")
		return
	}

	events, err := h.service.GetEventsForPeriod(r.Context(), userID, 0)
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"daily_events": events})
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var week = 7

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusBadRequest, "error user id")
		return
	}

	events, err := h.service.GetEventsForPeriod(r.Context(), userID, week)
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"weekly_events": events})
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.utils.WriteErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var month = 30

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusBadRequest, "error user id")
		return
	}

	events, err := h.service.GetEventsForPeriod(r.Context(), userID, month)
	if err != nil {
		log.Println(err)
		h.utils.WriteErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"monthly_events": events})
}
