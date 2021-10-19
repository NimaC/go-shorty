package storage

import "time"

type Service interface {
	Save(string, string, time.Time) error
	Load(string) (*Item, error)
	Close() error
}

type Item struct {
	Id      string `json:"id" redis:"id"`
	URL     string `json:"url" redis:"url"`
	Expires string `json:"expires" redis:"expires"`
	Visits  int    `json:"visits" redis:"visits"`
}
