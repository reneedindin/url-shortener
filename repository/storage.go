package repository

import "url-shortener/model"

type Storage interface {
	Save(url string, expireAt string) error
	Load(urlID string) (model.URL, error)
}
