package main

import (
	"log"
	"net/http"
	"short_url/internal/application/url"
	"short_url/internal/infrastructure/database/postgres"
	httpinfra "short_url/internal/infrastructure/http"
	"time"
)

func main() {
	dbPostgres := postgres.ConnectionDatabase()
	defer dbPostgres.Close()

	urlRepo := postgres.NewUrlRepository(dbPostgres) // this way, we can easily change de db
	urlSvc := url.NewService(urlRepo)
	urlHandler := httpinfra.NewUrlHandler(urlSvc)
	router := httpinfra.NewRouter(urlHandler)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	router.Run()
}
