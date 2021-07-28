package modules

import (
	"encoding/json"
	"net/http"
)

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status  string
		Message string
	}{"1", "success"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
