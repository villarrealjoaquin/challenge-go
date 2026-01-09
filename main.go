package main

import (
	"fmt"

	"educabot.com/bookshop/config"
	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/providers"
	"educabot.com/bookshop/repositories"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.New()
	router.SetTrustedProxies(nil)

	booksProvider := providers.NewHTTPBooksProvider(cfg.BooksAPIURL)
	booksRepo := repositories.NewBooksRepository(booksProvider)
	metricsService := services.NewMetricsService(booksRepo)
	metricsHandler := handlers.NewGetMetrics(metricsService)

	router.GET("/", metricsHandler.Handle())

	fmt.Println("Starting server on :3000")
	router.Run(":3000")
}
