package handlers

import (
	"encoding/json"
	"net/http"
)

// Respond generates a http respond to user
func Respond(w http.ResponseWriter, data interface{}, contentType string) {
	w.Header().Add("Content-Type", contentType)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
