package main

import (
	"log"
	"os"

	"github.com/NimaC/go-shorty/handler"
	storage "github.com/NimaC/go-shorty/storage/redis"
)

func main() {
	host, port := os.Getenv("SHORTY_HOST"), os.Getenv("SHORTY_PORT")
	if host == "" || port == "" {
		log.Fatal("set SHORTY_HOST, SHORTY_PORT in environment variables")
	}
	service, err := storage.New()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	router := handler.New(host, service)
	router.Run(host + ":" + port)
}
