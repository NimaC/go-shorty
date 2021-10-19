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
	host, port, password := os.Getenv("redishost"), os.Getenv("redisport"), os.Getenv("redispw")
	config := redisClientConfig{host, port, password}
	if host == "" || port == "" || password == "" {
		return &config, errors.New("set redis-host, redis-port and redis-pw in environment variables")
	}
	return &config, nil
}

func (r *myRedis) isUsed(id string) bool {
	_, err := r.client.Get(ctx, id).Result()
	if err != nil {
		return true
	} else {
		return false
	}
}

func (r *myRedis) Save(id string, url string, expires time.Time) error {
	if r.isUsed(id) {
		return errors.New("ID Collision. URL (probably) already in Redis DB")
	}
	shortLink := storage.Item{id, url, expires.Format("2006-01-02 15:04:05.728046 +0300 EEST"), 0}
	p, err := json.Marshal(shortLink)
	if err != nil {
		return err
	}

	_, err = r.client.Set(ctx, id, p, time.Duration(expires.Unix())).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *myRedis) Load(id string) (*storage.Item, error) {
	var item storage.Item
	err := r.client.HGetAll(ctx, id).Scan(item)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (r *myRedis) Close() error {
	return r.client.Close()
}
