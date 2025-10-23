package handler

import (
	"calendar-server/internal/model"
	"calendar-server/internal/service"
	"calendar-server/internal/utils"
	"encoding/json"
	"log"
	"net/http"
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
		h.utils.WriteErr(w, http.StatusServiceUnavailable, "service error")
		return
	}

	h.utils.WriteJSON(w, http.StatusOK, map[string]any{"result": "event created", "id": id})
}
