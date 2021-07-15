package repository

import "url-shortener/model"

type Storage interface {
	Save(urlID, url, expireAt string) error
	Load(urlID string) (model.URL, error)
}
