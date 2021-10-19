package storage

import (
	"fmt"
	"time"
)

//type redis struct{ pool *redisClient.Pool }

func New(host, port, password string) (storage.Service, error) {
	redis.NewClient()
	pool := &redisClient.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redisClient.Conn, error) {
			return redisClient.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &redis{pool}, nil
}
