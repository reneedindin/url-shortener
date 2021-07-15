package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"url-shortener/model"
	"url-shortener/repository"
	"url-shortener/util"
)

const urlHost = "http://localhost:8080/%s"

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
	w.Header().Set("Content-Type", "application/json")

	var shortURLInfo model.ShortURLInfo
	if err := json.NewDecoder(r.Body).Decode(&shortURLInfo); err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}

	urlID := util.GenUrlID(shortURLInfo.URL, shortURLInfo.ExpireAt)
	if err := h.storage.Save(urlID, shortURLInfo.URL, shortURLInfo.ExpireAt); err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}

	data := model.URL{
		ID: urlID,
		ShortUrl: fmt.Sprintf(urlHost, urlID),
	}
	Response(w, http.StatusOK, data)
}

func (h *handler)RedirectURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	urlID, ok := params["url_id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := h.storage.Load(urlID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.ShortUrl, http.StatusPermanentRedirect)
}
