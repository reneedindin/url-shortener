package handler

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
