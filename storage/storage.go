package storage

import (
	"io"
	"time"
)

type Service interface {
	io.Closer
	Save(string, string, time.Time, int) error
	Load(string) (*Item, error)
	IsUsed(string) bool
}

type Item struct {
	Id      string    `json:"id" redis:"id"`
	URL     string    `json:"url" redis:"url"`
	Expires time.Time `json:"expires" redis:"expires"`
	Visits  int       `json:"visits" redis:"visits"`
}
