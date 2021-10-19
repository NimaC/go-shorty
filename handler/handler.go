package handler

import (
	"encoding/json"
	"fmt"
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

func New(host string, storage storage.Service) *gin.Engine {
	router := gin.Default()

	h := handler{host, storage}
	router.POST("/encode/", responseHandler(h.encode))
	router.GET("/:shortLink", h.redirect)
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
		URL     string `json:"url"`
		Expires string `json:"expires"`
	}

	rawData, err := ctx.GetRawData()

	if err == nil {
		panic(err)
	}

	if err := json.Unmarshal(rawData, &input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}

	uri, err := url.ParseRequestURI(input.URL)

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid url")
	}

	layoutISO := "2006-01-02 15:04:05"
	expires, err := time.Parse(layoutISO, input.Expires)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid expiration date")
	}

	id := crypto.Encode([]byte(uri.String()))

	err = h.storage.Save(id, uri.String(), expires)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}

	u := url.URL{
		Host: h.host,
		Path: id,
	}

	fmt.Printf("Generated link: %v \n", u.String())

	return u.String(), http.StatusCreated, nil
}

func (h handler) decode(ctx *gin.Context) (interface{}, int, error) {
	code := ctx.Param("name")

	model, err := h.storage.Load(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return model, http.StatusOK, nil
}

func (h handler) redirect(ctx *gin.Context) {
	code := ctx.Param("name")

	item, err := h.storage.Load(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response{Data: nil, Success: err == nil})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, item.URL)
}
