package redis

import (
	"time"

	"github.com/go-redis/redis"
	errors "golang.org/x/xerrors"
	"url-shortener/model"
)

func NewRedisClient(cfg *Config) (client *redis.Client, err error) {
	timeout := time.Millisecond * time.Duration(cfg.Timeout)
	client = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  timeout,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
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

func (s *StorageRedis) SaveClientIP(clientIP string, expireAt time.Time) error {
	if status := s.client.Incr(clientIP); status.Err() != nil {
		return status.Err()
	}

	if status := s.client.ExpireAt(clientIP, expireAt); status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (s *StorageRedis) LoadClientIP(clientIP string) (int, error) {
	result := s.client.Get(clientIP)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return 0, nil
		}
		return 0, result.Err()
	}
	return result.Int()
}

func (s *StorageRedis) Save(urlID, url, expireAt string) error {
	if status := s.client.Set(urlID, url, 0); status.Err() != nil {
		return status.Err()
	}

	expireAtTime, _ := time.Parse(time.RFC3339, expireAt)
	if status := s.client.ExpireAt(urlID, expireAtTime); status.Err() != nil {
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
		ID:       urlID,
		ShortUrl: result.Val(),
	}
	return shortURL, nil
}
