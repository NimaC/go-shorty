package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/NimaC/go-shorty/storage"
	redis "github.com/go-redis/redis/v8"
)

type myRedis struct{ client *redis.Client }

type redisClientConfig struct {
	host     string
	port     string
	password string
}

var ctx = context.Background()

func New() (storage.Service, error) {
	clientConfig, err := getRedisClientConfig()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     clientConfig.host + ":" + clientConfig.port,
		Password: clientConfig.password,
		DB:       0,
	})
	return &myRedis{client}, nil
}

func getRedisClientConfig() (*redisClientConfig, error) {
	host, port, password := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PW")
	config := redisClientConfig{host, port, password}
	if host == "" || port == "" {
		return &config, errors.New("set REDIS_HOST, REDIS_PORT and REDIS_PW in environment variables")
	}
	return &config, nil
}

func (r *myRedis) IsUsed(id string) bool {
	_, err := r.client.Get(ctx, id).Result()
	if err == redis.Nil {
		return false
	} else {
		return true
	}
}

func (r *myRedis) Save(id string, url string, expires time.Time, count int) error {
	shortLink := storage.Item{id, url, expires, count}
	p, err := json.Marshal(shortLink)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, id, p, time.Duration(expires.UnixMilli())).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *myRedis) Load(id string) (*storage.Item, error) {
	var item storage.Item
	itemString, err := r.client.Get(ctx, id).Result()
	if err != nil {
		return &item, err
	}
	err = json.Unmarshal([]byte(itemString), &item)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (r *myRedis) Close() error {
	return r.client.Close()
}
