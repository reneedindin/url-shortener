package main

import (
	"log"
	"net/http"

	"url-shortener/handler"
	"url-shortener/repository/redis"
)

func main() {
	cfg := &redis.Config{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 50,
		Timeout:  10000,
	}
	client, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	storageRedis := redis.NewStorageRedis(client)
	router := handler.New(storageRedis)
	http.ListenAndServe(":8080", router)
}
