package app

import handler "github.com/demkowo/goquery/handlers"

func addRoutes(h handler.CollectorHandler) {
	router.POST("/", h.GatherAllLinks)
}
