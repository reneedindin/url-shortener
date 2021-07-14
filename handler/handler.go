package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"url-shortener/repository"
)

func New(storage repository.Storage) *mux.Router {
	handler := &handler{
		storage: storage,
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/urls", handler.UploadURL).Methods("POST")
	r.HandleFunc("/{url_id}", handler.RedirectURL).Methods("GET")
	http.ListenAndServe(":8080", r)
	return r
}

type handler struct {
	storage repository.Storage
}

func (h *handler)UploadURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
	fmt.Printf("%T\n", w)
}

func (h *handler)RedirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
	w.Header().Set("Content-Type", "application/json")

	var response Response

	json.NewEncoder(w).Encode(response)
}
