package utils

import (
	"encoding/json"
	"net/http"
)

type Utils struct {
}

func New() Utils {
	return Utils{}
}

func (u *Utils) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (u *Utils) WriteErr(w http.ResponseWriter, status int, msg string) {
	u.WriteJSON(w, status, map[string]string{"error": msg})
}
