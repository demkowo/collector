package handler

import (
	"log"
	"net/http"

	collector "github.com/demkowo/goquery/services"
	"github.com/gin-gonic/gin"
)

type CollectorHandler interface {
	FindAllLinks(*gin.Context)
}

type collectorHandler struct {
	service collector.CollectorService
}

func New(collector collector.CollectorService) CollectorHandler {
	log.Println("--- handler/User() ---")

	return &collectorHandler{
		service: collector,
	}
}

func (h collectorHandler) FindAllLinks(c *gin.Context) {
	log.Println("--- handler/FindAllLinks() ---")

	// Check the Content-Type header to determine how to decode the data
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		// Handling JSON data
		var data struct {
			Urls   []string `json:"urls"`
			Limits []string `json:"limits"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid JSON": err.Error()})
			return
		}

		h.service.FindAllLinks(data.Urls, data.Limits...)

		return
	}

	// Fetch values (url and limits)
	urls := c.PostFormArray("urls")
	limits := c.PostFormArray("limits")

	h.service.FindAllLinks(urls, limits...)
}
