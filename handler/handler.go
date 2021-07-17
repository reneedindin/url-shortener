package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"url-shortener/model"
	"url-shortener/repository"
	"url-shortener/util"
)

const (
	urlHost = "http://localhost:8080/%s"

	verifyUploadCountByDay = 100
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

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ips := strings.Split(strings.TrimSpace(xForwardedFor), ",")
	ip := ips[0]
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ip
}

func (h *handler) checkValidIP(r *http.Request) (bool, error) {
	count, err := h.storage.LoadClientIP(getClientIP(r))
	if err != nil {
		return false, err
	}
	if count >= verifyUploadCountByDay {
		return false, nil
	}
	return true, nil
}

func (h *handler) UploadURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortURLInfo model.ShortURLInfo
	if err := json.NewDecoder(r.Body).Decode(&shortURLInfo); err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}

	expireAt, err := time.Parse(time.RFC3339, shortURLInfo.ExpireAt)
	if err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}
	if expireAt.Before(time.Now()) {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	isValid, err := h.checkValidIP(r)
	if err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}
	if !isValid {
		Response(w, http.StatusUnauthorized, nil)
		return
	}

	urlID := util.GenUrlID(shortURLInfo.URL, shortURLInfo.ExpireAt)
	if err := h.storage.Save(urlID, shortURLInfo.URL, shortURLInfo.ExpireAt); err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err := h.storage.SaveClientIP(getClientIP(r), time.Now().Add(24*time.Hour)); err != nil {
		Response(w, http.StatusInternalServerError, nil)
		return
	}

	data := model.URL{
		ID:       urlID,
		ShortUrl: fmt.Sprintf(urlHost, urlID),
	}
	Response(w, http.StatusOK, data)
}

func (h *handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
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
