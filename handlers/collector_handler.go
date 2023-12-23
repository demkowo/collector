package handler

import (
	"fmt"
	"log"
	"net/http"

	collector "github.com/demkowo/goquery/services"
	"github.com/gin-gonic/gin"
)

type CollectorHandler interface {
	GatherAllLinks(*gin.Context)
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

func (h collectorHandler) GatherAllLinks(c *gin.Context) {
	log.Println("--- handler/GatherAllLinks() ---")

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

		if err := h.service.GatherAllLinks(data.Urls, data.Limits...); err != nil {
			fmt.Println("Can't gather links ", err)
			return
		}

		return
	}

	// Fetch values (url and limits)
	urls := c.PostFormArray("urls")
	limits := c.PostFormArray("limits")

	h.service.GatherAllLinks(urls, limits...)
}
