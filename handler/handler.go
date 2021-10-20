package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/NimaC/go-shorty/crypto"
	"github.com/NimaC/go-shorty/storage"
	"github.com/gin-gonic/gin"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"shortUrl"`
}

type handler struct {
	host    string
	storage storage.Service
}

var layoutISO string = "2006-01-02T15:04:05"

func New(host string, storage storage.Service) *gin.Engine {
	router := gin.Default()

	h := handler{host, storage}
	router.POST("/encode/", responseHandler(h.encode))
	router.GET("/:shortLink", h.redirect)
	router.GET("/loadinfo/:shortLink", responseHandler(h.loadInfo))
	return router
}

func responseHandler(h func(c *gin.Context) (interface{}, int, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, status, err := h(ctx)
		if err != nil {
			data = err.Error()
		}
		ctx.JSON(status, response{Data: data, Success: err == nil})
	}
}

func (h handler) encode(ctx *gin.Context) (interface{}, int, error) {
	var input struct {
		URL     string    `json:"url"`
		Expires time.Time `json:"expires"`
	}

	rawData, err := ctx.GetRawData()

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(rawData, &input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}

	uri, err := url.ParseRequestURI(input.URL)

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid url")
	}

	id := crypto.Encode([]byte(uri.String()))
	if h.storage.IsUsed(id) {
		return nil, http.StatusBadRequest, fmt.Errorf("ID Collision. URL (probably) already in Redis DB")
	}
	err = h.storage.Save(id, uri.String(), input.Expires, 0)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not store in database: %v", err)
	}

	u := url.URL{
		Host: h.host,
		Path: id,
	}

	fmt.Printf("Generated link: %v \n", u.String())

	return u.String(), http.StatusCreated, nil
}

func (h handler) loadInfo(ctx *gin.Context) (interface{}, int, error) {
	code := ctx.Param("shortLink")
	item, err := h.storage.Load(code)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	return item, http.StatusOK, nil
}

func (h handler) redirect(ctx *gin.Context) {
	code := ctx.Param("shortLink")
	item, err := h.storage.Load(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response{Data: nil, Success: err == nil})
		return
	}
	newVisits := item.Visits + 1
	err = h.storage.Save(item.Id, item.URL, item.Expires, newVisits)
	if err != nil {
		log.Println(err)
	}

	ctx.Redirect(http.StatusMovedPermanently, item.URL)
}
