package redis

import (
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

func (s *StorageRedis) Save(url string, expireAt string) error {
	return nil
}

func (s *StorageRedis) Load(urlID string) (model.URL, error) {
	return model.URL{}, nil
}
