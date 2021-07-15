package redis

import (
	"time"

	"github.com/go-redis/redis"
	errors "golang.org/x/xerrors"
	"url-shortener/model"
)


func NewRedisClient(cfg *Config) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
	})
	_, err = client.Ping().Result()
	if err != nil {
		return nil, errors.Errorf("fail to ping: %w", err)
	}
	return client, nil
}

type StorageRedis struct {
	client *redis.Client
}

func NewStorageRedis(client *redis.Client) (storage *StorageRedis) {
	return &StorageRedis{
		client: client,
	}
}

func (s *StorageRedis) Save(urlID, url, expireAt string) error {
	time, err := time.Parse(time.RFC3339, expireAt)
	if err != nil {
		return err
	}
	if status := s.client.Set(urlID, url, 0); status.Err() != nil {
		return status.Err()
	}

	if status := s.client.ExpireAt(urlID, time); status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (s *StorageRedis) Load(urlID string) (model.URL, error) {
	result := s.client.Get(urlID)
	if result.Err() != nil {
		return model.URL{}, result.Err()
	}

	shortURL := model.URL{
		ID: urlID,
		ShortUrl: result.Val(),
	}
	return shortURL, nil
}
