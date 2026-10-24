package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "idenova/identity/internal/adapters/http"
	"idenova/identity/internal/adapters/postgres"
	"idenova/identity/internal/adapters/postgres/db"
	"idenova/identity/internal/application"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal(err)
	}

	defer pool.Close()

	repo := postgres.NewRepository(db.New(pool))

	handler := httpadapter.NewHandler(
		application.NewCreateOrganization(repo),
		application.NewGetOrganization(repo),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /organizations", handler.CreateOrganization)
	mux.HandleFunc("GET /organizations/{id}", handler.GetOrganization)

	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Identity Service running on :8081")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP server failed: %v", err)
		return
	}
}
