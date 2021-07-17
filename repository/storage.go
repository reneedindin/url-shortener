package repository

import (
	"time"

	"url-shortener/model"
)

type Storage interface {
	SaveClientIP(clientIP string, expireAt time.Time) error
	LoadClientIP(clientIP string) (int, error)
	Save(urlID, url, expireAt string) error
	Load(urlID string) (model.URL, error)
}
