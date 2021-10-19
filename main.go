package main

import (
	"log"
	"os"

	"github.com/NimaC/go-shorty/handler"
	storage "github.com/NimaC/go-shorty/storage/redis"
)

func main() {
	host, port := os.Getenv("shortyhost"), os.Getenv("shortyport")
	if host == "" || port == "" {
		log.Fatal("set shortyhost, shortyport in environment variables")
	}
	service, err := storage.New()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	router := handler.New(host, service)
	router.Run(host + ":" + port)
}
